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
				dl("    is a FuncType with %d params\n", len(funcType.Params.List))
				params := make([]Parameter, 0, len(funcType.Params.List))
				for ip, param := range funcType.Params.List {
					dl("    %dth param has Names %v", ip, param.Names)
					dl("    %dth param is Type %v", ip, param.Type)
					if ident, ok := param.Type.(*ast.Ident); len(param.Names) > 0 && ok {
						dl("    %dth param was added", ip)
						params = append(params, NewParameter(param.Names[0].Name, ident.Name))
					} else {
						dl("    %dth param was NOT added", ip)
					}
				}
				methods = append(methods, NewMethod(field.Names[0].Name, params))
			}
		}
	}

	return NewInterface(methods), nil
}
