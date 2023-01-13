package insertionSort

import "testing"

func TestInsertionSort(t *testing.T) {
	var input = []int{15, 3, 9, 11, 1, 5, 17, 11}
	var output = []int{1, 3, 5, 9, 11, 11, 15, 17}

	InsertionSort(input)

	if !compareSlices(input, output) { // We can use reflect.DeepEqual as well
		t.Error("InsertionSort is not working properly.")
	}
}

func compareSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
