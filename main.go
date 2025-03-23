package main

import (
	"fmt"

	"dungeons-and-diagrams/board"

	"github.com/aclements/go-z3/z3"
	"github.com/charmbracelet/log"
)

func main() {
	log.SetLevel(log.DebugLevel)
	fmt.Println("Hello, World")
	cfg := z3.NewContextConfig()
	ctx := z3.NewContext(cfg)

	// Create symbolic integer variables
	a := ctx.IntConst("a")
	b := ctx.IntConst("b")
	c := ctx.FromInt(10, ctx.IntSort()).(z3.Int)
	// Create a solver
	slv := z3.NewSolver(ctx)

	// Add constraints
	slv.Assert(a.Add(b).Eq(c))
	sat, err := slv.Check()
	if !sat {
		fmt.Println("Not satisfiable: ", err)
		return
	}

	brd := board.NewBoard("Foo")
	brd.SetColTotals([8]int{1, 2, 3, 4, 5, 6, 7, 8})
	brd.SetCell(0, 0, board.Wall)
	brd.SetCell(0, 1, board.Wall)
	brd.SetCell(0, 2, board.Wall)
	brd.SetCell(0, 3, board.Wall)
	brd.SetCell(0, 4, board.Wall)
	brd.SetCell(0, 5, board.Wall)
	brd.SetCell(0, 6, board.Wall)
	brd.SetCell(0, 7, board.Wall)
	//brd.SetCell(5, 5, board.Treasure)
	//brd.SetCell(7, 6, board.Monster)
	//brd.SetCell(2, 2, board.Monster)
	//brd.SetCell(7, 5, board.Space)
	//brd.SetCell(7, 7, board.Space)
	fmt.Println(brd)

	nb, err := brd.Solve()
	if nb == nil {
		log.Errorf("Could not solve board. %s", err)
	}
	fmt.Println(nb)

}
