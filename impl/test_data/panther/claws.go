package panther

import (
	"io"

	"ultimatesoftware.com/accountstore/models"
	"ultimatesoftware.com/accountstore/utils"
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

type WithMap interface {
	TakeGiveMap(theMap map[string]io.Reader) map[int]string
}

type WithChannel interface {
	TakeGiveChannel(theChannel chan int) chan string
}

type WithEllipsis interface {
	TakeEllipsis(several ...int) int
}

type WithStars interface {
	GetAccounts(tenantId string, opts *utils.QueryOpts) ([]models.AccountSummary, error)
	GetTenants(tenantId string, filters *utils.QueryOpts, recursive bool) ([]models.TenantSummary, error)
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
