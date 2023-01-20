package main

import "fmt"

func main() {
	// var s = []int{}
	// var s = []int{-15}
	// var s = []int{8, 7}
	// var s = []int{10, 10, 10, 10, 10, 10, 10, 10}
	var s = []int{15, 3, 9, -1, 11, 1, -7, -8, 0, 5, 17, 11, 0, -7}
	// var s = []int{1, 2, 3, 7, 50, 34, 12, 33, 35, 99, 65, 12, 67, 12, 77, 32, 78, 89, 54, 34, 55, 33, 66}

	fmt.Printf("%v before\n", s)
	CountingSort(s)
	fmt.Printf("%v after\n", s)
}

func CountingSort(s []int) {
	if len(s) < 2 {
		return
	}

	if len(s) == 2 {
		if s[0] > s[1] {
			s[0], s[1] = s[1], s[0]
		}

		return
	}

	min, max := lookForMinAndMaxValues(s)
	size := max - min + 1
	fmt.Printf("min: %d, max: %d, size: %d\n", min, max, size)
	histogram := make([]int, size)

	var index int
	for _, v := range s {
		index = v - min
		histogram[index]++
	}

	index = 0
	for i, v := range histogram {
		value := i + min

		if v > 0 {
			for j := 0; j < v; j++ {
				s[index] = value
				index++
			}
		}
	}
}

func lookForMinAndMaxValues(s []int) (int, int) {
	min, max := s[0], s[0]

	for _, v := range s {
		if v < min {
			min = v
		} else if v > max {
			max = v
		}
	}

	return min, max
}
