package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func checkGridConsistent(grid []string) bool {
	if len(grid) == 0 {
		return true
	}
	width := len(grid[0])
	for _, row := range grid {
		if len(row) != width {
			return false
		}
	}
	return true
}

func hasRoll(grid []string, x, y, width, height int) bool {
	if x < 0 || x >= width || y < 0 || y >= height {
		return false
	}
	return grid[y][x] == '@'
}

func getNeighbors(grid []string, x, y, width, height int) int {
	directions := [][2]int{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1}}

	neighbors := 0
	for dir := range directions {
		dx := directions[dir][0]
		dy := directions[dir][1]
		nx := x + dx
		ny := y + dy
		if hasRoll(grid, nx, ny, width, height) {
			neighbors++
		}
	}

	return neighbors
}

func run(grid []string) int {
	if !checkGridConsistent(grid) {
		panic("Inconsistent grid")
	}

	width := len(grid[0])
	height := len(grid)
	if width == 0 || height == 0 {
		return 0
	}

	rolls := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if !hasRoll(grid, x, y, width, height) {
				continue
			}

			neighbors := getNeighbors(grid, x, y, width, height)
			if neighbors < 4 {
				rolls++
			}
		}
	}

	return rolls
}

func test(grid []string, exp_rolls int) bool {

	rolls := run(grid)

	if rolls == exp_rolls {
		fmt.Printf("✅Test passed: %v\n", grid)
		return true
	} else {
		fmt.Printf("❌Test failed: %v\n", grid)
		fmt.Printf("Actual total: %d\n", rolls)
		fmt.Printf("Expected total: %d\n", exp_rolls)
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

	success = test([]string{""}, 0) && success
	success = test([]string{"", ""}, 0) && success

	success = test([]string{"."}, 0) && success
	success = test([]string{"..", ".."}, 0) && success
	success = test([]string{"..."}, 0) && success
	success = test([]string{"...", "..."}, 0) && success
	success = test([]string{"...", "...", "..."}, 0) && success

	success = test([]string{"@"}, 1) && success
	success = test([]string{"@@"}, 2) && success
	success = test([]string{"@.", ".."}, 1) && success
	success = test([]string{".@", ".."}, 1) && success
	success = test([]string{"..", "@."}, 1) && success
	success = test([]string{"..", ".@"}, 1) && success
	success = test([]string{"@@", ".."}, 2) && success
	success = test([]string{"@.", "@."}, 2) && success
	success = test([]string{"@.", ".@"}, 2) && success
	success = test([]string{".@", "@."}, 2) && success
	success = test([]string{".@", ".@"}, 2) && success
	success = test([]string{"..", "@@"}, 2) && success
	success = test([]string{"@@", "@."}, 3) && success
	success = test([]string{"@@", ".@"}, 3) && success
	success = test([]string{"@.", "@@"}, 3) && success

	success = test([]string{".@.", "...", "..."}, 1) && success
	success = test([]string{"...", ".@.", "..."}, 1) && success
	success = test([]string{"...", "...", ".@."}, 1) && success
	success = test([]string{"...", "@..", "..."}, 1) && success
	success = test([]string{"...", "..@", "..."}, 1) && success
	success = test([]string{"...", "@.@", "..."}, 2) && success
	success = test([]string{".@.", "@@@", ".@."}, 4) && success

	success = test([]string{
		"..@@.@@@@.",
		"@@@.@.@.@@",
		"@@@@@.@.@@",
		"@.@@@@..@.",
		"@@.@@@@.@@",
		".@@@@@@@.@",
		".@.@.@.@@@",
		"@.@@@.@@@@",
		".@@@@@@@@.",
		"@.@.@@@.@.",
	}, 13) && success

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
