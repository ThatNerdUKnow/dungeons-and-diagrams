package board

import (
	"testing"
)

func TestChebyshevDistanceOffsets1(t *testing.T) {
	expected := make(map[[2]int]bool)
	expected_neighbors := [][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}
	for _, neighbor := range expected_neighbors {
		expected[neighbor] = true
	}
	result := chebyshevDistanceNeighbors(1)
	results := make(map[[2]int]bool)
	for _, neighbor := range result {
		results[neighbor] = true
	}
	t.Log(result)
	for pair := range expected {
		if !results[pair] {
			t.Errorf("%d,%d not found in expected", pair[0], pair[1])
		}
	}
}

func TestChebyshevDistanceOffsets0(t *testing.T) {
	expected := make(map[[2]int]bool)
	expected_neighbors := [][2]int{{0, 0}}
	for _, neighbor := range expected_neighbors {
		expected[neighbor] = true
	}
	result := chebyshevDistanceNeighbors(0)
	results := make(map[[2]int]bool)
	for _, neighbor := range result {
		results[neighbor] = true
	}
	t.Log(result)
	for pair := range expected {
		if !results[pair] {
			t.Errorf("%d,%d not found in expected", pair[0], pair[1])
		}
	}
}

func TestChebyshevDistanceOffsets2(t *testing.T) {
	expected := make(map[[2]int]bool)
	expected_neighbors := [][2]int{{-2, -2}, {-2, -1}, {-2, 0}, {-2, 1}, {-2, 2}, {-1, -2}, {-1, 2}, {0, -2}, {0, 2}, {1, -2}, {1, 2}, {2, -2}, {2, -1}, {2, 0}, {2, 1}, {2, 2}}
	for _, neighbor := range expected_neighbors {
		expected[neighbor] = true
	}
	result := chebyshevDistanceNeighbors(2)
	results := make(map[[2]int]bool)
	for _, neighbor := range result {
		results[neighbor] = true
	}
	t.Log(result)
	for pair := range expected {
		if !results[pair] {
			t.Errorf("%d,%d not found in expected", pair[0], pair[1])
		}
	}
}

func TestMonsterFail(t *testing.T) {
	b := NewBoard("TestMonsterFail")
	b.SetCell(0, 0, Monster)
	b.SetCell(0, 1, Monster)
	t.Log(b)
	nb, err := b.Solve()
	t.Log(nb)
	if nb != nil {
		t.Errorf("Should not be satisfiable %v", err)
	}
}

func TestMonster(t *testing.T) {
	b := NewBoard("TestMonster")
	b.SetCell(5, 5, Monster)
	t.Log(b)
	nb, err := b.Solve()
	t.Log(nb)
	if nb == nil {
		t.Errorf("Should be satisfiable %v", err)
	}
}

func TestTreasureFail(t *testing.T) {
	b := NewBoard("TestTreasureFail")
	b.SetCell(0, 0, Treasure)
	b.SetCell(0, 1, Treasure)
	t.Log(b)
	nb, err := b.Solve()
	t.Log(nb)
	if nb != nil {
		t.Errorf("Should not be satisfiable %v", err)
	}
}

func TestColCountSymbolic(t *testing.T) {
	b := NewBoard("TestColCountSymbolic")
	b.Build()

	b.SetCell(0, 0, Space)
	b.SetCell(0, 1, Wall)
	b.SetCell(0, 2, Wall)
	b.SetCell(0, 3, Wall)
	b.SetCell(0, 4, Wall)
	b.SetCell(0, 5, Wall)
	b.SetCell(0, 6, Space)
	b.SetCell(0, 7, Space)
	t.Log(b)
	b2, _ := b.Solve()
	if b2 == nil {
		t.Fatal("Should be satisfiable")
	}
	t.Log(b2)
	var sym []Cell
	for i := 0; i < BOARD_DIM; i++ {
		sym = append(sym, *Address(0, i, &b2.Cells))
	}
	if len(sym) != BOARD_DIM {
		t.Error("Off by one")
	}
	total := CountCellsConcrete(sym, int(Wall))

	b2ColTotal := b2.ColTotals[0]
	if b2ColTotal == nil {
		t.Fatal("col total should not be nil")
	}
	if total != *b2ColTotal {
		t.Errorf("Total does not match (got %d, expected %d)", *b2ColTotal, total)
	}

}

func TestColCountConcrete(t *testing.T) {
	b := NewBoard("TestColCountConcrete")
	b.Build()

	targetTotal := 5
	b.SetColTotal(0)(&targetTotal)
	b2, _ := b.Solve()
	if b2 == nil {
		t.Fatal("Should be satisfiable")
	}
	t.Log(b2)
	var sym []Cell
	for i := 0; i < BOARD_DIM; i++ {
		sym = append(sym, *Address(0, i, &b2.Cells))
	}
	if len(sym) != BOARD_DIM {
		t.Error("Off by one")
	}
	total := CountCellsConcrete(sym, int(Wall))

	b2ColTotal := b2.ColTotals[0]
	if b2ColTotal == nil {
		t.Fatal("col total should not be nil")
	}
	if total != *b2ColTotal {
		t.Errorf("Total does not match (got %d, expected %d)", *b2ColTotal, total)
	}

}
