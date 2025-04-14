package editor

import (
	"dungeons-and-diagrams/board"
	"fmt"

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

	m.table.ClearRows()
	cursor_coords := m.cursor.Coords()
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

	var header [board.BOARD_DIM + 1]string
	header[0] = " "
	for x := range board.BOARD_DIM {
		t := m.ColTotals[x]
		if t != nil {
			header[x+1] = fmt.Sprint(*t)
		} else {
			header[x+1] = "*"
		}
	}
	m.table.Row(header[:]...)

	for y := range board.BOARD_DIM {
		var row [board.BOARD_DIM + 1]string
		tr := m.RowTotals[y]
		if tr != nil {
			row[0] = fmt.Sprint(*tr)
		} else {
			row[0] = "*"
		}
		for x := range board.BOARD_DIM {
			row[x+1] = (*board.Address(x, y, &m.Cells)).String()
		}
		m.table.Row(row[:]...)
	}

	tr := m.table.Render()
	w := lipgloss.Width(tr)
	title := HeaderStyle.Width(w).Render(m.Name)
	return lipgloss.JoinVertical(lipgloss.Center, title, tr)
}
