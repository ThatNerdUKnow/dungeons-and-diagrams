package main

import (
	"dungeons-and-diagrams/board"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

func main() {

	// Set log output to a file instead of stdout
	f, err := os.Create("run.log")
	if err != nil {
		log.Fatal("could not create log file", "error", err)
	}
	defer f.Close()

	log.SetOutput(f)

	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
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
		log.Errorf("Could not solve board. %v", err)
	}
	fmt.Println(nb)

	brd = board.NewBoard("Graveyard of the vernal king")
	brd.SetColTotals([8]int{4, 2, 5, 0, 6, 2, 4, 2})
	brd.SetRowTotals([8]int{5, 2, 2, 1, 5, 3, 2, 5})

	brd.SetCell(0, 6, board.Monster)
	brd.SetCell(3, 7, board.Monster)
	brd.SetCell(5, 7, board.Monster)
	brd.SetCell(7, 7, board.Monster)
	brd.SetCell(6, 2, board.Treasure)
	brd.SetCell(2, 2, board.Monster)
	fmt.Println(brd)

	nb, err = brd.Solve()
	if nb == nil {
		log.Errorf("Could not solve board. %v", err)
	}
	fmt.Println(nb)

	brd = board.NewBoard("This is a test")
	brd.SetRowTotals([8]int{5, 2, 2, 1, 5, 3, 2, 5})
	fmt.Println(brd)

	nb, err = brd.Solve()
	if nb == nil {
		log.Errorf("Could not solve board. %v", err)
	}
	fmt.Println(nb)
}
