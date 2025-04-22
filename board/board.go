package board

import (
	"fmt"

	"github.com/aclements/go-z3/z3"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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

func (e Cell) String() string {
	return EntTranslations[e]
}

type Board struct {
	Name         string
	Cells        [BOARD_DIM][BOARD_DIM]Cell
	ColTotals    [BOARD_DIM]*int
	RowTotals    [BOARD_DIM]*int
	colSymbols   [BOARD_DIM]z3.Int
	rowSymbols   [BOARD_DIM]z3.Int
	symbols      [BOARD_DIM][BOARD_DIM]z3.Int
	space_labels [BOARD_DIM][BOARD_DIM]z3.Int
	ctx          *z3.Context
	slv          *z3.Solver
}

func NewBoard(name string) Board {
	var rowTotals [8]*int
	var colTotals [8]*int
	return Board{
		Name:      name,
		ColTotals: colTotals,
		RowTotals: rowTotals,
	}
}

func (b *Board) checkcells() {
	for x := range BOARD_DIM {
		for y := range BOARD_DIM {
			switch *Address(x, y, &b.Cells) {

			case Monster:
				{
					b.checkMonster(x, y)
				}
			case Space, Unknown:
				{
					// current checkspace implementation is unsuitable for most puzzles
					// only really need this if we're doing puzzle generation
					b.checkSpace(x, y)
				}
			case Treasure:
				b.checkTreasure(x, y)
			}
		}
	}
}

func (b *Board) SetRowTotals(totals [BOARD_DIM]int) {
	for i, v := range totals {
		b.SetRowTotal(i)(&v)
	}
}

func setTotalValue(r int, arr *[BOARD_DIM]*int) func(*int) {

	if arr == nil {
		log.Fatal("Totals Array is nil")
	}
	return func(i *int) {
		log.Debugf("Setting total %d to %d for array %v", r, i, arr)
		arr[r] = i
	}
}

func (b *Board) SetRowTotal(r int) func(*int) {
	return setTotalValue(r, &b.RowTotals)
}

func (b *Board) SetColTotal(c int) func(*int) {
	return setTotalValue(c, &b.ColTotals)
}

func (b *Board) SetColTotals(totals [BOARD_DIM]int) {
	for i, v := range totals {
		b.SetColTotal(i)(&v)
	}
}

func (b *Board) SetCell(x int, y int, cell Cell) error {
	logger := log.With("x", x, "y", y)
	logger.Infof("Setting cell to %s", cell)
	if !b.inBounds(x, y) {
		return fmt.Errorf("coordinates (%d,%d) are out of bounds", x, y)
	}

	*Address(x, y, &b.Cells) = cell
	return nil
}

func (b *Board) GetCell(x int, y int) Cell {
	return *Address(x, y, &b.Cells)
}

func (b Board) String() string {
	t := table.New()

	t.BorderColumn(false)

	var header [BOARD_DIM + 1]string
	header[0] = " "
	for x := range BOARD_DIM {
		t := b.ColTotals[x]
		if t != nil {
			header[x+1] = fmt.Sprint(*t)
		} else {
			header[x+1] = "*"
		}
	}
	t.Row(header[:]...)

	for y := range BOARD_DIM {
		var row [BOARD_DIM + 1]string
		tr := b.RowTotals[y]
		if tr != nil {
			row[0] = fmt.Sprint(*b.RowTotals[y])
		} else {
			row[0] = "*"
		}
		for x := range BOARD_DIM {
			row[x+1] = fmt.Sprint(*Address(x, y, &b.Cells))
		}
		t.Row(row[:]...)
	}
	return lipgloss.JoinVertical(lipgloss.Left, b.Name, t.Render())
}

func (b *Board) Build() {
	if b.ctx == nil {
		cfg := z3.NewContextConfig()
		b.ctx = z3.NewContext(cfg)
	}
	if b.slv == nil {
		b.slv = z3.NewSolver(b.ctx)
	} else {
		b.slv.Reset()
	}

	// setting up the board
	for x := range BOARD_DIM {
		for y := range BOARD_DIM {
			c := Address(x, y, &b.Cells)
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
			*Address(x, y, &b.symbols) = sym

			label_constname := fmt.Sprintf("cell_%d_%d_label", x, y)
			label_sym := b.ctx.IntConst(label_constname)
			*Address(x, y, &b.space_labels) = label_sym
		}
	}
	b.checkcols()
	b.checkrows()
	b.checkcells()
	log.Debug(b.slv)
}

// Check if current constraints are solvable
func (b Board) Check() (bool, error) {
	sat, err := b.slv.Check()
	if err != nil {
		log.Error(err)
	}
	log.With("sat", sat).Debug("")
	return sat, err
}

func (b Board) Solve() (*Board, error) {
	log.Infof("solving %s", b.Name)
	b.Build()

	sat, err := b.Check()
	if !sat {
		return nil, err
	}
	m := b.slv.Model()
	log.Debug(m)
	nb := NewBoard(fmt.Sprintf("%s (solved)", b.Name))
	nb.RowTotals = b.RowTotals
	for x := range BOARD_DIM {
		colval, _, _ := m.Eval(b.colSymbols[x], true).(z3.Int).AsInt64()
		colvalPtr := int(colval)
		nb.SetColTotal(x)(&colvalPtr)
		rowval, _, _ := m.Eval(b.rowSymbols[x], true).(z3.Int).AsInt64()
		rowvalPtr := int(rowval)
		nb.SetRowTotal(x)(&rowvalPtr)
		for y := range BOARD_DIM {
			v := m.Eval(*Address(x, y, &b.symbols), true).(z3.Int)
			val, _, _ := v.AsInt64()
			*Address(x, y, &nb.Cells) = Cell(val)
		}
	}
	log.Debugf("Solved\n%s", nb)
	return &nb, nil
}

func Address[T Cell | z3.Int](x int, y int, arr *[BOARD_DIM][BOARD_DIM]T) *T {
	if x > BOARD_DIM || y > BOARD_DIM || x < BOARD_MIN || y < BOARD_MIN {
		log.Fatalf("(%d,%d) is out of bounds", x, y)
	}
	return &arr[y][x]
}

func (b *Board) inBounds(x int, y int) bool {
	return (x < BOARD_DIM && x >= BOARD_MIN && y < BOARD_DIM && y >= BOARD_MIN)
}
