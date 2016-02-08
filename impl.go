package impl

import (
	"io/ioutil"
	"log"
	"os"
)

type Interface struct {
	Methods []Method
}

func NewInterface(m []Method) *Interface {
	return &Interface{m}
}

type Method struct {
	Name       string
	Parameters []Parameter
}

func NewMethod(name string, params []Parameter) Method {
	return Method{name, params}
}

type Parameter struct {
	Name string
	Type string
}

func NewParameter(name, typeName string) Parameter {
	return Parameter{name, typeName}
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

	dl("\n")
	dl("finished init()\n")
}
