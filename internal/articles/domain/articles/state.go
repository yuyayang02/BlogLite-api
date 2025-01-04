package articles

type State struct {
	s int
}

var (
	PublicState  = State{1}
	PrivateState = State{0}
)

var validStates = []State{
	PublicState,
	PrivateState,
}

func (s State) ToInt() int {
	return s.s
}
