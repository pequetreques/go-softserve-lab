package main

import (
	"fmt"
	"math"
)

type function func([]int, int) int

func main() {
	// var s = []int{}
	// var s = []int{-15}
	// var s = []int{8, 7}
	// var s = []int{15, 3, 9, 11, 1, 5, 17, 11}
	// var s = []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	// var s = []int{10, 10, 10, 10, 10, 10, 10, 10}
	// var s = []int{15, 3, 9, -1, 11, 1, -7, -8, 0, 5, 17, 11, 0, -7}
	// var s = []int{1, 2, 3, 7, 50, 34, 12, 33, 35, 99, 65, 12, 67, 12, 77, 32, 78, 89, 54, 34, 55, 33, 66}
	var s = []int{82, 253, 252, 7, 7, 130, 101, 82, 253, 0, 252, -7, 252, 7, 63, -95, 15, 154, 98, -29}

	fmt.Printf("%v before\n", s)
	bucketSort(s, msBits)
	fmt.Printf("%v after\n", s)
}

func bucketSort(s []int, f function) {
	if len(s) < 2 {
		return
	}

	if len(s) == 2 {
		if s[0] > s[1] {
			s[0], s[1] = s[1], s[0]
		}

		return
	}

	buckets := make([][]int, len(s))

	min, max := s[0], s[0]

	for i, v := range s {
		buckets[i] = make([]int, 0)

		if v < min {
			min = v
		} else if v > max {
			max = v
		}
	}

	if min == max {
		return
	}

	for _, si := range s {
		bucketIndex := f(s, si)
		buckets[bucketIndex] = append(buckets[bucketIndex], si)
	}

	var i int
	for _, bucket := range buckets {
		if len(bucket) > 0 {
			bucketSort(bucket, f)

			for _, v := range bucket {
				s[i] = v
				i++
			}

		}
	}
}

func msBits(s []int, si int) int {
	min, max := s[0], s[0]

	for _, v := range s {
		if v < min {
			min = v
		} else if v > max {
			max = v
		}
	}

	num := si - min
	den := max - min + 1
	div := float32(num) / float32(den)
	mul := div * float32(len(s))
	result := math.Floor(float64(mul))

	return int(result) // ((si - min) / (max - min + 1)) * n
}
