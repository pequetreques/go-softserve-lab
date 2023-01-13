package insertionSort

func InsertionSort(s []int) {
	var sorted int
	var unsorted int
	var box int
	var miniSorted int
	var miniUnsorted int

	unsorted = 1

	for i := 0; i < len(s)-1; i++ {
		sorted = unsorted - 1

		if s[sorted] > s[unsorted] {
			box = s[sorted]
			s[sorted] = s[unsorted]
			s[unsorted] = box
		}

		unsorted++

		if unsorted > 2 {
			miniUnsorted = unsorted - 1
			miniSorted = miniUnsorted - 1

			for miniSorted >= 0 {
				if s[miniSorted] > s[miniUnsorted] {
					box = s[miniSorted]
					s[miniSorted] = s[miniUnsorted]
					s[miniUnsorted] = box
				}

				miniUnsorted--
				miniSorted = miniUnsorted - 1
			}
		}
	}
}
