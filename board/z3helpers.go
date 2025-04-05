package board

import (
	"github.com/aclements/go-z3/z3"
	"github.com/charmbracelet/log"
)

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
