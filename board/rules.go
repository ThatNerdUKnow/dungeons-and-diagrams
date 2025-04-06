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

func (b *Board) checkTreasure(x int, y int) {
	room := append(chebyshevDistanceOffsets(1), chebyshevDistanceOffsets(0)...)
	walls := chebyshevDistanceOffsets(2)
	// find center of room
	var cond *z3.Bool = nil
RoomLoop:
	for _, start := range room {
		// s_x and s_y are a guess as to where the room center is
		s_x := start[0] + x
		s_y := start[1] + y
		logger := log.With("room_center", [2]int{s_x, s_y})
		logger.Debugf("Checking room %d,%d", s_x, s_y)
		var room_sym []z3.Int
		for _, off := range room {
			// o_x and o_y are immediate neighbors of room center (0 <= chebyshev distance >=1 )
			o_x := off[0] + s_x
			o_y := off[1] + s_y
			// o_r is a list of cells in the room to test. if any element in the room is out of bounds
			// skip any further checks and move on to the next room

			if b.inBounds(o_x, o_y) {
				room_sym = append(room_sym, *address(o_x, o_y, &b.symbols))
			} else {
				logger.Debugf("%d,%d out of bounds. skipping treasure room check", o_x, o_y)
				continue RoomLoop
			}
		}
		logger = logger.With("room_symbols", room_sym)
		logger.Debug("")
		var wall_sym []z3.Int
		for _, border := range walls {
			w_x := s_x + border[0]
			w_y := s_y + border[1]
			if b.inBounds(w_x, w_y) {
				wall_sym = append(wall_sym, *address(w_x, w_y, &b.symbols))
			}
		}
		logger = logger.With("wall_symbols", wall_sym)
		logger.Debug("")
		entrance_pred := func(cmp z3.Int) z3.Bool {
			return cmp.NE(b.intToConst(int(Wall)))
		}
		// neighbors of chebyshev distance 2 may only contain 1 neighbor that is not a wall
		entrance := b.countCells(wall_sym, entrance_pred, fmt.Sprintf("room_%d_%d_entrance", s_x, s_y)).Eq(b.intToConst(1))

		// each room must have exactly one treasure
		treasure_pred := func(c z3.Int) z3.Bool {
			return c.Eq(b.intToConst(int(Treasure)))
		}
		treasure_count := b.countCells(room_sym, treasure_pred, fmt.Sprintf("room_%d_%d_treasure_count", s_x, s_y)).Eq(b.intToConst(1))
		// each room must contain exactly 8 spaces
		space_pred := func(c z3.Int) z3.Bool {
			return c.Eq(b.intToConst(int(Space)))
		}
		space_count := b.countCells(room_sym, space_pred, fmt.Sprintf("room_%d_%d_space_count", s_x, s_y)).Eq(b.intToConst(9 - 1))

		inner_cond := entrance.And(treasure_count).And(space_count)
		if cond == nil {
			cond = &inner_cond
		} else {
			res := cond.Xor(inner_cond)
			cond = &res
		}
	}
	b.slv.Assert(*cond)
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
