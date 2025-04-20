package style

import "github.com/charmbracelet/lipgloss"

var BaseBorder = lipgloss.NewStyle()
var red = lipgloss.Color("#FF0000")
var SatBorder = BaseBorder.Foreground(Purple)
var UnsatBorder = BaseBorder.Foreground(red)
var HeaderStyle = lipgloss.NewStyle().Bold(true).Align(lipgloss.Right)
var CellStyle = lipgloss.NewStyle()
var SelectedStyle = func(s lipgloss.Style) lipgloss.Style {
	return s.Blink(true).Reverse(true)
}
