package editor

import (
	"dungeons-and-diagrams/board"
	"dungeons-and-diagrams/helpers"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	board.Board
	cursor helpers.Cursor2D
}

func (m Model) Init() tea.Cmd {
	return nil
}
