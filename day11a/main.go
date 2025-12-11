package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

type Device struct {
	name string
	out  []string
}

type DeviceMap map[string]*Device

func readDevices(input []string) DeviceMap {
	devices := make(DeviceMap)

	for _, line := range input {
		parts := strings.Split(line, ": ")
		name := parts[0]
		outputs := strings.Split(parts[1], " ")

		device, exists := devices[name]
		if !exists {
			device = &Device{name: name}
			devices[name] = device
		}
		device.out = outputs

		for _, outName := range outputs {
			if _, exists := devices[outName]; !exists {
				devices[outName] = &Device{name: outName}
			}
		}
	}

	return devices
}

func dfs(device *Device, devices DeviceMap, visited map[string]bool) int {
	if device.name == "out" {
		return 1
	}
	visited[device.name] = true
	totalPaths := 0
	for _, outName := range device.out {
		if !visited[outName] {
			totalPaths += dfs(devices[outName], devices, visited)
		}
	}
	visited[device.name] = false
	return totalPaths
}

func findPathsToOut(devices DeviceMap) int {
	startDevice := devices["you"]
	visited := make(map[string]bool)
	return dfs(startDevice, devices, visited)
}

func run(input []string) int {
	devices := readDevices(input)
	paths := findPathsToOut(devices)
	return paths
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

	success = test([]string{"you: out ccc"}, 1) && success
	success = test([]string{"you: bbb ccc", "bbb: ddd out"}, 1) && success
	success = test([]string{"you: bbb ccc", "bbb: ddd out", "ccc: you out"}, 2) && success
	success = test([]string{
		"aaa: you hhh",
		"you: bbb ccc",
		"bbb: ddd eee",
		"ccc: ddd eee fff",
		"ddd: ggg",
		"eee: out",
		"fff: out",
		"ggg: out",
		"hhh: ccc fff iii",
		"iii: out",
	}, 5) && success

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
