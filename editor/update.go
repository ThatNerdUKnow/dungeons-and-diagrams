package editor

import (
	"dungeons-and-diagrams/board"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type Solved struct{ board.Board }

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
				if m.InBoard() {
					m.SetCell(board.Unknown)
				} else if m.InHeaders() {
					m.SetHeader(nil)
				}
			}
		case key.Matches(msg, m.keymap.Numeric):
			{
				i, err := strconv.Atoi(msg.String())
				if err != nil {
					log.Fatal(err)
				}
				m.SetHeader(&i)
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
				nb, err := m.Solve()

				if nb != nil {
					f := func() tea.Msg {
						return Solved{Board: *nb}
					}
					log.Info("Sending solved update to model")
					return m, f
				} else {
					log.Fatalf("Could not solve board %v", err)
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
	m.Board.Build()
	m.UpdateTable()
	m.UpdateKeymap()
}

func (m *Model) SetHeader(i *int) {
	logger := log.With("i", i)
	if i != nil {
		if *i < 0 || *i > 8 {
			logger.Fatal("invalid header value. Appropriate range is 0-8")
		}
	}

	x, y := m.cursor.Coords()
	logger = logger.With("x", x, "y", y)
	if !m.InHeaders() {
		logger.Fatal("Cursor is not pointing at headers")
	}
	if x == 0 {
		logger.Info("Setting row total")
		m.Board.SetRowTotal(y - 1)(i)
	} else if y == 0 {
		logger.Info("Setting column total")
		m.Board.SetColTotal(x - 1)(i)
	}
	m.UpdateTable()
	m.Board.Build()
	m.UpdateKeymap()
}

// Is the cursor currently pointing inside the board
func (m Model) InBoard() bool {
	x, y := m.cursor.Coords()
	return x > 0 && y > 0
}

// is the cursor at 0,0
func (m Model) InCorner() bool {
	x, y := m.cursor.Coords()
	return x == 0 && y == 0
}

// is the cursor pointing at either row or column totals
func (m Model) InHeaders() bool {
	return !m.InBoard() && !m.InCorner()
}
