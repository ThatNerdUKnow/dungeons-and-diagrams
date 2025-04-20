package style

import "github.com/charmbracelet/lipgloss"

var (
	Purple = lipgloss.Color("99")
)

var keyBorderStyle = lipgloss.Border{
	Left:  "[",
	Right: "]"}

var KeyStyle = lipgloss.NewStyle().
	Border(keyBorderStyle, false, true).
	Foreground(Purple).
	BorderForeground(Purple)
