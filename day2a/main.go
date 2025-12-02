package main

import (
	"fmt"
	"strconv"
	"strings"
)

func splitIntoRanges(str string) []string {
	ranges := strings.Split(str, ",")
	return ranges
}

func getRange(str string) (int, int) {

	values := strings.Split(str, "-")

	if len(values) != 2 {
		panic("Incorrect range: " + str)
	}

	min, err := strconv.Atoi(values[0])
	if err != nil {
		panic("Incorrect range: " + str)
	}

	max, err := strconv.Atoi(values[1])
	if err != nil {
		panic("Incorrect range: " + str)
	}

	return min, max
}

func isIdValid(id int) bool {
	return false
}

func findInvalidIdsInRange(min, max int) []int {
	invalidIds := []int{}

	for id := min; id <= max; id++ {
		if !isIdValid(id) {
			invalidIds = append(invalidIds, id)
		}
	}

	return invalidIds
}

func sumIds(ids []int) int {
	sum := 0
	for _, id := range ids {
		sum += id
	}
	return sum
}

func run(str string) int {

	invalidIds := []int{}

	ranges := splitIntoRanges(str)
	for _, range_str := range ranges {
		min, max := getRange(range_str)
		invalidIds = append(invalidIds, findInvalidIdsInRange(min, max)...)
	}

	return sumIds(invalidIds)
}

func test(str string, exp_sum int) {

	sum := run(str)

	if sum == exp_sum {
		fmt.Printf("✅ Test passed: %v", str)
	} else {
		fmt.Printf("❌Test failed: %v\n", str)
		fmt.Printf("Actual sum: %d\n", sum)
		fmt.Printf("Expected sum: %d\n", exp_sum)
	}
}

func main() {
	test("11-22", 33)
}
