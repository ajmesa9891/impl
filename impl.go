package impl

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"impl/errs"
	"io"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

// BuildInterface generates a model Interface from the given internal
// or external  path. The path is expected to be in the format of
// <package>.<interface>. For example, "io.Reader" or
// "impl/test_data/panther.Clawable".
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

	for i, field := range interfaceType.Methods.List {
		dl("  %dth field with type %T and Names %v\n", i, field.Type, field.Names)
		funcType, isMethod := field.Type.(*ast.FuncType)
		if namesl := len(field.Names); namesl > 0 && isMethod {
			methods = append(methods, buildMethod(field.Names[0].Name, funcType))
		} else if ident, ok := field.Type.(*ast.Ident); ok {
			dl("    embedded interface field %q\n", ident.Name)
			embedded, err := BuildInterface(fmt.Sprintf("%s.%s", pkgPath, ident.Name))
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
func RenderInterface(i *Interface, receiver string, w io.Writer) error {
	var ugly bytes.Buffer
	methodTmpl, err := template.
		New("method").
		Funcs(template.FuncMap{"Receiver": func() string { return receiver }}).
		Parse(
		`func ({{Receiver}}) {{.Name}}` +
			`({{range .In}}{{.Name}} {{.Type}}, {{end}}) ` +
			`({{range .Out}}{{.Name}} {{.Type}}, {{end}}) {
			panic("TODO: implement this method") }`)
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
