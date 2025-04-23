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
			cells[j] = *Address(i, j, &b.symbols)
		} else {
			cells[j] = *Address(j, i, &b.symbols)
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
		cond := line_sym.Eq(b.intToConst(*line_total))
		log.Debug(cond)
		b.slv.Assert(cond)
	} else {
		*line_sym = b.ctx.IntConst(fmt.Sprintf("%s_%d_sum_unknown", sym_prefix, i))
		cond := b.colSymbols[i].GE(b.intToConst(0)).And(b.colSymbols[i].LE(b.intToConst(BOARD_DIM)))
		log.Debug(cond)
		b.slv.Assert(cond)
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
	if *Address(x, y, &b.Cells) != Monster {
		log.Fatalf("%d,%d is not a monster", x, y)
	}
	var neighbor_sym []z3.Int
	maxSpaceNeighbors := b.intToConst(1)
	for _, neighbor := range adjacent {
		n_x := x + neighbor[0]
		n_y := y + neighbor[1]
		if b.inBounds(n_x, n_y) {
			neighbor_sym = append(neighbor_sym, *Address(n_x, n_y, &b.symbols))
		} else {
			log.With("origin", [2]int{x, y}, "neighbor", [2]int{n_x, n_y}).Warn("skipping out of bounds neighbor")
		}

	}
	constname := fmt.Sprintf("monster_%d_%d_deadend", x, y)

	predCompose(z3.Bool.Or, nil)
	//sum_wall := b.countCells(neighbor_sym, b.predNE(Wall), &constname)
	sum_wall := b.countCells(neighbor_sym, predCompose(z3.Bool.Or, b.predEq(Space), b.predEq(Treasure)), &constname)
	sum_monster := b.countCells(neighbor_sym, b.predEq(Monster), nil)
	cond := sum_wall.Eq(maxSpaceNeighbors).And(sum_monster.Eq(b.intToConst(0)))
	//log.Debug(cond)
	b.slv.Assert(cond)
}

// space cells must have at least 2 neighbors that are NOT walls
func (b *Board) checkSpace(x int, y int) {
	// trying out a flood fill
	cell := *Address(x, y, &b.symbols)
	cell_label := *Address(x, y, &b.space_labels)
	var neighbors []z3.Int
	var preds []SymbolicPredicate
	//var neighbors_labels []z3.Int
	for _, neighbor := range adjacent {
		nx := neighbor[0] + x
		ny := neighbor[1] + y
		if b.inBounds(nx, ny) {
			neighbors = append(neighbors, *Address(nx, ny, &b.symbols))
			nlabel := *Address(nx, ny, &b.space_labels)
			// if neighbor is a space, then it must have a label one less than the current cell
			f := func(i z3.Int) z3.Bool {
				is_space := b.predEq(Space)(i)
				expected_label := cell_label.Sub(b.intToConst(1))
				nlabel_expected := nlabel.Eq(expected_label)
				return is_space.And(nlabel_expected)
			}
			preds = append(preds, f)
		}
	}
	constname := fmt.Sprintf("cell_%d_%d_neighbor", x, y)
	count := b.countCellsMany(neighbors, preds, &constname)
	// if current cell is a space and its label is NOT zero, then there must be at least once neighbor
	// with a label (this cell's label) - 1
	cond := b.predEq(Space)(cell).And(cell_label.NE(b.intToConst(0))).Implies(count.GE(b.intToConst(1)))
	b.slv.Assert(cond)
	/*
		cell := *Address(x, y, &b.Cells)
		if cell != Space && cell != Unknown {
			log.Fatalf("%d,%d is not a monster", x, y)
		}
		constname := fmt.Sprintf("space_%d_%d_hallway", x, y)
		var neighbor_sym []z3.Int
		minNonWallNeighbors := b.intToConst(2)
		sym := *Address(x, y, &b.symbols)
		cellIsSpace := sym.Eq(b.intToConst(int(Space)))
		for _, neighbor := range adjacent {
			n_x := x + neighbor[0]
			n_y := y + neighbor[1]
			if b.inBounds(n_x, n_y) {
				neighbor_sym = append(neighbor_sym, *Address(n_x, n_y, &b.symbols))
			}
		}

		nonWallNeighbors := b.countCells(neighbor_sym, b.predNE(Wall), &constname)
		cond := cellIsSpace.Implies(nonWallNeighbors.GE(minNonWallNeighbors))
		log.Debug(cond)
		b.slv.Assert(cond)*/
}

func (b *Board) checkTreasure(x int, y int) {
	room := append(chebyshevDistanceNeighbors(1), chebyshevDistanceNeighbors(0)...)
	walls := chebyshevDistanceNeighborsSansCorners(2)
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
				room_sym = append(room_sym, *Address(o_x, o_y, &b.symbols))
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
				wall_sym = append(wall_sym, *Address(w_x, w_y, &b.symbols))
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
func chebyshevDistanceNeighbors(d int) [][2]int {
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

func chebyshevDistanceNeighborsSansCorners(d int) [][2]int {
	var neighbors [][2]int
	cneighbors := chebyshevDistanceNeighbors(2)
	for _, neighbor := range cneighbors {
		x := ix.Abs(neighbor[0])
		y := ix.Abs(neighbor[1])
		if x == d && y == d {
			log.Debug("Skipping neighbor %s", neighbor)

		} else {
			log.Debug("Appending neighbor %s", neighbor)
			neighbors = append(neighbors, neighbor)
		}
	}
	return neighbors
}
