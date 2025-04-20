package solved

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	edit key.Binding
}

func NewKeyMap() KeyMap {
	edit := key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "return to editor"))
	return KeyMap{edit: edit}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.edit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.edit}}
}
