package main

import (
	"fmt"
	"sync/atomic"

	tea "github.com/charmbracelet/bubbletea"
)

type Game struct {
	update  int
	model   *Model
	program *tea.Program
	view    *View
	frame   *Buffer
	logger  *LoggerComponent
	title   *TitleComponent
	isDirty atomic.Bool
}

func NewGame(width int) *Game {
	g := &Game{}
	g.frame = NewBuffer(width-2, 25)
	g.model = NewModel()
	g.program = tea.NewProgram(g)
	g.view = NewView(g.model)
	g.logger = NewLoggerComponent(g, width-2)
	g.title = NewTitleComponent(width-2, "Game 734")

	return g
}

func (g *Game) Dirty() {
	if g.isDirty.CompareAndSwap(true, true) {
		return
	}
	go g.program.Send(0)
}

func (g *Game) Init() tea.Cmd {

	return nil
}

func (g *Game) Run() (tea.Model, error) {
	return g.program.Run()
}

func (g *Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	g.update++
	g.isDirty.Store(false)

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		g.logger.Log(fmt.Sprintf("%v", msg))

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q", "esc":
			return g, tea.Quit

		case "up":
		case "down":
		case "left":
		case "right":

		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return g, nil
}

func (g *Game) View() string {

	g.frame.Clear()
	g.frame.DrawBox(0, 0, g.frame.width-1, g.frame.height-1)
	g.frame.WriteBuffer(1, 1, g.title.Render())
	g.frame.WriteBuffer(1, 2, g.logger.Render())

	return g.frame.String()
}
