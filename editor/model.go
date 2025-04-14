package editor

import (
	"dungeons-and-diagrams/board"
	"dungeons-and-diagrams/helpers"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/table"
)

type Model struct {
	board.Board
	cursor helpers.Cursor2D
	table  *table.Table
}

func New() Model {
	b := board.NewBoard("New Dungeon")
	c := helpers.NewCursor2D(board.BOARD_DIM+1, board.BOARD_DIM+1)
	t := table.New()

	return Model{Board: b, cursor: c, table: t}
}

func (m Model) Init() tea.Cmd {
	return nil
}
