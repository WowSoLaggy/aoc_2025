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

func getStartCoords(input []string) (int, int) {
	for y, line := range input {
		sPos := strings.Index(line, "S")
		if sPos != -1 {
			return sPos, y
		}
	}
	return -1, -1
}

func isEmpty(input []string, x, y int) bool {
	if y < 0 || y >= len(input) || x < 0 || x >= len(input[0]) {
		return true
	}
	return input[y][x] == '.'
}
func isSplitter(input []string, x, y int) bool {
	if y < 0 || y >= len(input) || x < 0 || x >= len(input[0]) {
		return false
	}
	return input[y][x] == '^'
}

var cacheSplits [][]int

func makeTurnFrom(input []string, curX, curY int) int {
	if curY >= len(input) {
		return 0
	}

	if cacheSplits[curY][curX] != 0 {
		return cacheSplits[curY][curX]
	}

	if isEmpty(input, curX, curY+1) {
		splits := makeTurnFrom(input, curX, curY+1)
		cacheSplits[curY][curX] = splits
		return splits
	} else if isSplitter(input, curX, curY+1) {
		splitsA := makeTurnFrom(input, curX-1, curY+1)
		splitsB := makeTurnFrom(input, curX+1, curY+1)
		splits := splitsA + splitsB + 1
		cacheSplits[curY][curX] = splits
		return splits
	} else {
		panic("Beam blocked")
	}
}

func run(input []string) int {

	cacheSplits = make([][]int, len(input))
	for i := range cacheSplits {
		cacheSplits[i] = make([]int, len(input[0]))
	}

	curX, curY := getStartCoords(input)
	splits := makeTurnFrom(input, curX, curY)

	return splits + 1
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
		".S.",
		"...",
		"...",
		"...",
	}, 1) && success
	success = test([]string{
		".S.",
		"...",
		".^.",
		"...",
	}, 2) && success
	success = test([]string{
		".S.",
		"...",
		"^..",
		"...",
	}, 1) && success
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
	}, 40) && success

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
