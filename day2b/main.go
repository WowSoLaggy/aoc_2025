package main

import (
	"fmt"
	"os"
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
	idStr := strconv.Itoa(id)
	idLen := len(idStr)
	halfLen := idLen / 2

	for patternLen := 1; patternLen <= halfLen; patternLen++ {
		if idLen%patternLen != 0 {
			continue
		}

		pattern := idStr[:patternLen]
		patternRepeats := idLen / patternLen
		expectedStr := strings.Repeat(pattern, patternRepeats)
		if expectedStr == idStr {
			return false
		}
	}

	return true
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

func test(str string, exp_sum int) bool {

	sum := run(str)

	if sum == exp_sum {
		fmt.Printf("✅Test passed: %v\n", str)
		return true
	} else {
		fmt.Printf("❌Test failed: %v\n", str)
		fmt.Printf("Actual sum: %d\n", sum)
		fmt.Printf("Expected sum: %d\n", exp_sum)
		return false
	}
}

func readInput(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(data))
}

func main() {
	success := true

	success = test("1-1", 0) && success
	success = test("01-01", 0) && success
	success = test("11-22", 33) && success
	success = test("95-115", 210) && success
	success = test("998-1012", 2009) && success
	success = test("1188511880-1188511890", 1188511885) && success
	success = test("222220-222224", 222222) && success
	success = test("1698522-1698528", 0) && success
	success = test("446443-446449", 446446) && success
	success = test("38593856-38593862", 38593859) && success
	success = test("565653-565659", 565656) && success
	success = test("824824821-824824827", 824824824) && success
	success = test("2121212118-2121212124", 2121212121) && success

	success = test("38593856-38593862,565653-565659", 38593859+565656) && success
	success = test("446443-446449,38593856-38593862,565653-565659", 446446+38593859+565656) && success

	success = test("11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,"+
		"38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124", 4174379265) && success

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
