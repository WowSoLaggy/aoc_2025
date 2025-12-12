package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const PIECESIZE = 3

type Piece struct {
	id   int
	data [][]int
}

type Field struct {
	width, height int
	piecesUsed    []int
}

func readPieces(input []string) []Piece {
	pieces := []Piece{}

	// piece starts with "id:"
	for i := 0; i < len(input); {
		line := input[i]
		var id int
		_, err := fmt.Sscanf(line, "%d:", &id)
		if err != nil {
			break // reached end of pieces
		}

		piece := Piece{
			id: id,
		}
		piece.data = make([][]int, PIECESIZE)
		for r := 0; r < PIECESIZE; r++ {
			pieceLine := input[i+1+r]
			piece.data[r] = make([]int, PIECESIZE)
			for c := 0; c < PIECESIZE; c++ {
				if pieceLine[c] == '#' {
					piece.data[r][c] = 1
				} else {
					piece.data[r][c] = 0
				}
			}
		}
		pieces = append(pieces, piece)
		i += PIECESIZE + 2
	}

	return pieces
}

func readFields(input []string) []Field {
	fields := []Field{}

	for i := 0; i < len(input); i++ {
		line := input[i]

		// fields format: "WxH: id id id id id id"

		tokens := strings.Split(line, ":")

		var width, height int
		n, err := fmt.Sscanf(tokens[0], "%dx%d", &width, &height)
		if err != nil || n != 2 {
			continue
		}

		pieceIDs := strings.Fields(tokens[1])
		piecesUsed := []int{}
		for _, pid := range pieceIDs {
			id, _ := strconv.Atoi(pid)
			piecesUsed = append(piecesUsed, id)
		}
		field := Field{
			width:      width,
			height:     height,
			piecesUsed: piecesUsed,
		}
		fields = append(fields, field)
	}

	return fields
}

func getPieceArea(piece Piece) int {
	area := 0
	for r := 0; r < PIECESIZE; r++ {
		for c := 0; c < PIECESIZE; c++ {
			area += piece.data[r][c]
		}
	}
	return area
}

func getFieldArea(field Field) int {
	return field.width * field.height
}

func canFit(field Field, pieces []Piece) bool {
	totalPieceArea := 0
	for pid, piece := range pieces {
		pieceCount := field.piecesUsed[pid]
		totalPieceArea += getPieceArea(piece) * pieceCount
	}

	fieldArea := getFieldArea(field)

	return totalPieceArea <= fieldArea
}

func run(input []string) int {
	pieces := readPieces(input)
	fields := readFields(input)

	fieldsThatFit := 0
	for _, field := range fields {
		if canFit(field, pieces) {
			fieldsThatFit++
		}
	}

	return fieldsThatFit
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

	/* success = test([]string{
		"0:", "###", "##.", "##.", "",
		"1:", "###", "##.", ".##", "",
		"2:", ".##", "###", "##.", "",
		"3:", "##.", "###", "##.", "",
		"4:", "###", "#..", "###", "",
		"5:", "###", ".#.", "###", "",
		"4x4: 0 0 0 0 2 0",
	}, 1) && success
	success = test([]string{
		"0:", "###", "##.", "##.", "",
		"1:", "###", "##.", ".##", "",
		"2:", ".##", "###", "##.", "",
		"3:", "##.", "###", "##.", "",
		"4:", "###", "#..", "###", "",
		"5:", "###", ".#.", "###", "",
		"12x5: 1 0 1 0 2 2",
	}, 1) && success
	success = test([]string{
		"0:", "###", "##.", "##.", "",
		"1:", "###", "##.", ".##", "",
		"2:", ".##", "###", "##.", "",
		"3:", "##.", "###", "##.", "",
		"4:", "###", "#..", "###", "",
		"5:", "###", ".#.", "###", "",
		"12x5: 1 0 1 0 3 2",
	}, 0) && success
	success = test([]string{
		"0:", "###", "##.", "##.", "",
		"1:", "###", "##.", ".##", "",
		"2:", ".##", "###", "##.", "",
		"3:", "##.", "###", "##.", "",
		"4:", "###", "#..", "###", "",
		"5:", "###", ".#.", "###", "",
		"4x4: 0 0 0 0 2 0",
		"12x5: 1 0 1 0 2 2",
		"12x5: 1 0 1 0 3 2",
	}, 2) && success */

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
