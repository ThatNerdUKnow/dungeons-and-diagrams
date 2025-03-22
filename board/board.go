package board

import (
	"fmt"
	"strings"

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

/*
func (b Board) Solve() Board {
	cfg := z3.NewContextConfig()
	ctx := z3.NewContext(cfg)
	slv := z3.NewSolver(ctx)

	// constant values
	unk := ctx.FromInt(int64(Unknown), ctx.IntSort()).(z3.Int)
	max_sym := ctx.FromInt(int64(len(EntTranslations)), ctx.IntSort()).(z3.Int)
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
				slv.Assert(sym.Eq(wall).Or(sym.Eq(space)))
			default:
				sym = ctx.FromInt(int64(c), ctx.IntSort()).(z3.Int)
			}
			symbols[y][x] = sym
			// no cells may be "unknown"
			slv.Assert(sym.GT(unk))
			slv.Assert(sym.LT(max_sym))
		}
	}

	for x := range BOARD_DIM {
		for y := range BOARD_DIM {
			c := b.Cells[y][x]
			switch c {
			case Monster:

			case Treasure:
			}
		}
	}
}
*/
