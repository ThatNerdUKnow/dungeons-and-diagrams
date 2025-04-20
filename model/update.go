package model

import (
	"dungeons-and-diagrams/editor"
	"dungeons-and-diagrams/solved"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cm := m.getCurrentModel()
	out, cmd := cm.Update(msg)
	if cmd != nil {
		log.Warn("Command is not nil")
	}
	m.updateCurrentModel(out)

	switch msg := msg.(type) {
	case editor.Solved:
		{
			m.mode = solvedmode
			sm := solved.New(msg.Board)
			m.solved = &sm
		}
	case solved.ReturnToEditor:
		{
			m.mode = editmode
		}
	}
	return m, tea.Batch(cmd, nil)
}
