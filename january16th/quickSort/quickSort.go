package main

import "fmt"

type function func([]int) int

func main() {
	// var s = []int{}
	// var s = []int{-15}
	// var s = []int{8, 7}
	// var s = []int{10, 10, 10, 10, 10, 10, 10, 10}
	// var s = []int{15, 3, 9, -1, 11, 1, -7, -8, 0, 5, 17, 11, 0, -7}
	// var s = []int{1, 2, 3, 7, 50, 34, 12, 33, 35, 99, 65, 12, 67, 12, 77, 32, 78, 89, 54, 34, 55, 33, 66}
	var s = []int{82, 253, 252, 7, 7, 130, 101, 82, 253, 0, 252, -7, 252, 7, 63, -95, 15, 154, 98, -29}

	fmt.Printf("%v before\n", s)
	quickSort1(s, 0, len(s)-1)
	fmt.Printf("%v after\n", s)

	fmt.Println()
	// var t = []int{}
	// var t = []int{-15}
	// var t = []int{8, 7}
	// var t = []int{10, 10, 10, 10, 10, 10, 10, 10}
	// var t = []int{15, 3, 9, -1, 11, 1, -7, -8, 0, 5, 17, 11, 0, -7}
	// var t = []int{1, 2, 3, 7, 50, 34, 12, 33, 35, 99, 65, 12, 67, 12, 77, 32, 78, 89, 54, 34, 55, 33, 66}
	var t = []int{82, 253, 252, 7, 7, 130, 101, 82, 253, 0, 252, -7, 252, 7, 63, -95, 15, 154, 98, -29}

	fmt.Printf("%v before\n", t)
	quickSort2(t, getPivot)
	fmt.Printf("%v after\n", t)
}

func quickSort1(s []int, left, right int) {
	if len(s) < 2 {
		return
	}

	if len(s) == 2 {
		if s[0] > s[1] {
			s[0], s[1] = s[1], s[0]
		}

		return
	}

	if left >= right || left < 0 {
		return
	}

	p := partition(s, left, right)

	quickSort1(s, left, p-1)
	quickSort1(s, p+1, right)
}

func partition(s []int, left, right int) int {
	i := left
	pivot := s[right]

	for j := left; j < right; j++ {
		if s[j] <= pivot {
			s[i], s[j] = s[j], s[i]
			i++
		}
	}

	s[i], s[right] = s[right], s[i]

	return i
}

func quickSort2(s []int, f function) {
	if len(s) < 2 {
		return
	}

	if len(s) == 2 {
		if s[0] > s[1] {
			s[0], s[1] = s[1], s[0]
		}

		return
	}

	p := f(s)

	quickSort2(s[:p], f)
	quickSort2(s[p:], f)
}

func getPivot(s []int) int {
	i := 0
	right := len(s) - 1
	pivot := s[right]

	for j := 0; j < right; j++ {
		if s[j] <= pivot {
			s[i], s[j] = s[j], s[i]
			i++
		}
	}

	s[i], s[right] = s[right], s[i]

	return i
}
