package main

import (
	"fmt"
	"os"
	"strings"
)

var cur_state = 50
var zeros = 0

func rotate(diff int) {
	diff = diff % 100
	cur_state += diff
	cur_state = (cur_state + 100) % 100
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
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	return lines
}

func check_zero() {
	if cur_state == 0 {
		zeros++
	}
}

func main() {
	filename := "input.txt"
	instructions := read_input(filename)

	for _, instr := range instructions {
		rotate_from_str(instr)
		check_zero()
	}

	fmt.Printf("Final state: %d\n", cur_state)
	fmt.Printf("Number of times at zero: %d\n", zeros)
}
