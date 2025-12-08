package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func readLinesNumbers(input []string) [][]int {
	lines := make([][]int, len(input)-1)
	for i, line := range input {
		if i == len(input)-1 {
			break
		}
		parts := strings.Fields(line)
		nums := make([]int, len(parts))
		for j, part := range parts {
			var n int
			fmt.Sscanf(part, "%d", &n)
			nums[j] = n
		}
		lines[i] = nums
	}
	return lines
}

func readLinesOperations(input []string) []string {
	parts := strings.Fields(input[len(input)-1])
	return parts
}

func run(input []string) int {
	numbers := readLinesNumbers(input)
	operations := readLinesOperations(input)

	total := 0
	for col, ops := range operations {

		lineSum := numbers[0][col]

		for row := 1; row < len(numbers); row++ {
			switch ops {
			case "+":
				lineSum += numbers[row][col]
			case "*":
				lineSum *= numbers[row][col]
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

	success = test([]string{"2", "3", "4", "*"}, 24) && success
	success = test([]string{"2", "3", "4", "+"}, 9) && success
	success = test([]string{"2 1", "3  2", "4   3", "*  +"}, 30) && success
	success = test([]string{"123 328  51 64 ", " 45 64  387 23 ", "  6 98  215 314", "*   +   *   +  "}, 4277556) && success

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
