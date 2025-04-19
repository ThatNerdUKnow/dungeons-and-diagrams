package solved

import (
	"dungeons-and-diagrams/board"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/table"
)

type SolvedModel struct {
	board.Board
	table *table.Table
	help  help.Model
}

func New(b board.Board) SolvedModel {
	t := table.New()
	h := help.New()
	return SolvedModel{Board: b, table: t, help: h}
}

func (m SolvedModel) Init() tea.Cmd {
	return nil
}
