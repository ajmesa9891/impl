package panther

type Clawable interface {
	Hardness() int
	Puncture(strength int)
}

type Scenario interface {
	TwoTogether(i, j int) (a, b bool)
	TwoSeparate(i, j int) (a bool, b bool)
}
