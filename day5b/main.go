package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Range struct {
	start int
	end   int
}

func getRanges(input []string) []Range {
	var ranges []Range
	for _, line := range input {
		if line == "" {
			break
		}
		var start, end int
		fmt.Sscanf(line, "%d-%d", &start, &end)
		ranges = append(ranges, Range{start: start, end: end})
	}
	return ranges
}

func getSwitches(ranges []Range) map[int]int {
	switches := make(map[int]int)
	for _, r := range ranges {
		switches[r.start]++
		switches[r.end+1]--
	}
	return switches
}

func getSortedKeys(m map[int]int) []int {
	keys := []int{}
	for k := range m {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	return keys
}

func run(input []string) int {
	ranges := getRanges(input)
	switches := getSwitches(ranges)
	sortedKeys := getSortedKeys(switches)

	count := 0
	curSwitch := 0 // current number of active ranges
	prevKey := sortedKeys[0]
	for _, key := range sortedKeys {
		switch_of_current_key := switches[key]

		if curSwitch == 0 && switch_of_current_key > 0 {
			prevKey = key
		}

		curSwitch += switch_of_current_key

		if curSwitch < 0 {
			panic("curSwitch < 0")
		}

		if curSwitch == 0 {
			count += key - prevKey
		}
	}

	return count
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

	success = test([]string{"3-5"}, 3) && success
	success = test([]string{"3-5", "7-9"}, 6) && success
	success = test([]string{"3-5", "7-9", "8-9"}, 6) && success

	success = test([]string{
		"3-5",
		"10-14",
		"16-20",
		"12-18",
	}, 14) && success

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
