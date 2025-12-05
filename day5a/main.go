package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

type Range struct {
	start int
	end   int
}

func (r *Range) Contains(value int) bool {
	return value >= r.start && value <= r.end
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

func getIds(input []string) []int {
	var ids []int
	readingIds := false
	for _, line := range input {
		if line == "" {
			readingIds = true
			continue
		}
		if readingIds {
			var id int
			fmt.Sscanf(line, "%d", &id)
			ids = append(ids, id)
		}
	}
	return ids
}

func run(input []string) int {
	ranges := getRanges(input)
	ids := getIds(input)

	count := 0
	for _, id := range ids {
		for _, r := range ranges {
			if r.Contains(id) {
				count++
				break
			}
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

	success = test([]string{"3-5", "", "1"}, 0) && success
	success = test([]string{"3-5", "", "3"}, 1) && success
	success = test([]string{"3-5", "", "4"}, 1) && success
	success = test([]string{"3-5", "", "5"}, 1) && success
	success = test([]string{"3-5", "", "6"}, 0) && success
	success = test([]string{"3-5", "7-9", "", "6"}, 0) && success
	success = test([]string{"3-5", "7-9", "", "8"}, 1) && success
	success = test([]string{"3-5", "7-9", "8-9", "", "8"}, 1) && success

	success = test([]string{
		"3-5",
		"10-14",
		"16-20",
		"12-18",
		"",
		"1",
		"5",
		"8",
		"11",
		"17",
		"32",
	}, 3) && success

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
