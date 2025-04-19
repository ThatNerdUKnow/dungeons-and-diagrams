package editor

import "github.com/charmbracelet/lipgloss"

var keyBorderStyle = lipgloss.Border{
	Left:  "[",
	Right: "]"}

var purple = lipgloss.Color("99")

var keyStyle = lipgloss.NewStyle().
	Border(keyBorderStyle, false, true).
	Foreground(purple).
	BorderForeground(purple)
