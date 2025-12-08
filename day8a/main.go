package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

type Point3D struct {
	x int
	y int
	z int
}

func getDist(p1 Point3D, p2 Point3D) float64 {
	dx := p1.x - p2.x
	dy := p1.y - p2.y
	dz := p1.z - p2.z
	return math.Sqrt(float64(dx*dx + dy*dy + dz*dz))
}

func getPoints(input []string) []Point3D {
	points := []Point3D{}

	for _, line := range input {
		var x, y, z int
		fmt.Sscanf(line, "%d,%d,%d", &x, &y, &z)
		points = append(points, Point3D{x: x, y: y, z: z})
	}

	return points
}

func getDistances(points []Point3D) [][]float64 {
	dists := make([][]float64, len(points))
	for i := range dists {
		dists[i] = make([]float64, len(points))
	}

	for i := 0; i < len(points); i++ {
		for j := 0; j < len(points); j++ {
			if i != j {
				dists[i][j] = getDist(points[i], points[j])
				dists[j][i] = dists[i][j]
			} else {
				dists[i][j] = math.MaxFloat64
			}
		}
	}

	return dists
}

func getMinDist(dists [][]float64) (int, int, [][]float64) {
	minDist := math.MaxFloat64
	minI := -1
	minJ := -1

	for i := 0; i < len(dists); i++ {
		for j := 0; j < len(dists); j++ {
			if dists[i][j] < minDist {
				minDist = dists[i][j]
				minI = i
				minJ = j
			}
		}
	}

	dists[minI][minJ] = math.MaxFloat64
	dists[minJ][minI] = math.MaxFloat64

	return minI, minJ, dists
}

func isInCurcuit(curcuit []int, point int) bool {
	for _, p := range curcuit {
		if p == point {
			return true
		}
	}
	return false
}

func areConnected(curcuits [][]int, i int, j int) bool {
	for _, curcuit := range curcuits {
		if isInCurcuit(curcuit, i) && isInCurcuit(curcuit, j) {
			return true
		}
	}
	return false
}

func connectPoints(curcuits [][]int, i int, j int) [][]int {
	curcuitI := -1
	curcuitJ := -1

	for idx, curcuit := range curcuits {
		if isInCurcuit(curcuit, i) {
			curcuitI = idx
		}
		if isInCurcuit(curcuit, j) {
			curcuitJ = idx
		}
	}

	if curcuitI == -1 && curcuitJ == -1 {
		curcuits = append(curcuits, []int{i, j})
	} else if curcuitI != -1 && curcuitJ == -1 {
		curcuits[curcuitI] = append(curcuits[curcuitI], j)
	} else if curcuitI == -1 && curcuitJ != -1 {
		curcuits[curcuitJ] = append(curcuits[curcuitJ], i)
	} else if curcuitI != curcuitJ {
		curcuits[curcuitI] = append(curcuits[curcuitI], curcuits[curcuitJ]...)
		curcuits = append(curcuits[:curcuitJ], curcuits[curcuitJ+1:]...)
	}

	return curcuits
}

func multiplyBiggestCircuits(curcuits [][]int, numToMul int) int {
	curcuitsLengths := []int{}
	for _, curcuit := range curcuits {
		curcuitsLengths = append(curcuitsLengths, len(curcuit))
	}

	sort.Ints(curcuitsLengths)
	total := 1
	for i := 0; i < numToMul; i++ {
		total *= curcuitsLengths[len(curcuitsLengths)-1-i]
	}
	return total
}

func run(input []string, connections int, numToMul int) int {
	points := getPoints(input)
	dists := getDistances(points)
	curcuits := [][]int{}

	conn := 0
	for conn < connections {
		minI, minJ := -1, -1
		minI, minJ, dists = getMinDist(dists)
		if areConnected(curcuits, minI, minJ) {
			conn++
			continue
		}
		curcuits = connectPoints(curcuits, minI, minJ)
		conn++
	}

	res := multiplyBiggestCircuits(curcuits, numToMul)

	return res
}

func test(input []string, connections int, numToMul int, exp_output int) bool {

	output := run(input, connections, numToMul)

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
		"162,817,812",
		"57,618,57",
		"906,360,560",
		"592,479,940",
		"352,342,300",
		"466,668,158",
		"542,29,236",
		"431,825,988",
		"739,650,466",
		"52,470,668",
		"216,146,977",
		"819,987,18",
		"117,168,530",
		"805,96,715",
		"346,949,466",
		"970,615,88",
		"941,993,340",
		"862,61,35",
		"984,92,344",
		"425,690,689",
	}, 10, 3, 40) && success

	if success {
		fmt.Printf("-=-=-=-=-=-=-=-=-=-=-=-=-\n✅All tests passed!\n")
	} else {
		fmt.Printf("-=-=-=-=-=-=-=-=-=-=-=-=-\n❌Some tests failed!\n")
		return
	}

	input := readInput("input.txt")

	result := run(input, 1000, 3)
	fmt.Printf("Result: %d\n", result)
}
