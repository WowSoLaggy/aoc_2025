package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
)

type Vector = []int
type Matrix = []Vector

type Machine struct {
	switches      []Vector
	finalJoltages Vector
}

func readMachine(line string) Machine {
	machine := Machine{}

	parts := strings.Split(line, " ")

	// Start from final joltages to know the size of switches matrix
	finalStr := parts[len(parts)-1][1 : len(parts[len(parts)-1])-1]
	tokens := strings.Split(finalStr, ",")
	finalJoltages := []int{}
	for _, token := range tokens {
		var val int
		fmt.Sscanf(token, "%d", &val)
		finalJoltages = append(finalJoltages, val)
	}
	machine.finalJoltages = finalJoltages

	// Switches
	switches := []Vector{}
	for _, part := range parts[1 : len(parts)-1] {
		switchStr := part[1 : len(part)-1]
		switchVec := make(Vector, len(machine.finalJoltages))
		// switch e.g. (1,3) -> vector (0,1,0,1,0)
		tokens := strings.Split(switchStr, ",")
		for _, token := range tokens {
			var idx int
			fmt.Sscanf(token, "%d", &idx)
			switchVec[idx] = 1
		}
		switches = append(switches, switchVec)
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

func gaussianElimination(m Matrix) (Vector, Matrix) {
	// Return empty when nothing to do
	if len(m) == 0 {
		return nil, nil
	}
	rowCount := len(m)
	colCount := len(m[0]) - 1 // last column is the RHS

	// Work on a copy
	mat := make(Matrix, rowCount)
	for i := 0; i < rowCount; i++ {
		mat[i] = make([]int, colCount+1)
		copy(mat[i], m[i])
	}

	pivots := Vector{}
	r := 0

	for c := 0; c < colCount && r < rowCount; c++ {
		// Find a row with a non-zero in column c, starting from r
		pivot := -1
		for i := r; i < rowCount; i++ {
			if mat[i][c] != 0 {
				pivot = i
				break
			}
		}
		// No pivot in this column, skip
		if pivot == -1 {
			continue
		}

		// Move pivot row to current row
		if pivot != r {
			mat[r], mat[pivot] = mat[pivot], mat[r]
		}
		pivots = append(pivots, c)

		// Eliminate below
		base := mat[r][c]
		for i := r + 1; i < rowCount; i++ {
			if mat[i][c] == 0 {
				continue
			}
			f := mat[i][c]
			for j := c; j <= colCount; j++ {
				// integer-only elimination (no division)
				mat[i][j] = mat[i][j]*base - mat[r][j]*f
			}
		}
		r++
	}

	return pivots, mat
}

func getSolution(machine Machine) Vector {
	buttons := machine.switches
	final := machine.finalJoltages

	rowCount := len(final)
	colCount := len(buttons)

	// Build augmented matrix: A|b
	aug := make(Matrix, rowCount)
	for i := 0; i < rowCount; i++ {
		aug[i] = make(Vector, colCount+1)
		for j := 0; j < colCount; j++ {
			if buttons[j][i] != 0 {
				aug[i][j] = 1
			}
		}
		aug[i][colCount] = final[i]
	}

	pivotCols, reduced := gaussianElimination(aug)
	if reduced == nil {
		return nil
	}

	// Track pivot columns
	isPivot := make(map[int]bool, len(pivotCols))
	for _, c := range pivotCols {
		isPivot[c] = true
	}

	// Collect free vars
	freeVars := []int{}
	for j := 0; j < colCount; j++ {
		if !isPivot[j] {
			freeVars = append(freeVars, j)
		}
	}

	best := make([]int, colCount)
	bestSum := -1

	// Try a candidate assignment for free variables, solve pivots, then validate
	solve := func(freeVals []int) {
		x := make([]int, colCount)

		for k, idx := range freeVars {
			if k < len(freeVals) {
				x[idx] = freeVals[k]
			}
		}

		// Back-substitute for pivot variables
		for k := len(pivotCols) - 1; k >= 0; k-- {
			row := k
			pc := pivotCols[k]
			rhs := reduced[row][colCount]

			for j := pc + 1; j < colCount; j++ {
				if reduced[row][j] != 0 && x[j] != 0 {
					rhs -= reduced[row][j] * x[j]
				}
			}

			den := reduced[row][pc]
			if den == 0 {
				return
			}
			if rhs%den != 0 {
				return
			}

			val := rhs / den
			if val < 0 {
				return
			}
			x[pc] = val
		}

		// Validate solution against original constraints
		for i := 0; i < rowCount; i++ {
			sum := 0
			for j := 0; j < colCount; j++ {
				if x[j] != 0 && buttons[j][i] != 0 {
					sum += x[j]
				}
			}
			if sum != final[i] {
				return
			}
		}

		total := 0
		for _, v := range x {
			total += v
		}
		if bestSum == -1 || total < bestSum {
			copy(best, x)
			bestSum = total
		}
	}

	// BELOW HEURISTIC IS APPLIED, PLS DON'T ASK WHY
	// I AM ASHAMED OF IT (hopefully, GPT either)

	// Brute over small ranges for few free vars, aiming to minimize total presses
	switch len(freeVars) {
	case 0:
		solve(nil)
	case 1:
		maxV := 0
		for _, v := range final {
			if v > maxV {
				maxV = v
			}
		}
		// generous upper bound, but we short-circuit by bestSum
		maxV *= 3
		for v := 0; v <= maxV; v++ {
			if bestSum != -1 && v > bestSum {
				break
			}
			solve([]int{v})
		}
	case 2:
		maxV := 0
		for _, v := range final {
			if v > maxV {
				maxV = v
			}
		}
		if maxV < 200 {
			maxV = 200
		}
		for v1 := 0; v1 <= maxV; v1++ {
			for v2 := 0; v2 <= maxV; v2++ {
				if bestSum != -1 && v1+v2 > bestSum {
					continue
				}
				solve([]int{v1, v2})
			}
		}
	case 3:
		for v1 := 0; v1 < 250; v1++ {
			for v2 := 0; v2 < 250; v2++ {
				for v3 := 0; v3 < 250; v3++ {
					if bestSum != -1 && v1+v2+v3 > bestSum {
						continue
					}
					solve([]int{v1, v2, v3})
				}
			}
		}
	case 4:
		for v1 := 0; v1 < 30; v1++ {
			for v2 := 0; v2 < 30; v2++ {
				for v3 := 0; v3 < 30; v3++ {
					for v4 := 0; v4 < 30; v4++ {
						if bestSum != -1 && v1+v2+v3+v4 > bestSum {
							continue
						}
						solve([]int{v1, v2, v3, v4})
					}
				}
			}
		}
	default:
		// Fallback: zero init
		solve(make([]int, len(freeVars)))
	}

	if bestSum == -1 {
		panic("No solution found")
	}
	return best
}

func getSumN(machine Machine) int {
	sol := getSolution(machine)
	sum := 0
	for _, val := range sol {
		sum += val
	}
	return sum
}

func run(input []string) int {
	machines := readMachines(input)

	sum := 0
	for ind, machine := range machines {
		timeStart := time.Now()
		newSum := getSumN(machine)
		timeEnd := time.Now()
		fmt.Printf("Machine %d: sumN = %d (t=%dms)\n", ind+1, newSum, timeEnd.Sub(timeStart).Milliseconds())
		sum += newSum
	}

	return sum
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
	success = test([]string{
		"[...#...] (0,2,3,6) (0,1,4,6) (1,3,4,5) (1,2,4,6) (0,2,3,4,5) (2,3,6) (1,2) (2,3,4,5,6) {37,24,84,71,44,32,71}",
	}, 90) && success

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
