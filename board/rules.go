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

func (b *Board) checkline(i int, isCol bool, line_sym *z3.Int, line_total *int) {
	var cells [BOARD_DIM]z3.Int
	for j := range BOARD_DIM {
		if isCol {
			cells[j] = *address(i, j, &b.symbols)
		} else {
			cells[j] = *address(j, i, &b.symbols)
		}
	}
	// symbol representing the number of walls in this line
	wall_sum := b.countCells(cells[:], b.predEq(Wall), nil)

	// build final constraint
	var sym_prefix string
	if isCol {
		sym_prefix = "col"
	} else {
		sym_prefix = "row"
	}

	if line_total != nil {
		// if our line total is concrete, line_sym[i] and wall_sum should be equal
		*line_sym = b.ctx.IntConst(fmt.Sprintf("%s_%d_sum_%d", sym_prefix, i, *line_total))
		b.slv.Assert(line_sym.Eq(b.intToConst(*line_total)))
	} else {
		*line_sym = b.ctx.IntConst(fmt.Sprintf("%s_%d_sum_unknown", sym_prefix, i))
		b.slv.Assert(b.colSymbols[i].GE(b.intToConst(0)).And(b.colSymbols[i].LE(b.intToConst(BOARD_DIM))))
	}
	cond := line_sym.Eq(wall_sum)
	log.Debug(cond)
	b.slv.Assert(cond)
}

func (b *Board) checkcols() {
	for i := range BOARD_DIM {
		b.checkline(i, true, &b.colSymbols[i], b.ColTotals[i])
	}
}

func (b *Board) checkrows() {
	for i := range BOARD_DIM {
		b.checkline(i, false, &b.rowSymbols[i], b.RowTotals[i])
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
			log.With("origin", [2]int{x, y}, "neighbor", [2]int{n_x, n_y}).Warn("skipping out of bounds neighbor")
		}

	}
	constname := fmt.Sprintf("monster_%d_%d_deadend", x, y)
	sum := b.countCells(neighbor_sym, b.predNE(Wall), &constname)
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

	nonWallNeighbors := b.countCells(neighbor_sym, b.predNE(Wall), &constname)

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
		// neighbors of chebyshev distance 2 may only contain 1 neighbor that is not a wall
		constname := fmt.Sprintf("room_%d_%d_entrance", s_x, s_y)
		entrance := b.countCells(wall_sym, b.predNE(Wall), &constname).Eq(b.intToConst(1))

		// each room must have exactly one treasure
		treasure_count := b.countCells(room_sym, b.predEq(Treasure), &constname).Eq(b.intToConst(1))
		// each room must contain exactly 8 spaces
		space_count := b.countCells(room_sym, b.predEq(Space), &constname).Eq(b.intToConst(9 - 1))

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
