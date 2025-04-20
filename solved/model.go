package solved

import (
	"dungeons-and-diagrams/board"
	"dungeons-and-diagrams/helpers"
	"dungeons-and-diagrams/style"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type SolvedModel struct {
	board.Board
	table  *table.Table
	help   help.Model
	keymap KeyMap
}

func New(b board.Board) SolvedModel {
	t := table.New()
	t.BorderColumn(false)
	t.Border(lipgloss.DoubleBorder())
	h := help.New()

	h.Styles.ShortKey = style.KeyStyle
	h.Styles.FullKey = style.KeyStyle

	k := NewKeyMap()
	helpers.UpdateTable(b, t)
	return SolvedModel{Board: b, table: t, help: h, keymap: k}
}

func (m SolvedModel) Init() tea.Cmd {
	return nil
}
