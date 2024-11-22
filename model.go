package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	program *tea.Program
}

func NewModel() *Model {
	m := &Model{}
	return m
}

func (m *Model) Tick() {
	for {
		time.Sleep(time.Second / 30)
		m.program.Send(0)
	}
}

// func (m *Model) View() string {
// 	// return m.view.View() + m.logs.View()
// 	return ""
// }
