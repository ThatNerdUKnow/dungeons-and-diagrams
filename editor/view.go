package editor

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var BaseBorder = lipgloss.NewStyle()
var red = lipgloss.Color("#FF0000")
var SatBorder = BaseBorder.Foreground(purple)
var UnsatBorder = BaseBorder.Foreground(red)
var HeaderStyle = lipgloss.NewStyle().Bold(true).Align(lipgloss.Right)
var CellStyle = lipgloss.NewStyle()
var SelectedStyle = func(s lipgloss.Style) lipgloss.Style {
	return s.Blink(true).Reverse(true)
}

func (m Model) View() string {
	x, y := m.cursor.Coords()
	cursor_coords := [2]int{x, y}
	m.table.StyleFunc(func(row, col int) lipgloss.Style {
		var style lipgloss.Style
		if row == 0 || col == 0 {
			style = HeaderStyle
		} else {
			style = CellStyle
		}

		coords := [2]int{col, row}
		if coords == cursor_coords {
			return SelectedStyle(style)
		}
		return style
	})

	sat, err := m.Board.Check()
	if err != nil {
		log.Fatalf("%v", err)
	}

	if sat {
		m.table.BorderStyle(SatBorder)
	} else {
		m.table.BorderStyle(UnsatBorder)
	}

	tr := m.table.Render()
	h := m.help.View(m.keymap)
	w := lipgloss.Width(tr)
	boardTitle := HeaderStyle.Width(w).Render(m.Name)
	// saving this for my back pocket when I implement parsing for last call BBS level data
	//editorTitle := figure.NewFigure("Dungeons & Diagrams", "cosmic", true)

	return lipgloss.JoinVertical(lipgloss.Left, boardTitle, tr, h)
}
