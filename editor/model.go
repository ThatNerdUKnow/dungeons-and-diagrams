package editor

import (
	"dungeons-and-diagrams/board"
	"dungeons-and-diagrams/helpers"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/table"
)

type Model struct {
	board.Board
	cursor helpers.Cursor2D
	table  *table.Table
	keymap KeyMap
	help   help.Model
}

func New() Model {
	b := board.NewBoard("New Dungeon")
	c := helpers.NewCursor2D(board.BOARD_DIM+1, board.BOARD_DIM+1)
	t := table.New()
	m := NewKeyMap()
	h := help.New()

	keystyle := keyStyle
	h.Styles.FullKey = keystyle
	h.Styles.ShortKey = keystyle
	model := Model{Board: b, cursor: c, table: t, keymap: m, help: h}
	model.UpdateTable()
	model.cursor.X.Inc()
	model.cursor.Y.Inc()
	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}
