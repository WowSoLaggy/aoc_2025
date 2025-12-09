package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sort"
	"strings"
)

type Point struct {
	x int
	y int
}

type Line struct {
	start Point
	end   Point
}

type Poly = []Line

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func minMax(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}

func readPoints(input []string) []Point {
	points := []Point{}

	for _, line := range input {
		var x, y int
		fmt.Sscanf(line, "%d,%d", &x, &y)
		points = append(points, Point{x: x, y: y})
	}

	return points
}

func createPoly(points []Point) Poly {
	poly := Poly{}
	for i := 0; i < len(points); i++ {
		start := points[i]
		end := points[(i+1)%len(points)]
		poly = append(poly, Line{start: start, end: end})
	}
	return poly
}

func getRectArea(p1, p2 Point) int {
	width := abs(p2.x-p1.x) + 1
	height := abs(p2.y-p1.y) + 1
	return width * height
}

func getUniqueCoords(points []Point) ([]int, []int) {
	xMap := make(map[int]bool)
	yMap := make(map[int]bool)

	for _, p := range points {
		xMap[p.x] = true
		yMap[p.y] = true
	}

	xCoords := []int{}
	yCoords := []int{}

	for x := range xMap {
		xCoords = append(xCoords, x)
	}
	for y := range yMap {
		yCoords = append(yCoords, y)
	}

	sort.Ints(xCoords)
	sort.Ints(yCoords)

	return xCoords, yCoords
}

func saveGridToPng(grid [][]int, filename string) {
	height := len(grid)
	width := len(grid[0])

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if grid[y][x] == 1 {
				img.Set(x, y, color.RGBA{0, 255, 0, 255}) // Green
			} else if grid[y][x] == 2 {
				img.Set(x, y, color.RGBA{255, 0, 0, 255}) // Red
			} else {
				img.Set(x, y, color.RGBA{0, 0, 0, 255}) // Black
			}
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}

func isPointInPoly(x, y float64, poly Poly) bool {
	inside := false
	j := len(poly) - 1
	for i := 0; i < len(poly); i++ {
		xi := float64(poly[i].start.x)
		yi := float64(poly[i].start.y)
		xj := float64(poly[j].start.x)
		yj := float64(poly[j].start.y)

		intersect := ((yi > y) != (yj > y)) && (x < (xj-xi)*(y-yi)/(yj-yi)+xi)
		if intersect {
			inside = !inside
		}
		j = i
	}
	return inside
}

func buildGrid(xCoords, yCoords []int, poly Poly) [][]int {
	grid := make([][]int, len(yCoords))
	for i := range grid {
		grid[i] = make([]int, len(xCoords))
	}

	for gridY1 := 0; gridY1 < len(yCoords)-1; gridY1++ {
		for gridX1 := 0; gridX1 < len(xCoords)-1; gridX1++ {
			gridY2 := gridY1 + 1
			gridX2 := gridX1 + 1

			cellCenterX := float64(xCoords[gridX1]+xCoords[gridX2]) / 2.0
			cellCenterY := float64(yCoords[gridY1]+yCoords[gridY2]) / 2.0
			if isPointInPoly(cellCenterX, cellCenterY, poly) {
				grid[gridY1][gridX1] = 1
			}
		}
	}

	return grid
}

func isRectInGrid(p1, p2 Point, grid [][]int, xCoords, yCoords []int) bool {
	x1Index, x2Index := -1, -1
	y1Index, y2Index := -1, -1

	minX, maxX := minMax(p1.x, p2.x)
	minY, maxY := minMax(p1.y, p2.y)

	for i, x := range xCoords {
		if x == minX {
			x1Index = i
		}
		if x == maxX {
			x2Index = i
		}
	}
	for i, y := range yCoords {
		if y == minY {
			y1Index = i
		}
		if y == maxY {
			y2Index = i
		}
	}

	for y := y1Index; y < y2Index; y++ {
		for x := x1Index; x < x2Index; x++ {
			if grid[y][x] != 1 {
				return false
			}
		}
	}
	return true
}

func drawRectInGrid(p1, p2 Point, grid [][]int, xCoords, yCoords []int) {
	x1Index, x2Index := -1, -1
	y1Index, y2Index := -1, -1

	minX, maxX := minMax(p1.x, p2.x)
	minY, maxY := minMax(p1.y, p2.y)

	for i, x := range xCoords {
		if x == minX {
			x1Index = i
		}
		if x == maxX {
			x2Index = i
		}
	}
	for i, y := range yCoords {
		if y == minY {
			y1Index = i
		}
		if y == maxY {
			y2Index = i
		}
	}

	for y := y1Index; y < y2Index; y++ {
		for x := x1Index; x < x2Index; x++ {
			grid[y][x] = 2
		}
	}
}

func findMaxArea(points []Point, grid [][]int, xCoords, yCoords []int) int {
	maxArea := 0
	maxp1 := Point{}
	maxp2 := Point{}
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			area := getRectArea(points[i], points[j])

			if area > maxArea {
				if isRectInGrid(points[i], points[j], grid, xCoords, yCoords) {
					maxArea = area
					maxp1 = points[i]
					maxp2 = points[j]
				}
			}
		}
	}

	drawRectInGrid(maxp1, maxp2, grid, xCoords, yCoords)
	return maxArea
}

func run(input []string) int {
	points := readPoints(input)
	poly := createPoly(points)
	xCoords, yCoords := getUniqueCoords(points)
	grid := buildGrid(xCoords, yCoords, poly)
	saveGridToPng(grid, "grid1.png")

	maxArea := findMaxArea(points, grid, xCoords, yCoords)
	saveGridToPng(grid, "grid2.png")

	return maxArea
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

	success = test([]string{
		"7,1",
		"11,1",
		"11,7",
		"9,7",
		"9,5",
		"2,5",
		"2,3",
		"7,3",
	}, 24) && success

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
