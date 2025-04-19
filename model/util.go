package model

import (
	"dungeons-and-diagrams/editor"
	"dungeons-and-diagrams/solved"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func (m Model) getCurrentModel() tea.Model {
	log.Debugf("Current model is %s", m.modeString())
	var model tea.Model
	if m.mode == editmode {
		model = m.edit
	} else if m.mode == solvedmode {
		model = m.solved
	}
	return model
}

func (m *Model) updateCurrentModel(md tea.Model) {
	switch m.mode {
	case editmode:
		if md, ok := md.(editor.Model); ok {
			m.edit = md
		} else {
			log.Warn("Unexpected type when assigning to edit")
		}
	case solvedmode:
		if md, ok := md.(*solved.SolvedModel); ok {
			m.solved = md
		} else {
			log.Warn("Unexpected type when assigning to solved")
		}
	}
}

func (m Model) modeString() string {
	switch m.mode {
	case editmode:
		return "edit"
	case solvedmode:
		return "solved"
	}
	panic("Unreachable")
}
