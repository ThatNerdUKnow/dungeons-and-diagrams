package main

import (
	"fmt"

	"dungeons-and-diagrams/board"

	"github.com/aclements/go-z3/z3"
)

func main() {
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
	brd.SetCell(2, 3, board.Wall)
	brd.SetCell(5, 5, board.Treasure)
	brd.SetCell(7, 6, board.Monster)
	brd.SetCell(7, 5, board.Wall)
	brd.SetCell(7, 7, board.Wall)
	fmt.Println(brd)

}
