package solved

import "github.com/charmbracelet/lipgloss"

func (m SolvedModel) View() string {
	tr := m.table.Render()
	w := lipgloss.Width(tr)
	name := lipgloss.NewStyle().Width(w).Render(m.Board.Name)
	return lipgloss.JoinVertical(lipgloss.Left, name, tr, m.help.View(m.keymap))
}
