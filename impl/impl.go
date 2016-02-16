package impl

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

// Impl is the main entry point for the impl package. It writes scaffolding
// for the given interface using receiver as the receiver. The path is
// expected to be in the format of <package>.<interface>. For example,
// "io.Reader" or "impl/test_data/panther.Clawable".
func Impl(path string, receiver string, w io.Writer) error {
	iface, err := buildInterface(path)
	if err != nil {
		return err
	}
	return renderInterface(iface, receiver, w)
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