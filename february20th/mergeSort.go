package main

import "fmt"

func main() {
	// var s = []int{}
	// var s = []int{-15}
	// var s = []int{8, 7}
	// var s = []int{4, 3, 2, 1}
	// var s = []int{15, 3, 9, 11, 1, 5, 17, 11}
	// var s = []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	// var s = []int{10, 10, 10, 10, 10, 10, 10, 10}
	// var s = []int{15, 3, 9, -1, 11, 1, -7, -8, 0, 5, 17, 11, 0, -7}
	// var s = []int{1, 2, 3, 7, 50, 34, 12, 33, 35, 99, 65, 12, 67, 12, 77, 32, 78, 89, 54, 34, 55, 33, 66}
	// var s = []int{82, 253, 252, 7, 7, 130, 101, 82, 253, 0, 252, -7, 252, 7, 63, -95, 15, 154, 98, -29}
	var s = []int{82, 253, 252, 3, 3, 130, 101, 82, 7, 0, 128, -7, 252, 3, 63, -95, -15, 154, 98, -29}

	fmt.Printf("%v before\n", s)
	s = mergeSort(s)
	fmt.Printf("%v after\n", s)

}

func mergeSort(s []int) []int {
	if len(s) < 2 {
		return s
	}

	if len(s) == 2 {
		if s[0] > s[1] {
			s[0], s[1] = s[1], s[0]
		}

		return s
	}

	left := s[:len(s)/2]
	right := s[len(s)/2:]

	if len(s) <= 16 {
		fmt.Printf("Left: %v Right: %v\n", left, right)
	}

	left = mergeSort(left)
	right = mergeSort(right)

	return merge(left, right)
}

func merge(left, right []int) []int {
	var result []int

	j := 0
	k := 0

	for {
		if j < len(left) && k < len(right) {
			if left[j] < right[k] {
				result = append(result, left[j])
				j++
			} else {
				result = append(result, right[k])
				k++
			}
		} else {
			break
		}
	}

	for {
		if j < len(left) {
			result = append(result, left[j])
			j++
		} else {
			break
		}
	}

	for {
		if k < len(right) {
			result = append(result, right[k])
			k++
		} else {
			break
		}
	}

	return result
}
