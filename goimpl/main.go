package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ajmesa9891/impl/impl"
)

var commandPattern = regexp.MustCompile(`//go:generate\s*goimp\s*(.*)\n`)

const (
	cmdName = "goimp"
)

func logFatalUsage(args []string) {
	log.Fatalf("Must pass exactly 3 arguments:\n"+
		"  (1) the file name (perhaps $GOFILE if using go:generate)\n"+
		"  (2) interface path (e.g., sort.Interface)\n"+
		"  (3) the receiver (e.g., 'r *Receiver')\n"+
		"but got %d arguments: %q.\n"+
		"visit https://github.com/ajmesa9891/impl for more details.", len(args), args)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("impl: ")

	args := os.Args[1:]
	if len(args) < 3 {
		logFatalUsage(args)
	}

	file := filepath.Join(".", args[0])
	interfacePath := args[1]
	receiver := strings.Join(args[2:], "")
	var w bytes.Buffer

	err := impl.Impl(interfacePath, receiver, &w)
	if err != nil {
		log.Fatalf("could not build scaffolding for interface path %q: %s\n%",
			interfacePath, err)
	}
	err = writeInterfaceScaffolding(file, interfacePath, w.String())
	if err != nil {
		log.Fatalf("could not write scaffolding to file: %v\nscaffolding:\n%s", err, w.String())
	}
	log.Printf("wrote interface scaffolding for %q in file %q\n", interfacePath, file)
}

func writeInterfaceScaffolding(inputPath, interfacePath, scaffolding string) error {
	newContent := ""
	content, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("reading file %q: %s", inputPath, err)
	}

	for _, line := range strings.Split(string(content), "\n") {
		if strings.Contains(line, cmdName) && strings.Contains(line, interfacePath) {
			newContent += scaffolding
		} else {
			newContent += line + "\n"
		}
	}

	err = ioutil.WriteFile(inputPath, []byte(newContent), 0)
	if err != nil {
		return fmt.Errorf("writing file %q: %s", inputPath, err)
	}

	return nil
}
