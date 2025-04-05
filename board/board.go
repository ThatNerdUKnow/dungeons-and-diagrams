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

func (b *Board) checkcells() {
	for x := range BOARD_DIM {
		for y := range BOARD_DIM {
			switch *address(x, y, &b.Cells) {

			case Monster:
				{
					b.checkMonster(x, y)
				}
			case Space, Unknown:
				{
					b.checkSpace(x, y)
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

	if !b.inBounds(x, y) {
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

func address[T Cell | z3.Int](x int, y int, arr *[BOARD_DIM][BOARD_DIM]T) *T {
	if x > BOARD_DIM || y > BOARD_DIM || x < BOARD_MIN || y < BOARD_MIN {
		log.Fatalf("(%d,%d) is out of bounds", x, y)
	}
	return &arr[y][x]
}

func (b *Board) inBounds(x int, y int) bool {
	return (x < BOARD_DIM && x >= BOARD_MIN && y < BOARD_DIM && y >= BOARD_MIN)
}
