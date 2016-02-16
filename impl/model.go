package impl

type Interface struct {
	Methods []Method
}

func NewInterface(m []Method) *Interface {
	return &Interface{m}
}

type Method struct {
	Name string
	In   []Parameter
	Out  []Parameter
}

func NewMethod(name string, in []Parameter, out []Parameter) Method {
	return Method{name, in, out}
}

type Parameter struct {
	Name string
	Type string
}

// NewParameter creates a new parameter with the given name and type.
// An empty name creates an unnamed parameter, meant to be returned.
func NewParameter(name, typeName string) Parameter {
	return Parameter{name, typeName}
}
