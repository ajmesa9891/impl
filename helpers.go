package impl

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"impl/errs"
	"io"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
)

// parseImport splits impPath into the package part and the interface name part
// (e.g., splits "io.Reader" into "io" and "Reader")
func parseImport(impPath string) (pkgPath, interfaceName string, err error) {
	if len(strings.TrimSpace(impPath)) < 1 {
		return "", "", errs.NewInvalidImportFormatError("import path cannot be empty")
	}

	parts := strings.Split(impPath, ".")
	if len(parts) < 2 {
		return "", "", errs.NewInvalidImportFormatError(
			"interface must have at least two parts: package and name (e.g., \"io.Reader\") and had %d parts", len(parts))
	}

	pkgPath = strings.Trim(strings.Join(parts[:len(parts)-1], "."), ".")
	interfaceName = parts[len(parts)-1]
	return
}

// formatInterface formats the given path using "golang.org/x/tools/imports".
func formatInterface(path string) (string, error) {
	if len(strings.TrimSpace(path)) < 1 {
		return "", errs.NewEmptyInterfacePathError("invalid interface: empty interface path %q", path)
	}

	srcWithInterface := []byte(fmt.Sprintf("package p;var r %s", path))
	srcB, err := imports.Process("", srcWithInterface, nil)
	if err != nil {
		return "", errs.NewInvalidInterfacePathError("invalid interface: ", err)
	}

	src := string(srcB)
	i := strings.Index(src, "var r ") + len("var r ")
	parts := strings.Split(src[i:], "\n")
	if len(parts) < 1 {
		return "", fmt.Errorf("imports.Process behaved unexpectedly: expected a new line after the var declaration")
	}

	return parts[0], nil
}

// buildPackage returns a *build.Package from the given package path.
func buildPackage(pkgPath string) (pkg *build.Package, err error) {
	pkg, err = build.Import(pkgPath, "", 0)
	if err != nil {
		err = errs.NewCouldNotFindPackageError("could not find interface's package (%q): %s", pkgPath, err)
	}
	return
}

func interfaceTypeSpec(name string, pkg *build.Package) (ts *ast.TypeSpec, err error) {
	unparsedFiles := []string{}
	fset := token.NewFileSet()
	for _, fileName := range pkg.GoFiles {
		file, err := parser.ParseFile(fset, filepath.Join(pkg.Dir, fileName), nil, 0)
		if err != nil {
			unparsedFiles = append(unparsedFiles, fileName)
			continue
		}

		for _, decl := range file.Decls {
			if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
				for _, spec := range decl.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok && ts.Name.Name == name {
						return ts, nil
					}
				}
			}
		}
	}

	err = errs.NewInterfaceNotFoundError("could not find interface %q when parsing package %q",
		name, pkg.Name)
	if len(unparsedFiles) > 0 {
		err = errs.NewInterfaceNotFoundError("%s: the following files could not be parsed: %q",
			err, unparsedFiles)
	}
	return
}

func buildMethod(name string, funcType *ast.FuncType) Method {
	dl("    method has input parameters?\t%t - Adding them\n", funcType.Params != nil)
	in := buildParams(funcType.Params)
	dl("    method has results?\t%t - Adding them\n", funcType.Results != nil)
	out := buildParams(funcType.Results)
	return NewMethod(name, in, out)
}

func buildParams(fl *ast.FieldList) []Parameter {
	if fl == nil || fl.List == nil || len(fl.List) == 0 {
		dl("    nothing to add, empty list")
		return []Parameter{}
	}
	params := make([]Parameter, 0, len(fl.List))
	dl("    it has %d fields", len(fl.List))
	for ip, field := range fl.List {
		dl("    attempting to build parameters from %dth field of type %T with names %v",
			ip, field.Type, field.Names)
		typeName := getParamTypeName(field)
		if isUnamed := len(field.Names) == 0; isUnamed {
			params = append(params, NewParameter("", typeName))
			dl("    %dth unnamed field was added", ip)
		} else {
			// Multiple names indicate an "i, j int" situation.
			// 1 field, 1 type, multiple parameters.
			for jp, fieldName := range field.Names {
				params = append(params, NewParameter(fieldName.Name, typeName))
				dl("    %d-%dth field was added", ip, jp)
			}
		}
	}
	return params
}

