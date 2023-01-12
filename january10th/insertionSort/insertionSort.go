package main

import "fmt"

func main() {
	var s = []int{15, 3, 9, 11, 1, 5, 17, 11}

	fmt.Printf("%v before\n", s)
	iterations := insertionSort(s)
	fmt.Printf("%v after\n", s)
	fmt.Printf("Total iterations: %d\n", iterations)
}

func insertionSort(s []int) int {
	var sorted int
	var unsorted int
	var box int
	var miniSorted int
	var miniUnsorted int
	var iterations int

	unsorted = 1
	iterations = 0

	for i := 0; i < len(s)-1; i++ {
		sorted = unsorted - 1

		if s[sorted] > s[unsorted] {
			box = s[sorted]
			s[sorted] = s[unsorted]
			s[unsorted] = box
			fmt.Printf("%v main loop\n", s)
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
					fmt.Printf("%v nested loop\n", s)
				}

				miniUnsorted--
				miniSorted = miniUnsorted - 1
				iterations++
			}
		}

		iterations++
	}

	return iterations
}

func traverseSliceUsingIndex(s []int) {
	for i := 0; i < len(s); i++ {
		fmt.Println(s[i])
	}
}

func traverseSliceUsingRange(s []int) {
	for index, value := range s {
		fmt.Printf("index: %d, value: %d\n", index, value)
	}
}
