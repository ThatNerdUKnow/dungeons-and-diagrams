package editor

import (
	"github.com/charmbracelet/lipgloss"
)

var BaseBorder = lipgloss.NewStyle()
var green = lipgloss.Color("#00FF00")
var red = lipgloss.Color("#FF0000")
var SatBorder = BaseBorder.Foreground(green)
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
	tr := m.table.Render()
	h := m.help.View(m.keymap)
	w := lipgloss.Width(tr)
	title := HeaderStyle.Width(w).Render(m.Name)
	return lipgloss.JoinVertical(lipgloss.Left, title, tr, h)
}
