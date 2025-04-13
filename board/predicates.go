package board

import "github.com/aclements/go-z3/z3"

func (b Board) predEq(c Cell) func(z3.Int) z3.Bool {
	return func(i z3.Int) z3.Bool {
		return i.Eq(b.intToConst(int(c)))
	}
}

func (b Board) predNE(c Cell) func(z3.Int) z3.Bool {
	return func(i z3.Int) z3.Bool {
		return i.NE(b.intToConst(int(c)))
	}
}
