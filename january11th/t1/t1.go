package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
There is a sequence a1, a2, ..., an with some numbers.
All numbers are equal except for one. Write code to find it.
*/
func main() {
	// numbers := [11]int{7, 7, 7, 7, 7, 7, 7, 8, 7, 7, 7}
	// numbers := [11]int{8, 8, 8, 8, 8, 8, 8, 7, 8, 8, 8}
	// numbers := [11]int{7, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8}
	// numbers := [11]int{8, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7}
	// numbers := [11]int{7, 7, 8, 7, 7, 7, 7, 8, 7, 7, 7}
	// findTheDifferentOne(numbers[:])

	if len(os.Args) > 1 {
		processArguments(os.Args[1:])
	} else {
		processUserInput()
	}
}

/*
7 7 7 7 7 7 7 8 7 7 7
8 8 8 8 8 8 8 7 8 8 8
7 8 8 8 8 8 8 8 8 8 8
8 7 7 7 7 7 7 7 7 7 7
7 7 8 7 7 7 7 8 7 7 7
*/
func processArguments(s []string) {
	numbers := make([]int, len(s))

	for i, v := range s {
		numbers[i], _ = strconv.Atoi(v)
	}

	findTheDifferentOne(numbers)
}

func processUserInput() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Please type in integer values separated by single spaces: ")
	userInput, _ := reader.ReadString('\n')
	userInput = strings.Trim(userInput, "\n")
	numbers := strings.Split(userInput, " ")

	processArguments(numbers)
}

func findTheDifferentOne(numbers []int) {
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
		fmt.Printf("\nThere is a problem!\nWe cannot assure that there is a series of a given number plus another single number in given array: %v\n", numbers)
	}

	fmt.Printf("\ttimes1 ---> %d\n", times1)
	fmt.Printf("\ttimes2 ---> %d\n", times2)
}
