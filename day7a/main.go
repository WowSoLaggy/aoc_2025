package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

// 0 is emptiness
// 1 is a beam
// 2 is a splitter

func isSplitter(input []string, x, y int) bool {
	return input[y][x] == '^'
}

func isBeamAbove(lines [][]int, x, y int) bool {
	if y <= 0 {
		return false
	}
	return lines[y-1][x] == 1
}

func run(input []string) int {

	lines := [][]int{}

	splits := 0
	for y, line := range input {
		lineInt := make([]int, len(line))

		if y == 0 {
			sPos := strings.Index(line, "S")
			lineInt[sPos] = 1
		} else {
			for x := range line {
				if lineInt[x] == 1 {
					continue
				}
				if isBeamAbove(lines, x, y) {
					if isSplitter(input, x, y) {
						lineInt[x-1] = 1
						lineInt[x+1] = 1
						lineInt[x] = 2
						splits++
					} else {
						lineInt[x] = 1
					}
				}
			}
		}
		lines = append(lines, lineInt)
	}

	return splits
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

	success = test([]string{
		".......S.......",
		"...............",
		".......^.......",
		"...............",
		"......^.^......",
		"...............",
		".....^.^.^.....",
		"...............",
		"....^.^...^....",
		"...............",
		"...^.^...^.^...",
		"...............",
		"..^...^.....^..",
		"...............",
		".^.^.^.^.^...^.",
		"...............",
	}, 21) && success

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
