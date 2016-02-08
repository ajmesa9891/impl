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

// TODO: CONTINUE: by adding more cases for this test
func TestBuildInterface(t *testing.T) {
	cases := []struct {
		pkgPath       string
		interfaceName string
		wantInterface *Interface
		wantErr       error
	}{
		{
			"impl/test_data/panther",
			"Clawable",
			NewInterface(
				[]Method{
					NewMethod("Hardness", []Parameter{}),
					NewMethod("Puncture", []Parameter{NewParameter("strength", "int")}),
				},
			),
			nil,
		},
		// {"sort", "Interface", nil},

		// {"impl/test_data/panther", "UnexistentName", &errs.InterfaceNotFoundError{}},
		// {"impl/test_data/panther", "WithParseErrors", &errs.InterfaceNotFoundError{}},
	}

	for _, c := range cases {
		pkg, err := buildPackage(c.pkgPath)
		if err != nil {
			t.Errorf("buildInterface(...) failed precondition: could not load package with path %q", c.pkgPath)
		}
		ts, err := interfaceTypeSpec(c.interfaceName, pkg)
		if err != nil {
			t.Errorf("buildInterface(...) failed precondition: could load TypeSpec for \"%s.%s\"", c.pkgPath, c.interfaceName)
		}

		gotInterface, gotErr := buildInterface(ts)
		if reflect.TypeOf(gotErr) != reflect.TypeOf(c.wantErr) {
			t.Errorf(`buildInterface(<TypeSpec for %q>): wanted error type "%T", got "%T"`,
				c.interfaceName, c.pkgPath, c.wantErr, gotErr)
		} else if c.wantErr != nil {
			continue // The error match passed. Nothing more to test.
		} else if !reflect.DeepEqual(gotInterface, c.wantInterface) {
			t.Errorf("buildInterface(<TypeSpec for %q>)\ngot:\t%+v\nwanted:\t%+v",
				c.interfaceName, gotInterface, c.wantInterface)
		}
	}
}

// TODO: suing as repl. delete this test
// func TestExplore_DELETE_THIS_TEST(t *testing.T) {
// 	pkg, err := buildPackage("sort")
// 	if err != nil {
// 		t.Errorf("could not load package")
// 	}

// 	typeSpec, err := interfaceTypeSpec("Interface", pkg)
// 	if err != nil {
// 		t.Errorf("could not find interface")
// 	}

// 	fmt.Printf("type: %T\n", typeSpec.Type)
// 	interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
// 	fmt.Printf("is interface: %v\n", ok)
// 	if !ok {
// 		t.Errorf("not an interface")
// 	}

// 	fmt.Printf("methods.lists:\n")
// 	for _, m := range interfaceType.Methods.List {
// 		fmt.Printf("  field names: %s\n", m.Names)
// 		fmt.Printf("  field.Type: %T\n", m.Type)
// 		ft := m.Type.(*ast.FuncType)
// 		if len(ft.Params.List) > 0 {
// 			fmt.Printf("  field.Type.Params[0].Names: %v\n", ft.Params.List[0].Names)
// 		}

// 	}

// }
