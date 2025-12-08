package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func checkInputConsistent(input []string) {
	if len(input) <= 0 {
		panic("Input is empty")
	}
	numCols := len(input[0])
	for i, line := range input {
		if len(line) != numCols {
			panic(fmt.Sprintf("Inconsistent line length at line %d: expected %d, got %d", i, numCols, len(line)))
		}
	}
}

func isEmptyCol(input []string, col int) bool {
	for row := 0; row < len(input); row++ {
		char := input[row][col]
		if char != ' ' {
			return false
		}
	}
	return true
}

func readNumbersFromColumns(input []string) [][]int {
	numbers := [][]int{}
	numbers = append(numbers, []int{})

	numCols := len(input[0])
	for col := 0; col < numCols; col++ {

		if isEmptyCol(input, col) {
			numbers = append(numbers, []int{})
			continue
		}

		curNumberString := ""
		for row := 0; row < len(input)-1; row++ {
			char := input[row][col]
			curNumberString += string(char)
		}
		curNumber := 0
		fmt.Sscanf(curNumberString, "%d", &curNumber)
		numbers[len(numbers)-1] = append(numbers[len(numbers)-1], curNumber)
	}

	return numbers
}

func readLinesOperations(input []string) []string {
	parts := strings.Fields(input[len(input)-1])
	return parts
}

func run(input []string) int {

	checkInputConsistent(input)

	numbers := readNumbersFromColumns(input)
	operations := readLinesOperations(input)

	total := 0
	for opInd, ops := range operations {

		lineSum := numbers[opInd][0]

		for row := 1; row < len(numbers[opInd]); row++ {
			switch ops {
			case "+":
				lineSum += numbers[opInd][row]
			case "*":
				lineSum *= numbers[opInd][row]
			}
		}

		total += lineSum
	}

	return total
}

func test(input []string, exp_output int) bool {

	output := run(input)

	if output == exp_output {
		fmt.Printf("✅Test passed: %v\n", input)
		return true
	} else {
		fmt.Printf("❌Test failed: %v\n", input)
		fmt.Printf("Actual output: %d\n", output)
		fmt.Printf("Expected output: %d\n", exp_output)
		return false
	}
}

func readInput(filename string) []string {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	lines := strings.Split(string(data), "\n")
	return lines
}

func main() {
	success := true

	success = test([]string{"1", "2", "+"}, 12) && success
	success = test([]string{"13", "2 ", "+ "}, 15) && success
	success = test([]string{"13", " 2", "+ "}, 33) && success
	success = test([]string{"123 328  51 64 ", " 45 64  387 23 ", "  6 98  215 314", "*   +   *   +  "}, 3263827) && success

	if success {
		fmt.Printf("-=-=-=-=-=-=-=-=-=-=-=-=-\n✅All tests passed!\n")
	} else {
		fmt.Printf("-=-=-=-=-=-=-=-=-=-=-=-=-\n❌Some tests failed!\n")
		return
	}

	input := readInput("input.txt")

	result := run(input)
	fmt.Printf("Result: %d\n", result)
}
