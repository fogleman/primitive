package primitive

type State struct {
	Worker      *Worker
	Shape       Shape
	Alpha       int
	MutateAlpha bool
	Score       float64
}

func NewState(worker *Worker, shape Shape, alpha int) *State {
	var mutateAlpha bool
	if alpha == 0 {
		alpha = 128
		mutateAlpha = true
	}
	return &State{worker, shape, alpha, mutateAlpha, -1}
}

func (state *State) Energy() float64 {
	if state.Score < 0 {
		state.Score = state.Worker.Energy(state.Shape, state.Alpha)
	}
	return state.Score
}

func (state *State) DoMove() interface{} {
	rnd := state.Worker.Rnd
	oldState := state.Copy()
	state.Shape.Mutate(state.Worker)
	if state.MutateAlpha {
		state.Alpha = clampInt(state.Alpha+rnd.Intn(21)-10, 1, 255)
	}
	state.Score = -1
	return oldState
}

func (state *State) UndoMove(undo interface{}) {
	oldState := undo.(*State)
	state.Shape = oldState.Shape
	state.Alpha = oldState.Alpha
	state.Score = oldState.Score
}

func (state *State) Copy() Annealable {
	return &State{
		state.Worker, state.Shape.Copy(), state.Alpha, state.MutateAlpha, state.Score}
}
