package panther

import (
	"io"
)

type Clawable interface {
	Hardness() int
	Puncture(strength int)
}

type Scenario interface {
	TwoTogether(i, j int) (a, b bool)
	TwoSeparate(i, j int) (a bool, b bool)
}

type ExternalEmbedded interface {
	io.ReadWriter
}

type Type interface {
	Align() int

	FieldAlign() int

	Method(int) Method

	MethodByName(string) (Method, bool)

	NumMethod() int

	Name() string

	PkgPath() string

	Size() uintptr

	String() string

	Kind() Kind

	Implements(u Type) bool

	AssignableTo(u Type) bool

	ConvertibleTo(u Type) bool

	Comparable() bool

	Bits() int

	ChanDir() ChanDir

	IsVariadic() bool

	Elem() Type

	Field(i int) StructField

	FieldByIndex(index []int) StructField

	FieldByName(name string) (StructField, bool)

	FieldByNameFunc(match func(string) bool) (StructField, bool)

	In(i int) Type

	Key() Type

	Len() int

	NumField() int

	NumIn() int

	NumOut() int

	Out(i int) Type
}