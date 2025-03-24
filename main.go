package main

import (
	"dungeons-and-diagrams/board"
	"fmt"

	"github.com/charmbracelet/log"
)

func main() {
	log.SetLevel(log.ErrorLevel)

	brd := board.NewBoard("Tenaxxus's Gullet")
	brd.SetColTotals([8]int{4, 4, 2, 6, 2, 3, 4, 7})
	brd.SetRowTotals([8]int{7, 3, 4, 1, 7, 1, 6, 3})
	brd.SetCell(1, 2, board.Treasure)
	brd.SetCell(5, 0, board.Monster)
	brd.SetCell(0, 5, board.Monster)
	brd.SetCell(2, 7, board.Monster)
	brd.SetCell(7, 7, board.Monster)

	fmt.Println(brd)

	nb, err := brd.Solve()
	if nb == nil {
		log.Errorf("Could not solve board. %s", err)
	}
	fmt.Println(nb)

}
