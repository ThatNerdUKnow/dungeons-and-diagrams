package editor

import (
	"dungeons-and-diagrams/board"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type Build struct{}
type Solve struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Up):
			{
				m.cursor.Y.Dec()
			}
		case key.Matches(msg, m.keymap.Down):
			{
				m.cursor.Y.Inc()
			}
		case key.Matches(msg, m.keymap.Left):
			{
				m.cursor.X.Dec()
			}
		case key.Matches(msg, m.keymap.Right):
			{
				m.cursor.X.Inc()
			}
		case key.Matches(msg, m.keymap.Quit):
			{
				return m, tea.Quit
			}
		case key.Matches(msg, m.keymap.Delete):
			{

			}
		case key.Matches(msg, m.keymap.Numeric):
			{

			}
		case key.Matches(msg, m.keymap.Help):
			{
				m.help.ShowAll = !m.help.ShowAll
			}
		case key.Matches(msg, m.keymap.Space):
			{
				m.SetCell(board.Space)
			}
		case key.Matches(msg, m.keymap.Wall):
			{
				m.SetCell(board.Wall)
			}
		case key.Matches(msg, m.keymap.Monster):
			{
				m.SetCell(board.Monster)
			}
		case key.Matches(msg, m.keymap.Treasure):
			{
				m.SetCell(board.Treasure)
			}
		case key.Matches(msg, m.keymap.Solve):
			{

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

func (m *Model) SetCell(cell board.Cell) {
	x, y := m.cursor.CoordsOffset(-1, -1)
	// setcell will panic if x and y are out of bounds which is desired behavior
	m.Board.SetCell(x, y, cell)
}

func (m *Model) SetHeader(i int) {
	logger := log.With("i", i)
	if i < 0 || i > 8 {
		logger.Fatal("invalid header value. Appropriate range is 0-8")
	}

	x, y := m.cursor.Coords()

	if x == 0 && y == 0 {
		logger.With("x", x, "y", y).Fatal("Cursor is not pointing at headers")
	} else if x == 0 {
		m.Board.SetRowTotal(y)(i)
	} else if y == 0 {
		m.Board.SetColTotal(x)(i)
	} else {
		logger.With("x", x, "y", y).Fatal("Cursor is not pointing at headers")
	}
}
