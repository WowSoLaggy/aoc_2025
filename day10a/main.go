package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

type Machine struct {
	finalState int
	switches   []int
}

func readMachine(line string) Machine {
	machine := Machine{}

	parts := strings.Split(line, " ")

	// final state
	finalStateStr := parts[0][1 : len(parts[0])-1]
	finalState := 0
	for i, ch := range finalStateStr {
		if ch == '#' {
			finalState |= (1 << i)
		}
	}
	machine.finalState = finalState

	// switches
	switches := []int{}
	for _, part := range parts[1 : len(parts)-1] {
		switchStr := part[1 : len(part)-1]
		switchState := 0
		if switchStr != "" {
			switchIndices := strings.Split(switchStr, ",")
			for _, idxStr := range switchIndices {
				var idx int
				fmt.Sscanf(idxStr, "%d", &idx)
				switchState |= (1 << idx)
			}
		}
		switches = append(switches, switchState)
	}
	machine.switches = switches

	return machine
}

func readMachines(input []string) []Machine {
	machines := []Machine{}
	for _, line := range input {
		machine := readMachine(line)
		machines = append(machines, machine)
	}
	return machines
}

func getStateAfterSeq(seq []int, machine Machine) int {
	state := 0

	for _, btnInd := range seq {
		switchState := machine.switches[btnInd]
		state ^= switchState
	}

	return state
}

func isSeqValid(seq []int, machine Machine) bool {
	finalState := machine.finalState
	stateAfterSeq := getStateAfterSeq(seq, machine)
	return stateAfterSeq == finalState
}

func continuesSeqs(seqs [][]int, numSwitches int, seqMap map[string]int, stateMap map[int]bool, machine Machine) [][]int {
	newSeqs := [][]int{}
	for _, seq := range seqs {

		key := fmt.Sprint(seq)
		seqState, exists := seqMap[key]
		if !exists {
			panic("seq not in map")
		}

		for switchInd := 0; switchInd < numSwitches; switchInd++ {

			if len(seq) > 0 && seq[len(seq)-1] == switchInd {
				continue // never repeat the same switch
			}

			newState := seqState ^ machine.switches[switchInd]
			if _, exist := stateMap[newState]; exist {
				continue // already visited this state
			}

			newSeq := append([]int{}, seq...)
			newSeq = append(newSeq, switchInd)

			// store to maps
			seqMap[fmt.Sprint(newSeq)] = newState
			stateMap[newState] = true

			newSeqs = append(newSeqs, newSeq)
		}
	}
	return newSeqs
}

func getMinSwitches(machine Machine) int {
	if machine.finalState == 0 {
		return 0
	}

	seqs := make([][]int, 0)
	seqs = append(seqs, []int{})

	seqMap := make(map[string]int)
	seqMap["[]"] = 0

	stateMap := make(map[int]bool)
	stateMap[0] = true

	curLength := 0
	for {
		curLength++
		seqs = continuesSeqs(seqs, len(machine.switches), seqMap, stateMap, machine)

		for _, seq := range seqs {
			if isSeqValid(seq, machine) {
				return len(seq)
			}
		}
	}
}

func run(input []string) int {
	machines := readMachines(input)

	totalSwitches := 0
	for ind, machine := range machines {
		minSwitches := getMinSwitches(machine)
		totalSwitches += minSwitches
		fmt.Printf("Machine %d: min switches = %d\n", ind+1, minSwitches)
	}

	return totalSwitches
}

func test(input []string, exp_output int) bool {

	fmt.Printf("Test started: %v\n", input)
	output := run(input)

	if output == exp_output {
		fmt.Println("✅Test passed")
		return true
	} else {
		fmt.Println("❌Test failed:")
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

	success = test([]string{"[#] (0) {3}"}, 1) && success
	success = test([]string{"[.] (0) {3}"}, 0) && success
	success = test([]string{"[##] (0) (1) {3}"}, 2) && success
	success = test([]string{"[##] (0) (0,1) (1) {3}"}, 1) && success
	success = test([]string{"[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}"}, 2) && success
	success = test([]string{"[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}"}, 3) && success
	success = test([]string{"[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}"}, 2) && success
	success = test([]string{
		"[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}",
		"[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}",
		"[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}",
	}, 7) && success

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
