package editor

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Build struct{}
type Solve struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			{
				m.cursor.Y.Dec()
			}
		case "down":
			{
				m.cursor.Y.Inc()
			}
		case "left":
			{
				m.cursor.X.Dec()
			}
		case "right":
			{
				m.cursor.X.Inc()
			}
		case "enter":
			{
			}
		case "ctrl+c", "q":
			{
				return m, tea.Quit
			}
		}
	case Build:
		{

		}
	case Solve:
		{
			m.Solve()
		}
	}

	return m, nil
}
