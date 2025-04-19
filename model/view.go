package model

func (m Model) View() string {
	return m.getCurrentModel().View()
}
