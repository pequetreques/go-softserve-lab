package main

import "fmt"

func main() {
	// numbers := [11]int{7, 7, 7, 7, 7, 7, 7, 8, 7, 7, 7}
	// numbers := [11]int{8, 8, 8, 8, 8, 8, 8, 7, 8, 8, 8}
	// numbers := [11]int{7, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8}
	// numbers := [11]int{8, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7}
	numbers := [11]int{7, 7, 8, 7, 7, 7, 7, 8, 7, 7, 7}
	num1 := numbers[0]
	num2 := 0
	times1 := 0
	times2 := 0

	for i := 0; i < len(numbers); i++ {
		if num1 != numbers[i] {
			num2 = numbers[i]
			times2++
		} else {
			times1++
		}
	}

	if times1 == 1 {
		fmt.Printf("\n%d is different from a series of %d\n", num1, num2)
	} else if times2 == 1 {
		fmt.Printf("\n%d is different from a series of %d\n", num2, num1)
	} else {
		fmt.Println("\nThere is a problem!")
		fmt.Print("We cannot assure that there is a series of a given number plus another single number in 'numbers' array: ")
		fmt.Println(numbers)
	}

	fmt.Printf("\ttimes1 ---> %d\n", times1)
	fmt.Printf("\ttimes2 ---> %d\n", times2)
}
