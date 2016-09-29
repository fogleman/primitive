package primitive

import "math/rand"

type Scorer interface {
	Score(shape Shape) float64
}

type State struct {
	Worker      *Worker
	Shape       Shape
	Alpha       int
	MutateAlpha bool
	score       float64
	rnd         *rand.Rand
}

func NewState(worker *Worker, shape Shape, alpha int) *State {
	var mutateAlpha bool
	if alpha == 0 {
		alpha = 128
		mutateAlpha = true
	}
	return &State{model, buffer, alpha, shape, mutateAlpha, -1, rnd}
}

func (state *State) Energy() float64 {
	if state.score < 0 {
		state.score = state.Model.Energy(state.Alpha, state.Shape, state.Buffer)
	}
	return state.score
}

func (state *State) DoMove() interface{} {
	oldState := state.Copy()
	state.Shape.Mutate(state.rnd)
	if state.MutateAlpha {
		state.Alpha = clampInt(state.Alpha+state.rnd.Intn(21)-10, 1, 255)
	}
	state.score = -1
	return oldState
}

func (state *State) UndoMove(undo interface{}) {
	oldState := undo.(*State)
	state.Shape = oldState.Shape
	state.score = oldState.score
}

func (state *State) Copy() Annealable {
	return &State{
		state.Model, state.Buffer, state.Alpha, state.Shape.Copy(), state.MutateAlpha, state.score, state.rnd}
}
