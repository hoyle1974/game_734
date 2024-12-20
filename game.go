package main

import (
	"fmt"
	"sync/atomic"

	tea "github.com/charmbracelet/bubbletea"
)

type Game struct {
	views       int
	update      int
	model       *Model
	program     *tea.Program
	view        *View
	frame       *Buffer
	logger      *LoggerComponent
	playerstats *PlayerStatsComponent
	isDirty     atomic.Bool
	//display     *VirtualComponent
	display *RGBComponent
}

func NewGame(width, height int) *Game {
	g := &Game{}
	g.frame = NewBuffer(width-2, height)
	g.model = NewModel()
	g.program = tea.NewProgram(g)
	g.view = NewView(g.model)
	g.logger = NewLoggerComponent(g, width-2)
	g.playerstats = NewPlayerStatsComponent(g, height-g.logger.buffer.height-2)
	// g.display = NewVirtualComponent(500, 500, g)
	g.display = NewRGBComponent(120, 30)

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
	g.playerstats.SetStat("Updates", fmt.Sprintf("%d", g.update))
	g.playerstats.SetStat("Views", fmt.Sprintf("%d", g.views))

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
			// g.display.Move(0, -1)
		case "down":
			// g.display.Move(0, 1)
		case "left":
			// g.display.Move(-1, 0)
		case "right":
			// g.display.Move(1, 0)

		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return g, nil
}

func (g *Game) View() string {
	g.views++
	g.frame.Clear()
	g.frame.DrawBoxWithTitle(0, 0, g.frame.width-1, g.frame.height-1, "Game 734")
	g.frame.WriteBuffer(1, g.frame.height-g.logger.buffer.height-1, g.logger.Render())
	g.frame.WriteBuffer(g.frame.width-1-g.playerstats.buffer.width, 1, g.playerstats.Render())
	//g.frame.WriteBuffer(2, 2, g.display.Render(0, 0, g.frame.width-g.playerstats.buffer.width-4, g.frame.height-g.logger.buffer.height-3))
	g.frame.WriteBuffer(2, 2, g.display.Render())
	s := g.frame.String()
	l := len(s)
	return fmt.Sprintf("Data Length: %d\n%s", l, s)
}
