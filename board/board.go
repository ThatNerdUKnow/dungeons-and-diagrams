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
}

func NewBoard(name string) Board {
	return Board{
		Name: name,
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
	b.Cells[y][x] = cell
	return nil
}

func (b Board) String() string {
	var sb strings.Builder
	sb.WriteString("\"" + b.Name + "\"\n ")

	for _, sum := range b.ColTotals {
		sb.WriteString(fmt.Sprintf("%d ", sum))
	}
	sb.WriteString("\n")

	for y, col := range b.Cells {
		sb.WriteString(fmt.Sprint(b.RowTotals[y]))
		for _, cell := range col {
			//sb.WriteString(fmt.Sprintf("(%d,%d) ", x, y))
			sb.WriteString(fmt.Sprint(cell.string()))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (b Board) Solve() (*Board, error) {
	log.Infof("solving %s", b.Name)
	cfg := z3.NewContextConfig()
	ctx := z3.NewContext(cfg)
	slv := z3.NewSolver(ctx)

	// constant values
	//unk := ctx.FromInt(int64(Unknown), ctx.IntSort()).(z3.Int)
	//max_sym := ctx.FromInt(int64(len(EntTranslations)), ctx.IntSort()).(z3.Int)
	wall := ctx.FromInt(int64(Wall), ctx.IntSort()).(z3.Int)
	space := ctx.FromInt(int64(Space), ctx.IntSort()).(z3.Int)
	var symbols [BOARD_DIM][BOARD_DIM]z3.Int
	for x := range BOARD_DIM {
		for y := range BOARD_DIM {
			c := b.Cells[y][x]
			constname := fmt.Sprintf("cell_%d_%d", x, y)
			var sym z3.Int
			switch c {
			case Unknown:
				sym = ctx.IntConst(constname)
				space_or_wall := sym.Eq(wall).Xor(sym.Eq(space))
				//space_or_wall := sym.Eq(space).Or(sym.Eq(wall))
				slv.Assert(space_or_wall)
				log.Debugf("%d,%d %s", x, y, space_or_wall)
			default:
				sym = ctx.FromInt(int64(c), ctx.IntSort()).(z3.Int)
			}

			symbols[y][x] = sym
			//log.Tracef("created symbol %s", sym)
		}
	}

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
	}

	for col_x := range BOARD_DIM {

		col_sum := ctx.IntConst(fmt.Sprintf("col_%d_sum", col_x))
		for y := range BOARD_DIM {
			iswall := symbols[y][col_x].Eq(wall)
			AddBoolToInt(&col_sum, &iswall)
		}
		tot := ctx.FromInt(int64(b.ColTotals[col_x]), ctx.IntSort()).(z3.Int)
		assertion := col_sum.Eq(tot)
		log.Debug(assertion)
		slv.Assert(assertion)
	}

	sat, err := slv.Check()
	if !sat {
		return nil, err
	}
	m := slv.Model()
	//fmt.Println(m)
	nb := NewBoard(fmt.Sprintf("%s (solved)", b.Name))
	nb.SetColTotals(b.ColTotals)
	nb.SetRowTotals(b.RowTotals)
	for x := range BOARD_DIM {
		for y := range BOARD_DIM {
			v := m.Eval(symbols[y][x], true).(z3.Int)
			val, _, _ := v.AsInt64()
			nb.Cells[y][x] = Cell(val)
		}
	}

	return &nb, nil
}

func AddBoolToInt(i *z3.Int, b *z3.Bool) {
	//log.Debugf("Adding %s to %s", i, b)
	c := i.Context()
	bi := b.IfThenElse(c.FromInt(1, c.IntSort()), c.FromInt(0, c.IntSort()))
	tmp := i.Add(bi.(z3.Int))
	*i = tmp
}
