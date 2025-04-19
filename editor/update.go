package editor

import (
	"dungeons-and-diagrams/board"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Up):
			{
				m.cursor.Y.Dec()
				m.UpdateKeymap()
			}
		case key.Matches(msg, m.keymap.Down):
			{
				m.cursor.Y.Inc()
				m.UpdateKeymap()
			}
		case key.Matches(msg, m.keymap.Left):
			{
				m.cursor.X.Dec()
				m.UpdateKeymap()
			}
		case key.Matches(msg, m.keymap.Right):
			{
				m.cursor.X.Inc()
				m.UpdateKeymap()
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
				i, err := strconv.Atoi(msg.String())
				if err != nil {
					log.Fatal(err)
				}
				m.SetHeader(i)
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
				_, err := m.Solve()
				if err != nil {

				}
			}
		}

	}
	return m, nil
}

func (m *Model) SetCell(cell board.Cell) {
	x, y := m.cursor.CoordsOffset(-1, -1)
	// setcell will panic if x and y are out of bounds which is desired behavior
	m.Board.SetCell(x, y, cell)
	m.UpdateTable()
}

func (m *Model) SetHeader(i int) {
	logger := log.With("i", i)
	if i < 0 || i > 8 {
		logger.Fatal("invalid header value. Appropriate range is 0-8")
	}

	x, y := m.cursor.Coords()
	logger = logger.With("x", x, "y", y)
	if x == 0 && y == 0 {
		logger.Fatal("Cursor is not pointing at headers")
	} else if x == 0 {
		logger.Info("Setting row total")
		m.Board.SetRowTotal(y - 1)(i)
	} else if y == 0 {
		logger.Info("Setting column total")
		m.Board.SetColTotal(x - 1)(i)
	} else {
		logger.Fatal("Cursor is not pointing at headers")
	}
	m.UpdateTable()
}
