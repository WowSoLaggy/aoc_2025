package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

type Device struct {
	id  int
	out []int
}
type DeviceMap map[int]*Device

type Dict map[string]int

func readDevices(input []string, names Dict) DeviceMap {
	devices := make(DeviceMap)

	// perpare names dict
	for _, line := range input {
		parts := strings.Split(line, ": ")
		name := parts[0]
		_, exists := names[name]
		if exists {
			panic("duplicate name")
		}
		names[name] = len(names)
	}

	// if no "out" in names, add it
	_, exists := names["out"]
	if !exists {
		names["out"] = len(names)
	}

	for _, line := range input {
		parts := strings.Split(line, ": ")
		name := parts[0]
		outputs := strings.Split(parts[1], " ")

		id, exists := names[name]
		if !exists {
			panic("name not found")
		}

		device, exists := devices[id]
		if !exists {
			device = &Device{id: id}
			devices[id] = device
		}

		for _, outName := range outputs {
			outID, exists := names[outName]
			if !exists {
				panic("name not found")
			}
			device.out = append(device.out, outID)

			if _, exists := devices[outID]; !exists {
				devices[outID] = &Device{id: outID}
			}
		}
	}

	return devices
}

func sortTopological(devices DeviceMap) []*Device {
	sorted := []*Device{}
	visited := make(map[int]bool)
	var visit func(device *Device)
	visit = func(device *Device) {
		if visited[device.id] {
			return
		}
		visited[device.id] = true
		for _, outId := range device.out {
			visit(devices[outId])
		}
		sorted = append(sorted, device)
	}
	for _, device := range devices {
		visit(device)
	}
	return sorted
}

func findPathsToOut(devices DeviceMap, names Dict, startName string, finishName string) int {
	startId := names[startName]
	finishId := names[finishName]

	sortedDevices := sortTopological(devices)
	pathsCount := make(map[int]int)
	pathsCount[startId] = 1
	for i := len(sortedDevices) - 1; i >= 0; i-- {
		device := sortedDevices[i]
		for _, outId := range device.out {
			pathsCount[outId] += pathsCount[device.id]
		}
	}
	totalPaths := pathsCount[finishId]

	return totalPaths
}

func run(input []string) int {
	names := make(Dict)
	devices := readDevices(input, names)

	//paths := findPathsToOut(devices, names, "dac", "out") // 3420

	pathsSvrFft := findPathsToOut(devices, names, "svr", "fft")
	fmt.Printf("pathsSvrFft: %d\n", pathsSvrFft)
	pathsFftDac := findPathsToOut(devices, names, "fft", "dac")
	fmt.Printf("pathsFftDac: %d\n", pathsFftDac)
	pathsDacOut := findPathsToOut(devices, names, "dac", "out")
	fmt.Printf("pathsDacOut: %d\n", pathsDacOut)
	pathsSvrDac := findPathsToOut(devices, names, "svr", "dac")
	fmt.Printf("pathsSvrDac: %d\n", pathsSvrDac)
	pathsDacFft := findPathsToOut(devices, names, "dac", "fft")
	fmt.Printf("pathsDacFft: %d\n", pathsDacFft)
	pathsFftOut := findPathsToOut(devices, names, "fft", "out")
	fmt.Printf("pathsFftOut: %d\n", pathsFftOut)

	// 4277326870542356560 too high

	return pathsSvrFft*pathsFftDac*pathsDacOut + pathsSvrDac*pathsDacFft*pathsFftOut
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

	success = test([]string{
		"svr: aaa bbb",
		"aaa: fft",
		"fft: ccc",
		"bbb: tty",
		"tty: ccc",
		"ccc: ddd eee",
		"ddd: hub",
		"hub: fff",
		"eee: dac",
		"dac: fff",
		"fff: ggg hhh",
		"ggg: out",
		"hhh: out",
	}, 2) && success

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
