package solved

import (
	"dungeons-and-diagrams/board"
	"dungeons-and-diagrams/helpers"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
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
	h := help.New()
	k := NewKeyMap()
	helpers.UpdateTable(b, t)
	return SolvedModel{Board: b, table: t, help: h, keymap: k}
}

func (m SolvedModel) Init() tea.Cmd {
	return nil
}
