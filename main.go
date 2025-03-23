package main

import (
	"dungeons-and-diagrams/board"
	"fmt"
	"math/rand"

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
	perm := rand.Perm(8)
	var colTotals [8]int
	var rowTotals [8]int
	copy(colTotals[:], perm)
	perm = rand.Perm(8)
	copy(rowTotals[:], perm)
	brd.SetColTotals(colTotals)
	brd.SetRowTotals(rowTotals)
	fmt.Println(brd)

	nb, err := brd.Solve()
	if nb == nil {
		log.Errorf("Could not solve board. %s", err)
	}
	fmt.Println(nb)

}
