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

// This function works fine if only one int symbol needs to be accounted for,
// but isn't flexible enough to handle more dynamic constraints
// however, i'm keeping this function with this signature in order to prevent an API break
func (b *Board) countCells(cells []z3.Int, pred SymbolicPredicate, name *string) z3.Int {
	return b.countCellsMany(cells, []SymbolicPredicate{pred}, name)
}

// Allows more flexibility in the conditions for counting cells. each predicate function may capture information
// from its parent scope, which means that I can encode more complicated constraints in the predicates than I could otherwise
func (b *Board) countCellsMany(cells []z3.Int, preds []SymbolicPredicate, name *string) z3.Int {
	// if only one predicate is provided, apply the predicate to each cell
	// otherwise it's expected that there is a distinct predicate function for each cell
	// with 1:1 cardinality
	if len(cells) != len(preds) && len(preds) != 1 {
		log.Fatalf("predicate and cell count do not match")
	}
	singlePred := len(preds) == 1
	var sum z3.Int
	if name != nil {
		sum = b.ctx.IntConst(*name)
		b.slv.Assert(sum.Eq(b.intToConst(0)))
	} else {
		sum = b.intToConst(0)
	}

	for i, cell := range cells {
		var pred SymbolicPredicate
		if singlePred {
			pred = preds[0]
		} else {
			pred = preds[i]
		}
		p := pred(cell)
		AddBoolToInt(&sum, &p)
	}
	return sum
}
