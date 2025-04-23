package editor

import (
	"dungeons-and-diagrams/style"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	x, y := m.cursor.Coords()
	cursor_coords := [2]int{x, y}
	m.table.StyleFunc(func(row, col int) lipgloss.Style {
		var s lipgloss.Style
		if row == 0 || col == 0 {
			s = style.HeaderStyle
		} else {
			s = style.CellStyle
		}

		coords := [2]int{col, row}
		if coords == cursor_coords {
			return style.SelectedStyle(s)
		}
		return s
	})

	if m.sat {
		m.table.BorderStyle(style.SatBorder)
	} else {
		m.table.BorderStyle(style.UnsatBorder)
	}

	tr := m.table.Render()
	h := m.help.View(m.keymap)
	w := lipgloss.Width(tr)
	boardTitle := style.HeaderStyle.Width(w).Render(m.Name)
	// saving this for my back pocket when I implement parsing for last call BBS level data
	//editorTitle := figure.NewFigure("Dungeons & Diagrams", "cosmic", true)

	return lipgloss.JoinVertical(lipgloss.Left, boardTitle, tr, h)
}
