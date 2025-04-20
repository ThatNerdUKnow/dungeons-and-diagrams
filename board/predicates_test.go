package board

import (
	"github.com/aclements/go-z3/z3"
)

func CountCellsConcrete(sym []Cell, i int) int {
	total := 0
	for _, sym := range sym {

		if int(sym) == i {
			total++
		}
	}
	return total
}

func UpdateSymbols(sym *[]z3.Int, model *z3.Model) {
	for i, s := range *sym {
		v := model.Eval(s, true).(z3.Int)
		(*sym)[i] = v
	}
}

func CountCell(name *string, concreteTotal bool) (z3.Int, []z3.Int) {
	b := NewBoard("")
	b.Build()
	slv := z3.NewSolver(b.ctx)
	var sym []z3.Int
	var count z3.Int
	var total z3.Int
	for i := 0; i < 5; i++ {
		sym = append(sym, b.intToConst(int(Wall)))
	}
	count = b.countCells(sym, b.predEq(Wall), name)

	if concreteTotal {
		total = b.intToConst(5)
	} else {
		total = b.ctx.IntConst("Total_Symbolic")
	}

	slv.Assert(count.Eq(total))
	slv.Check()
	m := slv.Model()
	UpdateSymbols(&sym, m)
	total = m.Eval(total, true).(z3.Int)
	return total, sym
}

/*
func TestCountCellConcreteNameNil(t *testing.T) {
	totalsym, sym := CountCell(nil, true)
	total := CountCellsConcrete(sym, int(Wall))
	logger := log.With("sym", sym, "totalsym", totalsym, "total", total)
	logger.Info("Checking cell totals")

	v, _, _ := totalsym.AsInt64()
	if int(v) != total {
		t.Fatal("Totals do not match")
	}
}
*/
