package board

import (
	"fmt"

	"github.com/aclements/go-z3/z3"
	"github.com/adam-lavrik/go-imath/ix"
	"github.com/charmbracelet/log"
)

var (
	adjacent = [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
)

func (b *Board) checkcols() {
	for x := range BOARD_DIM {
		constname := fmt.Sprintf("col_%d_sum_%d", x, b.ColTotals[x])
		var cells [BOARD_DIM]z3.Int
		for y := range BOARD_DIM {
			cells[y] = *address(x, y, &b.symbols)
		}
		pred := func(cmp z3.Int) z3.Bool {
			return cmp.Eq(b.intToConst(int(Wall)))
		}
		sum := b.countCells(cells[:], pred, constname)
		cond := sum.Eq(b.intToConst(b.ColTotals[x]))
		log.Debug(cond)
		b.slv.Assert(cond)
	}
}

func (b *Board) checkrows() {
	for y := range BOARD_DIM {
		constname := fmt.Sprintf("row_%d_sum_%d", y, b.ColTotals[y])
		var cells [BOARD_DIM]z3.Int
		for x := range BOARD_DIM {
			cells[x] = *address(x, y, &b.symbols)
		}
		pred := func(cmp z3.Int) z3.Bool {
			return cmp.Eq(b.intToConst(int(Wall)))
		}
		sum := b.countCells(cells[:], pred, constname)
		cond := sum.Eq(b.intToConst(b.RowTotals[y]))
		log.Debug(cond)
		b.slv.Assert(cond)
	}
}

// Monster must be surrounded by at least 3 walls (or edge of map)
func (b *Board) checkMonster(x int, y int) {
	if *address(x, y, &b.Cells) != Monster {
		log.Fatalf("%d,%d is not a monster", x, y)
	}
	var neighbor_sym []z3.Int
	maxSpaceNeighbors := b.intToConst(1)
	for _, neighbor := range adjacent {
		n_x := x + neighbor[0]
		n_y := y + neighbor[1]
		if b.inBounds(n_x, n_y) {
			neighbor_sym = append(neighbor_sym, *address(n_x, n_y, &b.symbols))
		} else {
			log.Warnf("neighbor %d,%d is out of bounds. skipping", n_x, n_y)
		}

	}
	constname := fmt.Sprintf("monster_%d_%d_deadend", x, y)
	pred := func(cmp z3.Int) z3.Bool {
		return cmp.NE(b.intToConst(int(Wall)))
	}
	sum := b.countCells(neighbor_sym, pred, constname)
	cond := sum.Eq(maxSpaceNeighbors)
	log.Debug(cond)
	b.slv.Assert(cond)
}

// space cells must have at least 2 neighbors that are NOT walls
func (b *Board) checkSpace(x int, y int) {
	cell := *address(x, y, &b.Cells)
	if cell != Space && cell != Unknown {
		log.Fatalf("%d,%d is not a monster", x, y)
	}
	constname := fmt.Sprintf("space_%d_%d_hallway", x, y)
	var neighbor_sym []z3.Int
	minNonWallNeighbors := b.intToConst(2)
	sym := *address(x, y, &b.symbols)
	cellIsSpace := sym.Eq(b.intToConst(int(Space)))
	for _, neighbor := range adjacent {
		n_x := x + neighbor[0]
		n_y := y + neighbor[1]
		if b.inBounds(n_x, n_y) {
			neighbor_sym = append(neighbor_sym, *address(n_x, n_y, &b.symbols))
		}
	}
	pred := func(cmp z3.Int) z3.Bool {
		return cmp.NE(b.intToConst(int(Wall)))
	}

	nonWallNeighbors := b.countCells(neighbor_sym, pred, constname)

	b.slv.Assert(cellIsSpace.Implies(nonWallNeighbors.GE(minNonWallNeighbors)))
}

func checkTreasure() {
	//room := append(chebyshevDistanceOffsets(1), chebyshevDistanceOffsets(0))
}

// Returns a list of offsets representing neighbors that are n chebyshev distance away
func chebyshevDistanceOffsets(d int) [][2]int {
	var neighbors [][2]int
	for x := -d; x <= d; x++ {
		if ix.Abs(x) != d {
			neighbors = append(neighbors, [2]int{x, d})
			neighbors = append(neighbors, [2]int{x, -d})
		} else {
			for y := -d; y <= d; y++ {
				neighbors = append(neighbors, [2]int{x, y})
			}
		}
	}
	return neighbors
}
