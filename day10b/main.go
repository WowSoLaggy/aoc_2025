package main

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Machine struct {
	switches      []int
	finalJoltages []int
}

type SeqCost struct {
	seq  []int
	cost int
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func readMachine(line string) Machine {
	machine := Machine{}

	parts := strings.Split(line, " ")

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

	// final joltages in format: {2,3,5,7} -> store just as an int array
	finalStr := parts[len(parts)-1][1 : len(parts[len(parts)-1])-1]
	tokens := strings.Split(finalStr, ",")
	finalJoltages := []int{}
	for _, token := range tokens {
		var val int
		fmt.Sscanf(token, "%d", &val)
		finalJoltages = append(finalJoltages, val)
	}
	machine.finalJoltages = finalJoltages

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

func key(seq []int) string {
	keyStr := ""
	for i, val := range seq {
		if i > 0 {
			keyStr += ","
		}
		keyStr += fmt.Sprintf("%d", val)
	}
	return keyStr
}

func applySwitch(currentJoltage []int, switchState int) []int {
	for i := 0; i < len(currentJoltage); i++ {
		if (switchState & (1 << i)) != 0 {
			currentJoltage[i]++
		}
	}
	return currentJoltage
}

func getJoltagesAfterSeq(seq []int, machine Machine, cache map[string][]int) []int {
	currentJoltage := make([]int, len(machine.finalJoltages))
	if len(seq) == 0 {
		return currentJoltage
	}

	// Check cache first
	seqKey := key(seq)
	cachedState, exists := cache[seqKey]
	if exists {
		return cachedState
	}

	// State for previous seq shall always be already cached
	if len(seq) > 0 {
		prevSeq := seq[:len(seq)-1]
		prevKey := key(prevSeq)
		cachedState, exists := cache[prevKey]
		if exists {
			copy(currentJoltage, cachedState)
		} else {
			panic("Previous sequence state not found in cache")
		}
	}

	// Apply only the last switch
	currentJoltage = applySwitch(currentJoltage, machine.switches[seq[len(seq)-1]])

	// Cache the current state
	if cache != nil {
		cache[seqKey] = make([]int, len(currentJoltage))
		copy(cache[seqKey], currentJoltage)
	}

	return currentJoltage
}

func isSeqValid(seq []int, machine Machine, cache map[string][]int) bool {
	joltagesAfterSeq := getJoltagesAfterSeq(seq, machine, cache)
	return slices.Equal(joltagesAfterSeq, machine.finalJoltages)
}

func getSeqCost(seq []int, machine Machine, cache map[string][]int) int {
	joltagesAfterSeq := getJoltagesAfterSeq(seq, machine, cache)
	cost := 0
	for i := 0; i < len(joltagesAfterSeq); i++ {
		cost += abs(joltagesAfterSeq[i] - machine.finalJoltages[i])
	}
	return cost
}

func checkSeqExceedsFinalJoltages(seq []int, machine Machine, cache map[string][]int) bool {
	joltagesAfterSeq := getJoltagesAfterSeq(seq, machine, cache)
	for i := 0; i < len(joltagesAfterSeq); i++ {
		if joltagesAfterSeq[i] > machine.finalJoltages[i] {
			return true
		}
	}
	return false
}

func continuesSeq(seq []int, numSwitches int, machine Machine, cache map[string][]int) [][]int {
	newSeqs := [][]int{}
	for switchInd := 0; switchInd < numSwitches; switchInd++ {
		newSeq := append([]int{}, seq...)
		newSeq = append(newSeq, switchInd)

		if checkSeqExceedsFinalJoltages(newSeq, machine, cache) {
			continue
		}

		newSeqs = append(newSeqs, newSeq)
	}
	return newSeqs
}

func checkSwitchesAreRequired(machine Machine) bool {
	for _, finalJoltage := range machine.finalJoltages {
		if finalJoltage != 0 {
			return true
		}
	}
	return false
}

func sort(seqs []SeqCost) {
	slices.SortFunc(seqs, func(a, b SeqCost) int {
		return a.cost - b.cost
	})
}

func getMinSwitches(machine Machine) int {
	if !checkSwitchesAreRequired(machine) {
		return 0
	}

	seqs := make([]SeqCost, 0)
	seqs = append(seqs, SeqCost{seq: []int{}, cost: getSeqCost([]int{}, machine, nil)})

	cache := make(map[string][]int)
	cache[""] = make([]int, len(machine.finalJoltages))

	iteration := 0
	for {
		iteration++

		// we have a list of candidate seqs
		// sort them by cost (the less the better) and proceed with the best one
		sort(seqs)
		bestSeq := seqs[0]
		if bestSeq.cost == 0 {
			return len(bestSeq.seq)
		}

		if iteration > 1000 && iteration%1000 == 0 {
			fmt.Printf("Iteration %d, best cost: %d\n", iteration, bestSeq.cost)
		}

		newSeqs := continuesSeq(bestSeq.seq, len(machine.switches), machine, cache)
		seqs = seqs[1:] // remove the best seq
		for _, newSeq := range newSeqs {
			newCost := getSeqCost(newSeq, machine, cache)
			seqs = append(seqs, SeqCost{seq: newSeq, cost: newCost})
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

	success = test([]string{"[#] (0) {0}"}, 0) && success
	success = test([]string{"[#] (0) {1}"}, 1) && success
	success = test([]string{"[#] (0) {2}"}, 2) && success
	success = test([]string{"[#] (0) {3}"}, 3) && success
	success = test([]string{"[##] (0) (1) {2,3}"}, 5) && success
	success = test([]string{"[##] (0) (0,1) (1) {2,3}"}, 3) && success
	success = test([]string{"[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}"}, 10) && success
	success = test([]string{"[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}"}, 12) && success
	success = test([]string{"[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}"}, 11) && success
	success = test([]string{
		"[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}",
		"[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}",
		"[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}",
	}, 33) && success

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
