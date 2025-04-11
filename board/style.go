package board

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	purple = lipgloss.Color("99")
	gray   = lipgloss.Color("245")
)

var t = table.New().
	Border(lipgloss.RoundedBorder()).
	BorderStyle(lipgloss.NewStyle().Foreground(purple))
