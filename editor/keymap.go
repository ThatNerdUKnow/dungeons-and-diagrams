package editor

import (
	"dungeons-and-diagrams/board"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	/* Mutually exclusive with Cell bindings */
	Numeric key.Binding

	/* Always Enabled */
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Quit  key.Binding
	Help  key.Binding

	/* Contextually enabled */
	Delete key.Binding
	Solve  key.Binding

	/* Mutually Exclusive with numeric */
	Treasure key.Binding
	Monster  key.Binding
	Space    key.Binding
	Wall     key.Binding
}

func NewKeyMap() KeyMap {
	numbers := key.NewBinding(key.WithKeys(genNumbers()...), key.WithHelp("0-8", "Enter a number"), key.WithDisabled())
	up := key.NewBinding(key.WithKeys("up"), key.WithHelp("↑", "move up"))
	down := key.NewBinding(key.WithKeys("down"), key.WithHelp("↓", "move down"))
	left := key.NewBinding(key.WithKeys("left"), key.WithHelp("←", "move left"))
	right := key.NewBinding(key.WithKeys("right"), key.WithHelp("→", "move right"))
	space := key.NewBinding(key.WithKeys("1"), key.WithHelp("1", fmt.Sprintf("insert space (%s)", board.Space)))
	wall := key.NewBinding(key.WithKeys("2"), key.WithHelp("2", fmt.Sprintf("insert wall (%s)", board.Wall)))
	monster := key.NewBinding(key.WithKeys("3"), key.WithHelp("3", fmt.Sprintf("insert monster (%s)", board.Monster)))
	treasure := key.NewBinding(key.WithKeys("4"), key.WithHelp("4", fmt.Sprintf("insert treasure (%s)", board.Treasure)))
	quit := key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q", "quit"))
	solve := key.NewBinding(key.WithKeys("enter", "return"), key.WithHelp("↵", "solve the board"))
	delete := key.NewBinding(key.WithKeys("backspace", "delete"), key.WithHelp("backspace/delete", "delete selected element"))
	help := key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help"))
	return KeyMap{
		Numeric:  numbers,
		Up:       up,
		Down:     down,
		Left:     left,
		Right:    right,
		Space:    space,
		Wall:     wall,
		Monster:  monster,
		Treasure: treasure,
		Quit:     quit,
		Solve:    solve,
		Delete:   delete,
		Help:     help,
	}
}

func genNumbers() []string {
	var keys []string
	for n := '0'; n <= '8'; n++ {
		keys = append(keys, fmt.Sprint(n))
	}
	return keys
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Help}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Quit, k.Help}, {k.Up, k.Down, k.Left, k.Right}, {k.Numeric, k.Space, k.Wall, k.Monster, k.Treasure}, {k.Delete, k.Solve}}
}

func (m Model) UpdateKeymap() {

	x, y := m.cursor.Coords()
	cursorTopLeft := x == 0 && y == 0

	// true: cursor is inside the board cells - false: cursor is in the header section
	var cursorInBoard bool
	if x == 0 || y == 0 {
		cursorInBoard = false
	} else {
		cursorInBoard = true
	}

	m.keymap.Numeric.SetEnabled(!cursorInBoard && !cursorTopLeft)
	m.keymap.Space.SetEnabled(cursorInBoard && !cursorTopLeft)
	m.keymap.Wall.SetEnabled(cursorInBoard && !cursorTopLeft)
	m.keymap.Monster.SetEnabled(cursorInBoard && !cursorTopLeft)
	m.keymap.Treasure.SetEnabled(cursorInBoard && !cursorTopLeft)

}
