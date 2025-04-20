package editor

import (
	"dungeons-and-diagrams/board"
	"dungeons-and-diagrams/helpers"
	"dungeons-and-diagrams/style"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type Model struct {
	board.Board
	cursor helpers.Cursor2D
	table  *table.Table
	keymap KeyMap
	help   help.Model
}

func New(b board.Board) Model {
	b.Build()

	m := NewKeyMap()
	h := help.New()
	c := helpers.NewCursor2D(board.BOARD_DIM+1, board.BOARD_DIM+1)
	t := table.New()

	t.BorderColumn(false)
	t.Border(lipgloss.DoubleBorder())
	keystyle := style.KeyStyle
	h.Styles.FullKey = keystyle
	h.Styles.ShortKey = keystyle
	model := Model{Board: b, cursor: c, table: t, keymap: m, help: h}

	model.cursor.X.Inc()
	model.cursor.Y.Inc()
	model.UpdateTable()
	model.UpdateKeymap()
	model.help.ShowAll = true
	return model
}

func Default() Model {
	b := board.NewBoard("")
	return New(b)
}

func (m Model) Init() tea.Cmd {
	return nil
}
