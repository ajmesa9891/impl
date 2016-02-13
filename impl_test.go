package impl

import (
	"bytes"
	"reflect"
	"testing"
)

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
						[]Parameter{NewParameter("p", "[]byte")},
						[]Parameter{NewParameter("n", "int"), NewParameter("err", "error")}),
					NewMethod(
						"Write",
						[]Parameter{NewParameter("p", "[]byte")},
						[]Parameter{NewParameter("n", "int"), NewParameter("err", "error")}),
				},
			),
			nil,
		},
		// TODO: make this case work! (embedding and interface of a different packages)
		// {
		// 	// embedded with an interface from another package
		// 	"impl/test_data/panther.ExternalEmbedded",
		// 	NewInterface(
		// 		[]Method{
		// 			NewMethod(
		// 				"Read",
		// 				[]Parameter{NewParameter("p", "byte[]")},
		// 				[]Parameter{NewParameter("n", "int"), NewParameter("err", "error")}),
		// 			NewMethod(
		// 				"Write",
		// 				[]Parameter{NewParameter("p", "byte[]")},
		// 				[]Parameter{NewParameter("n", "int"), NewParameter("err", "error")}),
		// 		},
		// 	),
		// 	nil,
		// },
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

func TestRenderInterface(t *testing.T) {
	cases := []struct {
		iface      *Interface
		receiver   string
		wantErr    error
		wantSource string
	}{
		{
			NewInterface(
				[]Method{
					NewMethod(
						"Read",
						[]Parameter{NewParameter("p", "[]byte")},
						[]Parameter{NewParameter("n", "int"), NewParameter("err", "error")}),
				},
			),
			"r *Repo",
			nil,
			`func (r *Repo) Read(p []byte) (n int, err error) {
	panic("TODO: implement this method")
}

`,
		},
		{
			NewInterface(
				[]Method{
					NewMethod(
						"Simplest",
						[]Parameter{},
						[]Parameter{}),
					NewMethod(
						"Read",
						[]Parameter{NewParameter("p", "[]byte")},
						[]Parameter{NewParameter("n", "int"), NewParameter("err", "error")}),
					NewMethod(
						"Unnamed",
						[]Parameter{NewParameter("p", "[]byte"), NewParameter("r", "io.Reader")},
						[]Parameter{NewParameter("", "int"), NewParameter("", "error")}),
				},
			),
			"rec receiver",
			nil,
			`func (rec receiver) Simplest() {
	panic("TODO: implement this method")
}

func (rec receiver) Read(p []byte) (n int, err error) {
	panic("TODO: implement this method")
}

func (rec receiver) Unnamed(p []byte, r io.Reader) (int, error) {
	panic("TODO: implement this method")
}

`,
		},
	}

	for _, c := range cases {
		var w bytes.Buffer
		gotErr := RenderInterface(c.iface, c.receiver, &w)
		if reflect.TypeOf(gotErr) != reflect.TypeOf(c.wantErr) {
			t.Errorf("RenderInterface(<interface>, %s, <writer>) == %T, wanted error: %T.\n%q", c.receiver, gotErr, c.wantErr, gotErr)
		} else if c.wantErr != nil {
			continue // got the error we wanted
		} else if gotSrc := w.String(); c.wantSource != gotSrc {
			t.Errorf("RenderInterface(<interface>, %s, <writer>) == \n\"%s\"\n, wanted: \n\"%s\"\n",
				c.receiver, gotSrc, c.wantSource)
		}
	}
}
