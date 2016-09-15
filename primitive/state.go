package primitive

import "image"

type State struct {
	Model  *Model
	Buffer *image.RGBA
	Shape  Shape
}

func NewState(model *Model, buffer *image.RGBA, shape Shape) *State {
	return &State{model, buffer, shape}
}

func (state *State) Energy() float64 {
	return state.Model.Energy(state.Shape, state.Buffer)
}

func (state *State) DoMove() interface{} {
	oldShape := state.Shape.Copy()
	state.Shape.Mutate()
	return oldShape
}

func (state *State) UndoMove(undo interface{}) {
	state.Shape = undo.(Shape)
}

func (state *State) Copy() Annealable {
	return &State{state.Model, state.Buffer, state.Shape.Copy()}
}
