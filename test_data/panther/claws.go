package panther

type Clawable interface {
	Hardness() int
	Puncture(strength int)
}
