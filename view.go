package main

type View struct {
	model *Model
}

func NewView(m *Model) *View {
	v := &View{model: m}

	return v
}

func (g *View) View() string {
	return ""
}
