package helpers

import (
	"dungeons-and-diagrams/board"
	"fmt"

	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/log"
)

// Update table rows to reflect board state
func UpdateTable(b board.Board, tb *table.Table) {
	log.Info("Updating table")
	tb.ClearRows()
	var header [board.BOARD_DIM + 1]string
	header[0] = " "
	for x := range board.BOARD_DIM {
		t := b.ColTotals[x]
		if t != nil {
			header[x+1] = fmt.Sprint(*t)
		} else {
			header[x+1] = "*"
		}
	}
	tb.Row(header[:]...)

	for y := range board.BOARD_DIM {
		var row [board.BOARD_DIM + 1]string
		tr := b.RowTotals[y]
		if tr != nil {
			row[0] = fmt.Sprint(*tr)
		} else {
			row[0] = "*"
		}
		for x := range board.BOARD_DIM {
			row[x+1] = (*board.Address(x, y, &b.Cells)).String()
		}
		tb.Row(row[:]...)
	}
}
