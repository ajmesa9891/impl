package impl

import (
	"fmt"
	"go/build"
	"impl/errs"
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
