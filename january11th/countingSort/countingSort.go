package main

import "fmt"

func main() {
	var s = []int{15, 3, 9, 11, 1, 5, 17, 11}

	fmt.Printf("%v before\n", s)
	countingSort(s)
	fmt.Printf("%v after\n", s)
}

func countingSort(s []int) {
	maxValue := lookForMaxValue(s)
	memory := make([]int, maxValue)

	for _, v := range s {
		mIndex := v - 1
		memory[mIndex]++
	}

	sIndex := 0
	for i, v := range memory {
		if v > 0 {
			for j := 0; j < v; j++ {
				s[sIndex] = i + 1
				sIndex++
			}
		}
	}
}

func lookForMaxValue(s []int) int {
	var maxValue int

	for _, v := range s {
		if v > maxValue {
			maxValue = v
		}
	}

	return maxValue
}
