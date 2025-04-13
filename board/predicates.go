package board

import "github.com/aclements/go-z3/z3"

type SymbolicPredicate = func(z3.Int) z3.Bool

// Variadic symbolic predicate function that evaluates to true if cell is one of the provided variants
func (b Board) predEq(cs ...Cell) SymbolicPredicate {
	return func(i z3.Int) z3.Bool {
		var cond *z3.Bool
		for _, c := range cs {
			inner_cond := i.Eq(b.intToConst(int(c)))
			if cond == nil {
				cond = &inner_cond
			} else {
				tmp := cond.Or(inner_cond)
				cond = &tmp
			}
		}
		return *cond
	}
}

// Variadic symbolic predicate function that evaluates to true if cell is not any of the provided variants
func (b Board) predNE(cs ...Cell) SymbolicPredicate {
	return func(i z3.Int) z3.Bool {
		var cond *z3.Bool
		for _, c := range cs {
			inner_cond := i.NE(b.intToConst(int(c)))
			if cond == nil {
				cond = &inner_cond
			} else {
				tmp := cond.And(inner_cond)
				cond = &tmp
			}
		}
		return *cond
	}
}
