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
	result := chebyshevDistanceOffsets(1)
	results := make(map[[2]int]bool)
	for _, neighbor := range result {
		results[neighbor] = true
	}
	t.Log(result)
	for pair, _ := range expected {
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
	result := chebyshevDistanceOffsets(0)
	results := make(map[[2]int]bool)
	for _, neighbor := range result {
		results[neighbor] = true
	}
	t.Log(result)
	for pair, _ := range expected {
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
	result := chebyshevDistanceOffsets(2)
	results := make(map[[2]int]bool)
	for _, neighbor := range result {
		results[neighbor] = true
	}
	t.Log(result)
	for pair, _ := range expected {
		if !results[pair] {
			t.Errorf("%d,%d not found in expected", pair[0], pair[1])
		}
	}
}