func getParamTypeName(field *ast.Field) (typeName string) {
	switch fieldType := field.Type.(type) {
	case *ast.Ident:
		typeName = fieldType.Name
	case *ast.ArrayType:
		if ident, ok := fieldType.Elt.(*ast.Ident); ok {
			typeName = "[]" + ident.Name
		}
		dl("    field of type %T with Elt %T was NOT added",
			field.Type, fieldType.Elt)
	case *ast.StarExpr:
		if ident, ok := fieldType.X.(*ast.Ident); ok {
			typeName = "*" + ident.Name
			break
		}
		dl("    field of type %T was NOT added", field.Type)
	default:
		dl("    field of type %T was NOT added", field.Type)
	}
	return
}

// BuildInterface generates a model Interface from the given internal
// or external  path. The path is expected to be in the format of
// <package>.<interface>. For example, "io.Reader" or
// "impl/test_data/panther.Clawable".
func buildInterface(path string) (*Interface, error) {
	pkgPath, interfaceName, err := parseImport(path)
	if err != nil {
		return nil, err
	}
	interfaceName, err = formatInterface(interfaceName)
	if err != nil {
		return nil, err
	}
	pkg, err := buildPackage(pkgPath)
	if err != nil {
		return nil, err
	}
	typeSpec, err := interfaceTypeSpec(interfaceName, pkg)
	if err != nil {
		return nil, err
	}
	interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
	if !ok {
		return nil, errs.NewNotAnInterfaceError("%q is not an interface", typeSpec.Name.Name)
	}

	dl("Going through %d fields of %q\n", len(interfaceType.Methods.List), typeSpec.Name.Name)
	methods := make([]Method, 0, len(interfaceType.Methods.List))

	for i, field := range interfaceType.Methods.List {
		dl("  %dth field with type %T and Names %v\n", i, field.Type, field.Names)
		funcType, isMethod := field.Type.(*ast.FuncType)
		if namesl := len(field.Names); namesl > 0 && isMethod {
			methods = append(methods, buildMethod(field.Names[0].Name, funcType))
		} else if ident, ok := field.Type.(*ast.Ident); ok {
			dl("    embedded interface field %q\n", ident.Name)
			embedded, err := buildInterface(fmt.Sprintf("%s.%s", pkgPath, ident.Name))
			if err != nil {
				dl("      error building embedded interface %q: %s\n", ident.Name, err.Error())
				return nil, err
			}
			dl("    adding %d methods from embedded interface\n", len(embedded.Methods))
			for _, m := range embedded.Methods {
				methods = append(methods, m)
			}
		} else {
			dl("    unexpected field %q was not processed\n", ident.Name)
		}
	}

	return NewInterface(methods), nil
}

// RenderInterface writes scaffolding for the given interface using receiver
// as the receiver. It formats the source using goformat and inserts a panic
// where the implementation should go.
func renderInterface(i *Interface, receiver string, w io.Writer) error {
	var ugly bytes.Buffer
	methodTmpl, err := template.
		New("method").
		Funcs(template.FuncMap{"Receiver": func() string { return receiver }}).
		Parse(
		"func ({{Receiver}}) {{.Name}}" +
			"({{range .In}}{{.Name}} {{.Type}}, {{end}}) " +
			"{{if ne (len .Out) 0}}({{range .Out}}{{.Name}} {{.Type}}, {{end}}){{end}} {\n" +
			"panic(\"TODO: implement this method\") }\n\n")
	if err != nil {
		return fmt.Errorf("error building template (methods %v): %s\n", i.Methods, err)
	}

	for _, m := range i.Methods {
		dl("rendering method %q\n", m.Name)
		err := methodTmpl.Execute(&ugly, m)
		if err != nil {
			return fmt.Errorf("error rendering method %q (%v): %s\n", m.Name, m, err)
		}
	}
	pretty, err := format.Source(ugly.Bytes())
	if err != nil {
		return fmt.Errorf("error formatting source:\n%s\n: %s\n",
			ugly.Bytes(), err.Error())
	}
	_, err = w.Write(pretty)
	if err != nil {
		return fmt.Errorf("error writing the formatted source: %s\n", err)
	}
	return nil
}
