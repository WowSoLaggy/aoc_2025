package main

import (
	"fmt"
	"os"
	"strings"
)

var cur_state = 50
var zeros = 0
var zero_passes = 0

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func rotate(diff int) {
	if diff == 0 {
		return
	}

	// cur_state [0..99]

	full_100_zero_passes := abs(diff) / 100
	rem := diff % 100

	// rem [-99..99]

	if rem == 0 {
		zero_passes += full_100_zero_passes
		if cur_state == 0 {
			zeros++
			zero_passes--
		}
		return
	}

	init_state := cur_state
	cur_state += rem

	// cur_state [-99..198]

	zero_passes += full_100_zero_passes

	if cur_state < 0 {

		cur_state += 100
		if init_state > 0 {
			zero_passes++
		}

	} else if cur_state == 0 {

		zeros++

	} else if cur_state == 100 {

		cur_state = 0
		zeros++

	} else if cur_state > 100 {

		cur_state -= 100
		zero_passes++
	}

	// cur_state [0..99]
}

func get_direction(ch rune) int {
	switch ch {
	case 'R':
		return 1
	case 'L':
		return -1
	default:
		panic("Invalid direction character")
	}
}

func get_amount(diff_str string) int {
	diff := 0
	fmt.Sscanf(diff_str[1:], "%d", &diff)
	return diff
}

func rotate_from_str(diff_str string) {
	dir_rune := rune(diff_str[0])
	dir := get_direction(dir_rune)

	amount := get_amount(diff_str)
	rotate(dir * amount)
}

func read_input(filename string) []string {
	// mock
	//return []string{"L150"}

	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	return lines
}

func run(init_value int, instructions []string) {
	cur_state = init_value
	zeros = 0
	zero_passes = 0

	for _, instr := range instructions {
		rotate_from_str(instr)
	}
}

func test(init_value int, instructions []string, exp_state, exp_zeros, exp_zeros_and_passes int) {

	run(init_value, instructions)

	zeros_and_passes := zeros + zero_passes
	if cur_state != exp_state || zeros != exp_zeros || zeros_and_passes != exp_zeros_and_passes {
		fmt.Printf("Test failed for: %d --> %v\n", init_value, instructions)
		fmt.Printf("Expected state: %d, got: %d\n", exp_state, cur_state)
		fmt.Printf("Expected zeros: %d, got: %d\n", exp_zeros, zeros)
		fmt.Printf("Expected zeros_and_passes: %d, got: %d\n", exp_zeros_and_passes, zeros_and_passes)
		fmt.Printf("-=-=-=-=-=-=-=-=-\n")
	}
}

func main() {

	// Run tests

	test(50, []string{"R49"}, 99, 0, 0)
	test(50, []string{"R50"}, 0, 1, 1)
	test(50, []string{"R51"}, 1, 0, 1)

	test(50, []string{"L49"}, 1, 0, 0)
	test(50, []string{"L50"}, 0, 1, 1)
	test(50, []string{"L51"}, 99, 0, 1)

	test(0, []string{"R99"}, 99, 0, 0)
	test(0, []string{"R100"}, 0, 1, 1)
	test(0, []string{"R101"}, 1, 0, 1)

	test(0, []string{"L99"}, 1, 0, 0)
	test(0, []string{"L100"}, 0, 1, 1)
	test(0, []string{"L101"}, 99, 0, 1)

	test(50, []string{"R150"}, 0, 1, 2)
	test(50, []string{"L150"}, 0, 1, 2)
	test(50, []string{"R250"}, 0, 1, 3)
	test(50, []string{"L250"}, 0, 1, 3)
	test(50, []string{"R1000"}, 50, 0, 10)
	test(50, []string{"R1001"}, 51, 0, 10)

	// Run PRD

	filename := "input.txt"
	instructions := read_input(filename)

	run(50, instructions)

	fmt.Printf("Final state: %d\n", cur_state)
	fmt.Printf("Number of times at zero: %d\n", zeros)
	fmt.Printf("Number of times at or passed zero: %d\n", zeros+zero_passes)
}
