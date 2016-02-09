// TODO: need to take care of embedded interfaces!
package impl

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"impl/errs"
	"path/filepath"
	"strings"

	"golang.org/x/tools/imports"
)

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
					if spec, ok := spec.(*ast.TypeSpec); ok && spec.Name.Name == name {
						return spec, nil
					}
				}
			}
		}
	}

	err = errs.NewInterfaceNotFoundError("could not find %q when parsing package %q",
		name, pkg.Name)
	if len(unparsedFiles) > 0 {
		err = errs.NewInterfaceNotFoundError("%s: the following files could not be parsed: %q",
			err, unparsedFiles)
	}
	return
}

func buildInterface(ts *ast.TypeSpec) (*Interface, error) {
	interfaceType, ok := ts.Type.(*ast.InterfaceType)
	if !ok {
		return NewInterface(nil), errs.NewNotAnInterfaceError("%q is not an interface type", ts.Name)
	}

	dl("Going through %d methods\n", len(interfaceType.Methods.List))
	methods := make([]Method, 0, len(interfaceType.Methods.List))
	for i, field := range interfaceType.Methods.List {
		dl("  %dth method with Names %v\n", i, field.Names)
		if namesl := len(field.Names); namesl > 0 {
			if funcType, ok := field.Type.(*ast.FuncType); ok {
				dl("    has parameters?\t%t - Adding them\n", funcType.Params != nil)
				in := buildParams(funcType.Params)
				dl("    has results?\t%t - Adding them\n", funcType.Results != nil)
				out := buildParams(funcType.Results)
				methods = append(methods, NewMethod(field.Names[0].Name, in, out))
			}
		}
	}

	return NewInterface(methods), nil
}

func buildParams(fl *ast.FieldList) []Parameter {
	if fl == nil || fl.List == nil || len(fl.List) == 0 {
		dl("    nothing to add, empty list")
		return []Parameter{}
	}
	params := make([]Parameter, 0, len(fl.List))
	dl("    it has %d fields", len(fl.List))
	for ip, field := range fl.List {
		if ident, ok := field.Type.(*ast.Ident); ok {
			// No names indicate an unnamed return parameter.
			if len(field.Names) == 0 {
				params = append(params, NewParameter("", ident.Name))
				dl("    %dth unnamed field was added", ip)
				continue
			}
			// Multiple names indicate an "i, j int" situation.
			// 1 field, 1 type, multiple parameters.
			for jp, fieldName := range field.Names {
				params = append(params, NewParameter(fieldName.Name, ident.Name))
				dl("    %d-%dth field was added", ip, jp)
			}
		} else {
			dl("    %dth field was NOT added", ip)
		}
	}
	return params
}
