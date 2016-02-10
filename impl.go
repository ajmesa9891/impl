// TODO: improve comments and debugging
// TODO: make embedding and interface of a different packages work.
package impl

import (
	"fmt"
	"go/ast"
	"impl/errs"
	"io/ioutil"
	"log"
	"os"
)

// TODO: CONTINUE: make embedded interfaces work.
func BuildInterface(path string) (*Interface, error) {
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

	// TODO: CONTINUE HERE AFTER THE OTHER CONTINUE: check for different types. Make embedded interfaces work.
	for i, field := range interfaceType.Methods.List {
		dl("  %dth field with type %T and Names %v\n", i, field.Type, field.Names)
		funcType, isMethod := field.Type.(*ast.FuncType)
		if namesl := len(field.Names); namesl > 0 && isMethod {
			methods = append(methods, buildMethod(field.Names[0].Name, funcType))
		} else if ident, ok := field.Type.(*ast.Ident); ok {
			dl("    embedded interface field %q\n", ident.Name)
			embedded, err := BuildInterface(fmt.Sprintf("%s.%s", pkgPath, ident.Name))
			if err != nil {
				dl("      error building embedded interface %q: %s", ident.Name, err.Error())
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

// debugL is the debug logger
var debugL *log.Logger

// dl is a convenience function that logs a formatted string
// using the debug logger
var dl func(format string, v ...interface{})

func init() {
	debug := false // TODO: get this from the command line!
	w := ioutil.Discard
	if debug {
		w = os.Stderr
	}

	debugL = log.New(w, "impl: DEBUG: ", log.Lshortfile)
	dl = debugL.Printf

	dl("finished init()\n")
}
