package board

import (
	"fmt"
	"strings"

	"github.com/aclements/go-z3/z3"
	"github.com/charmbracelet/log"
	"github.com/enescakir/emoji"
)

type Cell int

const (
	Unknown Cell = iota
	Space
	Wall
	Monster
	Treasure
)

const (
	BOARD_DIM = 8
	BOARD_MIN = 0
)

var EntTranslations = map[Cell]string{
	Unknown:  string(emoji.QuestionMark),
	Space:    string(emoji.SmallBlueDiamond),
	Wall:     string(emoji.Brick),
	Monster:  string(emoji.Ogre),
	Treasure: string(emoji.GemStone),
}

func (e Cell) string() string {
	return EntTranslations[e]
}

type Board struct {
	Name      string
	Cells     [BOARD_DIM][BOARD_DIM]Cell
	ColTotals [BOARD_DIM]int
	RowTotals [BOARD_DIM]int
	symbols   [BOARD_DIM][BOARD_DIM]z3.Int
	ctx       *z3.Context
	slv       *z3.Solver
}

func NewBoard(name string) Board {
	return Board{
		Name: name,
	}
}

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

func (b *Board) checkcells() {
	neighbors := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for x := range BOARD_DIM {
		for y := range BOARD_DIM {
			switch *address(x, y, &b.Cells) {
			// Monster must be surrounded by at least 3 walls (or edge of map)
			case Monster:
				{
					var neighbor_sym []z3.Int
					maxSpaceNeighbors := b.intToConst(1)
					for _, neighbor := range neighbors {
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
			case Space, Unknown:
				{
					constname := fmt.Sprintf("space_%d_%d_hallway", x, y)
					var neighbor_sym []z3.Int
					minNonWallNeighbors := b.intToConst(2)
					cell := *address(x, y, &b.symbols)
					cellIsSpace := cell.Eq(b.intToConst(int(Space)))
					for _, neighbor := range neighbors {
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
			}
		}
	}
}

func (b *Board) SetRowTotals(totals [BOARD_DIM]int) {
	b.RowTotals = totals
}

func (b *Board) SetColTotals(totals [BOARD_DIM]int) {
	b.ColTotals = totals
}

func (b *Board) SetCell(x int, y int, cell Cell) error {
	if (x < BOARD_MIN || y < BOARD_MIN || x > len(b.Cells)) || y > len(b.Cells[0]) {
		return fmt.Errorf("coordinates (%d,%d) are out of bounds", x, y)
	}

	*address(x, y, &b.Cells) = cell
	return nil
}

func (b Board) String() string {
	var sb strings.Builder
	sb.WriteString("\"" + b.Name + "\"\n ")

	for _, sum := range b.ColTotals {
		sb.WriteString(fmt.Sprintf("%d ", sum))
	}
	sb.WriteString("\n")

	for y := range BOARD_DIM {
		sb.WriteString(fmt.Sprint(b.RowTotals[y]))
		for x := range BOARD_DIM {
			cell := address(x, y, &b.Cells)
			sb.WriteString(fmt.Sprint(cell.string()))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (b Board) Solve() (*Board, error) {
	log.Infof("solving %s", b.Name)
	cfg := z3.NewContextConfig()
	b.ctx = z3.NewContext(cfg)
	b.slv = z3.NewSolver(b.ctx)

	// setting up the board
	for x := range BOARD_DIM {
		for y := range BOARD_DIM {
			c := address(x, y, &b.Cells)
			constname := fmt.Sprintf("cell_%d_%d", x, y)
			var sym z3.Int
			switch *c {
			case Unknown:
				sym = b.ctx.IntConst(constname)
				space_or_wall := sym.Eq(b.intToConst(int(Wall))).Xor(sym.Eq(b.intToConst(int(Space))))
				b.slv.Assert(space_or_wall)
				log.Debugf("%d,%d %s", x, y, space_or_wall)
			default:
				sym = b.intToConst(int(*c))
			}

			*address(x, y, &b.symbols) = sym
			log.Debugf("created symbol %s", sym)
		}
	}

	/*
		for x := range BOARD_DIM {
			for y := range BOARD_DIM {
				c := b.Cells[y][x]
				constname := fmt.Sprintf("cell_%d_%d", x, y)

				switch c {
				case Monster:
					//hallway_constraint := ctx.BoolConst(fmt.Sprintf("%s_hallway_constraint", constname))
					hallway_sum := ctx.IntConst(fmt.Sprintf("%s_hallway_sum", constname))
					dim := BOARD_DIM - 1
					if x < dim {
						//hallway_constraint.Implies()
						cond := symbols[y][x+1].Eq(space)
						AddBoolToInt(&hallway_sum, &cond)
						//hallway_constraint = hallway_constraint.Xor(cond)
					}
					if x > 0 {
						cond := symbols[y][x-1].Eq(space)
						AddBoolToInt(&hallway_sum, &cond)
						//hallway_constraint = hallway_constraint.Xor(cond)
					}
					if y < dim {
						cond := symbols[y+1][x].Eq(space)
						AddBoolToInt(&hallway_sum, &cond)
						//hallway_constraint = hallway_constraint.Xor(cond)
					}
					if y > 0 {
						cond := symbols[y-1][x].Eq(space)
						AddBoolToInt(&hallway_sum, &cond)
						//hallway_constraint = hallway_constraint.Xor(cond)
					}
					//log.Debug(hallway_constraint)
					//slv.Assert(hallway_constraint)
					slv.Push()
					slv.Assert(hallway_sum.Eq(ctx.FromInt(1, ctx.IntSort()).(z3.Int)))
				}
			}
		}*/

	b.checkcols()
	b.checkrows()
	b.checkcells()
	log.Info(b.slv)
	sat, err := b.slv.Check()
	if !sat {
		return nil, err
	}
	m := b.slv.Model()
	log.Debug(m)
	nb := NewBoard(fmt.Sprintf("%s (solved)", b.Name))
	nb.SetColTotals(b.ColTotals)
	nb.SetRowTotals(b.RowTotals)
	for x := range BOARD_DIM {
		for y := range BOARD_DIM {
			v := m.Eval(*address(x, y, &b.symbols), true).(z3.Int)
			val, _, _ := v.AsInt64()
			*address(x, y, &nb.Cells) = Cell(val)
		}
	}

	return &nb, nil
}

func AddBoolToInt(i *z3.Int, b *z3.Bool) {
	c := i.Context()
	bi := b.IfThenElse(c.FromInt(1, c.IntSort()), c.FromInt(0, c.IntSort()))
	tmp := i.Add(bi.(z3.Int))
	*i = tmp
}

// Create an int const to represent a specific cell variant
func (b Board) intToConst(c int) z3.Int {
	if b.ctx == nil || b.slv == nil {
		log.Fatalf("Board %s does not have configured context or solver", b.Name)
	}
	return b.ctx.FromInt(int64(c), b.ctx.IntSort()).(z3.Int)
}

func address[T Cell | z3.Int](x int, y int, arr *[BOARD_DIM][BOARD_DIM]T) *T {
	if x > BOARD_DIM || y > BOARD_DIM || x < BOARD_MIN || y < BOARD_MIN {
		log.Fatalf("(%d,%d) is out of bounds", x, y)
	}
	return &arr[y][x]
}

func (b *Board) countCells(cells []z3.Int, pred func(z3.Int) z3.Bool, name string) z3.Int {
	sum := b.ctx.IntConst(name)
	b.slv.Assert(sum.Eq(b.intToConst(0)))
	//cmp := b.intToConst(int(t))
	for _, cell := range cells {
		p := pred(cell)
		AddBoolToInt(&sum, &p)
	}
	return sum
}

func (b *Board) inBounds(x int, y int) bool {
	return (x < BOARD_DIM && x >= BOARD_MIN && y < BOARD_DIM && y >= BOARD_MIN)
}
