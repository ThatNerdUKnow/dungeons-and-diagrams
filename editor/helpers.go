package editor

import (
	"dungeons-and-diagrams/board"
	"fmt"

	"github.com/charmbracelet/log"
)

// Update table rows to reflect board state
func (m *Model) UpdateTable() {
	log.Info("Updating table")
	m.table.ClearRows()
	var header [board.BOARD_DIM + 1]string
	header[0] = " "
	for x := range board.BOARD_DIM {
		t := m.ColTotals[x]
		if t != nil {
			header[x+1] = fmt.Sprint(*t)
		} else {
			header[x+1] = "*"
		}
	}
	m.table.Row(header[:]...)

	for y := range board.BOARD_DIM {
		var row [board.BOARD_DIM + 1]string
		tr := m.RowTotals[y]
		if tr != nil {
			row[0] = fmt.Sprint(*tr)
		} else {
			row[0] = "*"
		}
		for x := range board.BOARD_DIM {
			row[x+1] = (*board.Address(x, y, &m.Cells)).String()
		}
		m.table.Row(row[:]...)
	}
}
