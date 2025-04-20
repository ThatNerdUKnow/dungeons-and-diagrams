package board

import (
	"log"

	"github.com/aclements/go-z3/z3"
)

type SymbolicPredicate = func(z3.Int) z3.Bool
type SymbolicCombinator = func(z3.Bool, ...z3.Bool) z3.Bool
type SymbolicComparator = func(z3.Int, z3.Int) z3.Bool
type PredicateCombinator = func(SymbolicCombinator, ...SymbolicPredicate) SymbolicPredicate

// symbolic predicate function that evaluates to true if cell is one of the provided variants
func (b Board) predEq(c Cell) SymbolicPredicate {
	return b.pred(z3.Int.Eq)(c)
}

// symbolic predicate function that evaluates to true if cell is not any of the provided variants
func (b Board) predNE(c Cell) SymbolicPredicate {
	return b.pred(z3.Int.NE)(c)
}

// Symbolic predicate function generator. p is a comparator function between symbolic ints and c is a constant cell
func (b Board) pred(p SymbolicComparator) func(Cell) SymbolicPredicate {
	return func(c Cell) SymbolicPredicate {
		return func(i z3.Int) z3.Bool {
			return p(i, b.intToConst(int(c)))
		}
	}
}

// Compose multiple symbolic predicates using a symbolic combinator to create a new composite symbolic predicate
func predCompose(c SymbolicCombinator, ps ...SymbolicPredicate) SymbolicPredicate {
	if len(ps) == 0 {
		log.Fatal("Predicate list is empty")
	}
	return func(i z3.Int) z3.Bool {
		var cond *z3.Bool
		for _, p := range ps {
			if cond == nil {
				tmp := p(i)
				cond = &tmp
			} else {
				*cond = c(*cond, p(i))
			}
		}
		return *cond
	}
}
