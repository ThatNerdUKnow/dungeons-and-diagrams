package solved

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type ReturnToEditor struct{}

func (m SolvedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.edit):
			{
				f := func() tea.Msg {
					return ReturnToEditor{}
				}
				return m, f
			}
		}
	}
	return m, nil
}
