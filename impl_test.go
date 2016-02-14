package impl

import (
	"impl/errs"

	"bytes"
	"reflect"
	"testing"
)

func TestImpl(t *testing.T) {
	cases := []struct {
		interfacePath string
		receiver      string
		wantErr       error
		wantSource    string
	}{
		{
			"impl/test_data/panther.Clawable",
			"r *Repo",
			nil,
			`func (r *Repo) Hardness() int {
	panic("TODO: implement this method")
}

func (r *Repo) Puncture(strength int) {
	panic("TODO: implement this method")
}

`,
		},
		{
			"sort.Interface",
			"m *MusicList",
			nil,
			`func (m *MusicList) Len() int {
	panic("TODO: implement this method")
}

func (m *MusicList) Less(i int, j int) bool {
	panic("TODO: implement this method")
}

func (m *MusicList) Swap(i int, j int) {
	panic("TODO: implement this method")
}

`,
		},
		{
			"io.ReadWriter",
			"",
			nil,
			`func () Read(p []byte) (n int, err error) {
	panic("TODO: implement this method")
}

func () Write(p []byte) (n int, err error) {
	panic("TODO: implement this method")
}

`,
		},
		{
			"net/http.Handler",
			"s Server",
			nil,
			`func (s Server) ServeHTTP(ResponseWriter, *Request) {
	panic("TODO: implement this method")
}

`,
		},
		{
			"encoding/json.Marshaler",
			"b *Banana",
			nil,
			`func (b *Banana) MarshalJSON() ([]byte, error) {
	panic("TODO: implement this method")
}

`,
		},
		{
			"os.FileInfo",
			"src Source",
			nil,
			`func (src Source) Name() string {
	panic("TODO: implement this method")
}

func (src Source) Size() int64 {
	panic("TODO: implement this method")
}

func (src Source) Mode() FileMode {
	panic("TODO: implement this method")
}

func (src Source) ModTime() time.Time {
	panic("TODO: implement this method")
}

func (src Source) IsDir() bool {
	panic("TODO: implement this method")
}

func (src Source) Sys() interface{} {
	panic("TODO: implement this method")
}

`,
		},
		{
			"io.NonExistent",
			"f *os.File",
			&errs.InterfaceNotFoundError{},
			"",
		},
		{
			"",
			"f *os.File",
			&errs.InvalidImportFormatError{},
			"",
		},
	}

	for _, c := range cases {
		var w bytes.Buffer
		gotErr := Impl(c.interfacePath, c.receiver, &w)
		if reflect.TypeOf(gotErr) != reflect.TypeOf(c.wantErr) {
			t.Errorf("Impl(%q, %q, <writer>) == %T, wanted error: %T.\n%q",
				c.interfacePath, c.receiver, gotErr, c.wantErr, gotErr)
		} else if c.wantErr != nil {
			continue // got the error we wanted
		} else if gotSrc := w.String(); c.wantSource != gotSrc {
			t.Errorf("Impl(%q, %q, <writer>) == \n\"%s\"\n, wanted: \n\"%s\"\n",
				c.interfacePath, c.receiver, gotSrc, c.wantSource)
		}
	}
}
