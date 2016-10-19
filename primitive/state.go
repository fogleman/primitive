package primitive

type State struct {
	Worker      *Worker
	Shape       Shape
	Alpha       int
	MutateAlpha bool
	Score       float64
}

// NewState uses a reference to a Worker type, a Shape type, and an
// alpha value to create a new State type. If 0, the new alpha value
// is set to roughly 50% opacity.
// Returns the address of the new State
func NewState(worker *Worker, shape Shape, alpha int) *State {
	var mutateAlpha bool
	if alpha == 0 {
		alpha = 128
		mutateAlpha = true
	}
	return &State{worker, shape, alpha, mutateAlpha, -1}
}

// Energy, implemented by State, returns the state's score, or calculates
// a new one by passing the state's shape type and alpha value to its Worker object.
// Return's the value of the state's score.
func (state *State) Energy() float64 {
	if state.Score < 0 {
		state.Score = state.Worker.Energy(state.Shape, state.Alpha)
	}
	return state.Score
}

// DoMove saves the current state, and then changes it by
// mutating the state's Shape object.
// Returns the previous state to allow reversal.
func (state *State) DoMove() interface{} {
	rnd := state.Worker.Rnd
	oldState := state.Copy()
	state.Shape.Mutate()
	if state.MutateAlpha {
		state.Alpha = clampInt(state.Alpha+rnd.Intn(21)-10, 1, 255)
	}
	state.Score = -1
	return oldState
}

// UndoMove returns the current state's Shape, Alpha, and Score to
// their values in the previous state.
func (state *State) UndoMove(undo interface{}) {
	oldState := undo.(*State)
	state.Shape = oldState.Shape
	state.Alpha = oldState.Alpha
	state.Score = oldState.Score
}

// Copy creates a new State type identical to the State type passed as argument.
// Returns a reference to the new copy State object.
func (state *State) Copy() Annealable {
	return &State{
		state.Worker, state.Shape.Copy(), state.Alpha, state.MutateAlpha, state.Score}
}
