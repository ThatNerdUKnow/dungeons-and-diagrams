package board

import "testing"

func TestColTotals(t *testing.T) {
	b := NewBoard("")
	totals := [8]int{1, 2, 3, 4, 5, 6, 7, 8}
	b.SetColTotals(totals)
	var totals_cmp [8]int
	for i, v := range b.ColTotals {
		totals_cmp[i] = *v
	}
	if totals != totals_cmp {
		t.Errorf("%v != %v", totals, totals_cmp)
	}
}

func TestRowTotals(t *testing.T) {
	b := NewBoard("")
	totals := [8]int{1, 2, 3, 4, 5, 6, 7, 8}
	b.SetColTotals(totals)
	var totals_cmp [8]int
	for i, v := range b.ColTotals {
		totals_cmp[i] = *v
	}
	if totals != totals_cmp {
		t.Errorf("%v != %v", totals, totals_cmp)
	}
}
