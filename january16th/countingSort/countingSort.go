package main

import "fmt"

func main() {
	var s = []int{15, 3, 9, -1, 11, 1, -7, -8, 0, 5, 17, 11, 0, -7}

	fmt.Printf("%v before\n", s)
	CountingSort(s)
	fmt.Printf("%v after\n", s)
}

func CountingSort(s []int) {
	var index int

	min, max := lookForMinAndMaxValues(s)
	size := max - min + 1
	fmt.Printf("min: %d, max: %d, size: %d\n", min, max, size)
	counters := make([]int, size)

	for _, v := range s {
		index = v - min
		counters[index]++
	}

	index = 0
	for i, v := range counters {
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
	var min int
	var max int

	for _, v := range s {
		if v < min {
			min = v
		} else if v > max {
			max = v
		}
	}

	return min, max
}
