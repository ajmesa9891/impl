package impl

import (
	"impl/errs"
	"reflect"
	"strings"
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

func TestInterfaceTypeSpec_FindsIt(t *testing.T) {
	cases := []struct {
		pkgPath       string
		interfaceName string
		wantErr       error
	}{
		{"impl/test_data/panther", "Clawable", nil},
		{"sort", "Interface", nil},

		{"impl/test_data/panther", "UnexistentName", &errs.InterfaceNotFoundError{}},
		{"impl/test_data/panther", "WithParseErrors", &errs.InterfaceNotFoundError{}},
	}

	for _, c := range cases {
		pkg, err := buildPackage(c.pkgPath)
		if err != nil {
			t.Errorf("interfaceTypeSpec(...) failed precondition: could load package with path %q", c.pkgPath)
		}
		gotSpec, gotErr := interfaceTypeSpec(c.interfaceName, pkg)
		if reflect.TypeOf(gotErr) != reflect.TypeOf(c.wantErr) {
			t.Errorf(`interfaceTypeSpec(%q, %q): wanted error type "%T", got "%T"`,
				c.interfaceName, c.pkgPath, c.wantErr, gotErr)
		} else if c.wantErr != nil {
			continue // The error match passed. Nothing more to test.
		} else if gotSpec == nil {
			t.Errorf(`interfaceTypeSpec(%q, %q) == (nil, _), wanted (<TypeSpec with name %q>, _)`,
				c.interfaceName, c.pkgPath, c.interfaceName)
		} else if gotSpec.Name.Name != c.interfaceName {
			t.Errorf("interfaceTypeSpec(%q, %q) == (%q, _), wanted (%q, _)",
				c.interfaceName, c.pkgPath, gotSpec.Name.Name, c.interfaceName)
		}
	}
}

func TestInterfaceTypeSpec_ReportsItNicely(t *testing.T) {
	interfaceName := "WithParseErrors"
	pkgPath := "impl/test_data/panther"
	wantErr := &errs.InterfaceNotFoundError{}
	fileNameWithError := "with_parse_errors.go"

	pkg, err := buildPackage(pkgPath)
	if err != nil {
		t.Errorf("interfaceTypeSpec(...) failed precondition: could not load package with path %q", pkgPath)
	}
	_, gotErr := interfaceTypeSpec(interfaceName, pkg)
	if gotErr == nil {
		t.Errorf(`interfaceTypeSpec(%q, %q): wanted error type "%T", got "%T"`,
			interfaceName, pkgPath, wantErr, gotErr)
	}

	if !strings.Contains(gotErr.Error(), "could not be parsed") {
		t.Errorf(`interfaceTypeSpec(%q, %q): did not report some files could not be parsed`,
			interfaceName, pkgPath)
	}

	if !strings.Contains(gotErr.Error(), fileNameWithError) {
		t.Errorf(`interfaceTypeSpec(%q, %q): did not report %q could not be parsed`,
			interfaceName, pkgPath, fileNameWithError)
	}
}

func TestBuildInterface(t *testing.T) {
	cases := []struct {
		interfacePath string
		wantInterface *Interface
		wantErr       error
	}{
		{
			"impl/test_data/panther.Clawable",
			NewInterface(
				[]Method{
					NewMethod("Hardness", []Parameter{}, []Parameter{NewParameter("", "int")}),
					NewMethod("Puncture", []Parameter{NewParameter("strength", "int")}, []Parameter{}),
				},
			),
			nil,
		},
		{
			"impl/test_data/panther.Scenario",
			NewInterface(
				[]Method{
					NewMethod(
						"TwoTogether",
						[]Parameter{NewParameter("i", "int"), NewParameter("j", "int")},
						[]Parameter{NewParameter("a", "bool"), NewParameter("b", "bool")}),
					NewMethod(
						"TwoSeparate",
						[]Parameter{NewParameter("i", "int"), NewParameter("j", "int")},
						[]Parameter{NewParameter("a", "bool"), NewParameter("b", "bool")}),
				},
			),
			nil,
		},
		{
			"sort.Interface",
			NewInterface(
				[]Method{
					NewMethod(
						"Len",
						[]Parameter{},
						[]Parameter{NewParameter("", "int")}),
					NewMethod(
						"Less",
						[]Parameter{NewParameter("i", "int"), NewParameter("j", "int")},
						[]Parameter{NewParameter("", "bool")}),
					NewMethod(
						"Swap",
						[]Parameter{NewParameter("i", "int"), NewParameter("j", "int")},
						[]Parameter{}),
				},
			),
			nil,
		},
		{
			"io.ReadWriter", // embedded interfaces
			NewInterface(
				[]Method{
					NewMethod(
						"Read",
						[]Parameter{NewParameter("p", "byte[]")},
						[]Parameter{NewParameter("n", "int"), NewParameter("err", "error")}),
					NewMethod(
						"Write",
						[]Parameter{NewParameter("p", "byte[]")},
						[]Parameter{NewParameter("n", "int"), NewParameter("err", "error")}),
				},
			),
			nil,
		},
	}

	for _, c := range cases {
		gotInterface, gotErr := BuildInterface(c.interfacePath)
		if reflect.TypeOf(gotErr) != reflect.TypeOf(c.wantErr) {
			t.Errorf(`buildInterface(%q): wanted error type "%T", got "%T": %q`,
				c.interfacePath, c.wantErr, gotErr, gotErr.Error())
		} else if c.wantErr != nil {
			continue // The error match passed. Nothing more to test.
		} else if !reflect.DeepEqual(gotInterface, c.wantInterface) {
			t.Errorf("buildInterface(%q)\ngot:\t%+v\nwanted:\t%+v",
				c.interfacePath, gotInterface, c.wantInterface)
		}
	}
}
