package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func getMaxDigitAndPos(digits []int) (int, int) {
	maxDigit := -1
	maxPos := -1
	for ind, val := range digits {
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

func getBankMax(bank string) int {
	ints := []int{}
	for _, ch := range bank {
		ints = append(ints, int(ch-'0'))
	}

	maxDigit, pos := getMaxDigitAndPos(ints[:len(ints)-1])
	maxDigit2, _ := getMaxDigitAndPos(ints[pos+1:])

	return maxDigit*10 + maxDigit2
}

func run(banks []string) int {
	total := 0
	for _, bank := range banks {
		bankMax := getBankMax(bank)
		total += bankMax
	}
	return total
}

func test(banks []string, exp_total int) bool {

	total := run(banks)

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

	success = test([]string{"987654321111111"}, 98) && success
	success = test([]string{"811111111111119"}, 89) && success
	success = test([]string{"234234234234278"}, 78) && success
	success = test([]string{"818181911112111"}, 92) && success
	success = test([]string{"987654321111111", "811111111111119", "234234234234278", "818181911112111"}, 357) && success

	if success {
		fmt.Printf("-=-=-=-=-=-=-=-=-=-=-=-=-\n✅All tests passed!\n")
	} else {
		fmt.Printf("-=-=-=-=-=-=-=-=-=-=-=-=-\n❌Some tests failed!\n")
		return
	}

	input := readInput("input.txt")
	fmt.Printf("Input: %v\n", input)

	result := run(input)
	fmt.Printf("Result: %d\n", result)
}
