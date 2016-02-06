package impl

import (
	"fmt"
	"golang.org/x/tools/imports"
	"strings"
)

// formatInterface formats the given path using "golang.org/x/tools/imports"
func formatInterface(path string) (string, error) {
	if len(path) < 1 {
		return "", NewInvalidInterfacePathError("invalid interface: empty interface path %q", path)
	}

	srcWithInterface := []byte(fmt.Sprintf("package p;var r %s", path))
	srcB, err := imports.Process("", srcWithInterface, nil)
	if err != nil {
		return "", fmt.Errorf("invalid interface: ", err)
	}

	src := string(srcB)
	i := strings.Index(src, "var r ") + len("var r ")
	parts := strings.Split(src[i:], "\n")
	if len(parts) < 1 {
		return "", fmt.Errorf("imports.Process behaved unexpectedly: expected a new line after the var declaration")
	}

	return parts[0], nil
}
