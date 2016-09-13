package tri

type State struct {
	Model *Model
	Shape Shape
}

func NewState(model *Model, shape Shape) *State {
	return &State{model, shape}
}

func (state *State) Energy() float64 {
	lines := state.Shape.Rasterize()
	c := state.Model.computeColor(lines, 128)
	s := state.Model.computeScore(lines, c)
	return s
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
	return &State{state.Model, state.Shape.Copy()}
}
