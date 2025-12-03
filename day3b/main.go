package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strings"
)

func getMaxDigitAndPos(digits []int, startPos, endPos int) (int, int) {
	maxDigit := -1
	maxPos := -1
	for ind, val := range digits {
		if ind < startPos || ind >= endPos {
			continue
		}
		if val > maxDigit {
			maxDigit = val
			maxPos = ind
		}
	}
	if maxDigit == -1 || maxPos == -1 {
		panic("No max digit found")
	}
	return maxDigit, maxPos
}

func getBankMax(bank string, digUsed int) int {
	ints := []int{}
	for _, ch := range bank {
		ints = append(ints, int(ch-'0'))
	}

	digits := []int{}
	last_pos := 0
	for ind := 0; ind < digUsed; ind++ {

		endPos := len(ints) - (digUsed - (ind + 1))
		maxDigit, pos := getMaxDigitAndPos(ints, last_pos, endPos)
		last_pos = pos + 1
		digits = append(digits, maxDigit)
	}

	finalDigit := 0
	for ind, d := range digits {
		finalDigit = finalDigit + int(math.Pow10(digUsed-1-ind))*d
	}

	return finalDigit
}

func run(banks []string, digUsed int) int {
	total := 0
	for _, bank := range banks {
		bankMax := getBankMax(bank, digUsed)
		total += bankMax
	}
	return total
}

func test(banks []string, digUsed int, exp_total int) bool {

	total := run(banks, digUsed)

	if total == exp_total {
		fmt.Printf("✅Test passed: %v\n", banks)
		return true
	} else {
		fmt.Printf("❌Test failed: %v\n", banks)
		fmt.Printf("Actual total: %d\n", total)
		fmt.Printf("Expected total: %d\n", exp_total)
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

	success = test([]string{"987654321111111"}, 2, 98) && success
	success = test([]string{"811111111111119"}, 2, 89) && success
	success = test([]string{"234234234234278"}, 2, 78) && success
	success = test([]string{"818181911112111"}, 2, 92) && success
	success = test([]string{"987654321111111", "811111111111119", "234234234234278", "818181911112111"}, 2, 357) && success

	success = test([]string{"987654321111111"}, 12, 987654321111) && success
	success = test([]string{"811111111111119"}, 12, 811111111119) && success
	success = test([]string{"234234234234278"}, 12, 434234234278) && success
	success = test([]string{"818181911112111"}, 12, 888911112111) && success
	success = test([]string{"987654321111111", "811111111111119", "234234234234278", "818181911112111"}, 12, 3121910778619) && success

	if success {
		fmt.Printf("-=-=-=-=-=-=-=-=-=-=-=-=-\n✅All tests passed!\n")
	} else {
		fmt.Printf("-=-=-=-=-=-=-=-=-=-=-=-=-\n❌Some tests failed!\n")
		return
	}

	input := readInput("input.txt")

	result2 := run(input, 2)
	fmt.Printf("Result (2 digits): %d\n", result2)

	result12 := run(input, 12)
	fmt.Printf("Result (12 digits): %d\n", result12)
}
