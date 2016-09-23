package primitive

import "image"

type State struct {
	Model  *Model
	Buffer *image.RGBA
	Shape  Shape
	Score  float64
}

func NewState(model *Model, buffer *image.RGBA, shape Shape) *State {
	return &State{model, buffer, shape, -1}
}

func (state *State) Energy() float64 {
	if state.Score < 0 {
		state.Score = state.Model.Energy(state.Shape, state.Buffer)
	}
	return state.Score
}

func (state *State) DoMove() interface{} {
	oldState := state.Copy()
	state.Shape.Mutate()
	state.Score = -1
	return oldState
}

func (state *State) UndoMove(undo interface{}) {
	oldState := undo.(*State)
	state.Shape = oldState.Shape
	state.Score = oldState.Score
}

func (state *State) Copy() Annealable {
	return &State{state.Model, state.Buffer, state.Shape.Copy(), state.Score}
}
