package impl

import (
	"impl/errs"
	"reflect"
	"testing"
)

func TestFormatInterface(t *testing.T) {
	cases := []struct {
		in       string
		wantPath string
		wantErr  error
	}{
		{"io.Blinger", "io.Blinger", nil},
		{"io.Reader", "io.Reader", nil},
		{"imports.Madeuper", "imports.Madeuper", nil},

		{"", "", &errs.EmptyInterfacePathError{}},
		{" \n		", "", &errs.EmptyInterfacePathError{}},
		{"io..Reader", "", &errs.InvalidInterfacePathError{}},
	}
	for _, c := range cases {
		got, err := formatInterface(c.in)
		if reflect.TypeOf(c.wantErr) != reflect.TypeOf(err) {
			t.Errorf("formatInterface(%q): wanted error type \"%T\", got \"%T\"", c.in, c.wantErr, err)
		} else if got != c.wantPath {
			t.Errorf("formatInterface(%q) == %q, want %q", c.in, got, c.wantPath)
		}
	}
}

func TestParseImportInterface(t *testing.T) {
	cases := []struct {
		in            string
		wantPkg       string
		wantInterface string
		wantErr       error
	}{
		{"io.Blinger", "io", "Blinger", nil},
		{"io.Reader", "io", "Reader", nil},
		{"imports.Madeuper", "imports", "Madeuper", nil},
		{"io..Reader", "io", "Reader", nil},

		{"Reader", "", "", &errs.InvalidImportFormatError{}},
		{"", "", "", &errs.InvalidImportFormatError{}},
		{" 	\n", "", "", &errs.InvalidImportFormatError{}},
	}
	for _, c := range cases {
		gotPkg, gotInterface, gotErr := parseImport(c.in)
		if reflect.TypeOf(c.wantErr) != reflect.TypeOf(gotErr) {
			t.Errorf("parseImport(%q): wanted error type \"%T\", got \"%T\"", c.in, c.wantErr, gotErr)
		} else if gotPkg != c.wantPkg || gotInterface != c.wantInterface {
			t.Errorf("parseImport(%q) == (%q, %q, %T), want (%q, %q, %T)",
				c.in, gotPkg, gotInterface, gotErr, c.wantPkg, c.wantInterface, c.wantErr)
		}
	}
}

func TestBuildPackage(t *testing.T) {
	cases := []struct {
		in            string
		wantedPkgName string
		wantErr       error
	}{
		{"io", "io", nil},

		{"nonexistent", "", &errs.CouldNotFindPackageError{}},
	}
	for _, c := range cases {
		gotPkg, gotErr := buildPackage(c.in)
		if reflect.TypeOf(c.wantErr) != reflect.TypeOf(gotErr) {
			t.Errorf("buildPackage(%q): wanted error type \"%T\", got \"%T\"", c.in, c.wantErr, gotErr)
		} else if gotPkg.Name != c.wantedPkgName {
			t.Errorf("buildPackage(%q) == (%q, %T), want (%q, %T)",
				c.in, gotPkg.Name, gotErr, c.wantedPkgName, c.wantErr)
		}
	}
}
