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
		sum := b.ctx.IntConst(constname)
		for y := range BOARD_DIM {
			cond := address(x, y, &b.symbols).Eq(b.intToConst(int(Wall)))
			//cond := b.symbols[y][x].Eq(b.intToConst(int(Wall)))
			AddBoolToInt(&sum, &cond)
		}
		//assertion_const := b.ctx.BoolConst(fmt.Sprintf("%s-satisfied", constname))
		assertion := sum.Eq(b.intToConst(b.ColTotals[x]))
		log.Debug(assertion)
		b.slv.Assert(assertion)
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
	//b.Cells[y][x] = cell
	return nil
}

func (b Board) String() string {
	var sb strings.Builder
	sb.WriteString("\"" + b.Name + "\"\n ")

	for _, sum := range b.ColTotals {
		sb.WriteString(fmt.Sprintf("%d ", sum))
	}
	sb.WriteString("\n")

	/*
		for y, col := range b.Cells {
			sb.WriteString(fmt.Sprint(b.RowTotals[y]))
			for _, cell := range col {
				//sb.WriteString(fmt.Sprintf("(%d,%d) ", x, y))
				sb.WriteString(fmt.Sprint(cell.string()))
			}
			sb.WriteString("\n")
		}*/

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
			//c := b.Cells[y][x]
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
			//nb.Cells[y][x] = Cell(val)
		}
	}

	return &nb, nil
}

func AddBoolToInt(i *z3.Int, b *z3.Bool) {
	//log.Debugf("Adding %s to %s", i, b)
	c := i.Context()
	//bi := b.IfThenElse(i.Add(c.FromInt(1, c.IntSort()).(z3.Int)), i)
	//*i = bi.(z3.Int)
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
