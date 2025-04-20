package solved

import "github.com/charmbracelet/lipgloss"

func (m SolvedModel) View() string {
	tr := m.table.Render()
	return lipgloss.JoinVertical(lipgloss.Left, tr, m.help.View(m.keymap))
}
