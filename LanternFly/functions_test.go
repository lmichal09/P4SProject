package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"testing"
)

type CheckDeadTest struct {
	fly    []Fly
	result bool
}

type CopyFlyTest struct {
	fly    Fly
	result Fly
}

func TestCheckDead(t *testing.T) {
	tests := ReadCheckDeadTests("Tests/CheckDead/")

	for _, test := range tests {
		// run the test
		result := CheckDead(test.fly)

		// check if the result if correct
		if result != test.result {
			t.Errorf("CheckDead(%v) = %v, want %v", test.fly, result, test.result)
		}
	}
}

func ReadCheckDeadTests(directory string) []CheckDeadTest {
	// Read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")

	tests := make([]CheckDeadTest, 0)
	for _, inputFile := range inputFiles {
		// Read the input file
		tests = append(tests, ReadCheckDeadTest(directory+"input/"+inputFile.Name()))
	}

	outputFiles := ReadDirectory(directory + "output")

	if len(outputFiles) != len(tests) {
		panic("Number of input and output files do not match")
	}

	for i, outputFile := range outputFiles {
		tests[i].result = ReadBoolFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadCheckDeadTest(filename string) CheckDeadTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var flies []Fly
	scanner := bufio.NewScanner(file)

	// Read each line for Fly data
	for scanner.Scan() {
		line := scanner.Text()
		// Assuming the line is not empty and not the result line
		if line != "" && !strings.HasPrefix(line, "Result:") {
			fly := ReadFly(line)
			flies = append(flies, fly)
		} else {
			break // Break out of the loop when the result line or empty line is encountered
		}
	}

	// The next line is assumed to be the result
	var result bool
	if scanner.Scan() {
		resultLine := scanner.Text()
		var err error
		result, err = strconv.ParseBool(resultLine)
		if err != nil {
			panic(fmt.Sprintf("Error parsing result: %v", err))
		}
	}

	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return CheckDeadTest{fly: flies, result: result}
}

func TestCopyFly(t *testing.T) {
	tests := ReadCopyFlyTests("Tests/CopyFly/")

	for _, test := range tests {
		// run the test
		result := CopyFly(test.fly)

		// check if the result if correct
		if result != test.result {
			t.Errorf("CopyFly(%v) = %v, want %v", test.fly, result, test.result)
		}
	}
}

func ReadCopyFlyTests(directory string) []CopyFlyTest {
	// Read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")

	tests := make([]CopyFlyTest, 0)
	for _, inputFile := range inputFiles {
		// Read the input file
		tests = append(tests, ReadCopyFlyTest(directory+"input/"+inputFile.Name()))
	}

	outputFiles := ReadDirectory(directory + "output")

	if len(outputFiles) != len(tests) {
		panic("Number of input and output files do not match")
	}

	for i, outputFile := range outputFiles {
		tests[i].result = ReadFly(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadCopyFlyTest(filename string) CopyFlyTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read each line for Fly data
	scanner.Scan()
	line := scanner.Text()
	fly := ReadFly(line)

	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return CopyFlyTest{fly: fly}
}

func ReadDirectory(dir string) []fs.DirEntry {
	//read in all files in the given directory
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	return files
}

func ReadFly(line string) Fly {
	fly := Fly{}

	parts := strings.Split(line, " ")

	if len(parts) != 5 {
		fmt.Println("parts", len(parts))
		panic("Invalid number of parts in line")
	}

	posX, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		panic(err)
	}

	posY, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		panic(err)
	}

	fly.position = OrderedPair{x: posX, y: posY}

	fmt.Println("position", fly.position)

	stage, err := strconv.Atoi(parts[2])
	if err != nil {
		panic(err)
	}

	fly.stage = stage

	energy, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		panic(err)
	}

	fly.energy = energy

	isAlive, err := strconv.ParseBool(parts[4])
	if err != nil {
		fmt.Printf("Error parsing isAlive: %v\n", parts[4])
		panic(err)
	}

	fly.isAlive = isAlive

	return fly
}

func ReadBoolFromFile(filename string) bool {
	//read in the bool from the given file
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	//read in the bool
	scanner.Scan()
	b, err := strconv.ParseBool(scanner.Text())
	if err != nil {
		panic(err)
	}

	return b
}

func ParseOrderedPair(s string) OrderedPair {
	// Parse the string into an OrderedPair
	parts := strings.Split(s, " ")
	if len(parts) != 2 {
		panic("Invalid number of parts in line")
	}

	x, _ := strconv.ParseFloat(parts[0], 64)
	y, _ := strconv.ParseFloat(parts[1], 64)

	return OrderedPair{x: x, y: y}
}

func ParseColor(s string) Color {
	// Parse the string into a Color
	parts := strings.Split(s, " ")
	if len(parts) != 3 {
		panic("Invalid number of parts in line")
	}

	red, _ := strconv.ParseUint(parts[0], 10, 8)
	green, _ := strconv.ParseUint(parts[1], 10, 8)
	blue, _ := strconv.ParseUint(parts[2], 10, 8)

	return Color{red: uint8(red), green: uint8(green), blue: uint8(blue)}
}
func TestComputeDegreeDay(t *testing.T) {
	tests := ReadComputeDegreeDayTests("Tests/ComputeDegreeDay/")

	for _, test := range tests {
		// run the test
		result := ComputeDegreeDay(test.fly)

		// check if the result if correct
		if result != test.result {
			t.Errorf("ComputeDegreeDay(%v) = %v, want %v", test.fly, result, test.result)
		}
	}
}

func ReadComputeDegreeDayTests(directory string) []ComputeDegreeDayTest {
	// Read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")

	tests := make([]ComputeDegreeDayTest, 0)
	for _, inputFile := range inputFiles {
		// Read the input file
		tests = append(tests, ReadComputeDegreeDayTest(directory+"input/"+inputFile.Name()))
	}

	outputFiles := ReadDirectory(directory + "output")

	if len(outputFiles) != len(tests) {
		panic("Number of input and output files do not match")
	}

	for i, outputFile := range outputFiles {
		tests[i].result = ReadBoolFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadComputeDegreeDayTest(filename string) ComputeDegreeDayTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var flies []Fly
	scanner := bufio.NewScanner(file)

	// Read each line for Fly data
	for scanner.Scan() {
		line := scanner.Text()
		// Assuming the line is not empty and not the result line
		if line != "" && !strings.HasPrefix(line, "Result:") {
			fly := ReadFly(line)
			flies = append(flies, fly)
		} else {
			break // Break out of the loop when the result line or empty line is encountered
		}
	}

	// The next line is assumed to be the result
	var result bool
	if scanner.Scan() {
		resultLine := scanner.Text()
		var err error
		result, err = strconv.ParseBool(resultLine)
		if err != nil {
			panic(fmt.Sprintf("Error parsing result: %v", err))
		}
	}

	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return ComputeDegreeDayTest{fly: flies, result: result}
}
func TestComputeFecundity(t *testing.T) {
	tests := ReadComputeFecundityTests("Tests/ComputeFecundity/")

	for _, test := range tests {
		// run the test
		result := ComputeFecundity(test.fly)

		// check if the result if correct
		if result != test.result {
			t.Errorf("ComputeFecundity(%v) = %v, want %v", test.fly, result, test.result)
		}
	}
}

func ReadComputeFecundityTests(directory string) []ComputeFecundityTest {
	// Read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")

	tests := make([]ComputeFecundityTest, 0)
	for _, inputFile := range inputFiles {
		// Read the input file
		tests = append(tests, ReadComputeFecundityTest(directory+"input/"+inputFile.Name()))
	}

	outputFiles := ReadDirectory(directory + "output")

	if len(outputFiles) != len(tests) {
		panic("Number of input and output files do not match")
	}

	for i, outputFile := range outputFiles {
		tests[i].result = ReadBoolFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadComputeFecundityTest(filename string) ComputeFecundityTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var flies []Fly
	scanner := bufio.NewScanner(file)

	// Read each line for Fly data
	for scanner.Scan() {
		line := scanner.Text()
		// Assuming the line is not empty and not the result line
		if line != "" && !strings.HasPrefix(line, "Result:") {
			fly := ReadFly(line)
			flies = append(flies, fly)
		} else {
			break // Break out of the loop when the result line or empty line is encountered
		}
	}

	// The next line is assumed to be the result
	var result bool
	if scanner.Scan() {
		resultLine := scanner.Text()
		var err error
		result, err = strconv.ParseBool(resultLine)
		if err != nil {
			panic(fmt.Sprintf("Error parsing result: %v", err))
		}
	}

	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return ComputeFecundityTest{fly: flies, result: result}
}
func TestComputeMorality(t *testing.T) {
	tests := ReadComputeMoralityTests("Tests/ComputeMorality/")

	for _, test := range tests {
		// run the test
		result := ComputeMorality(test.fly)

		// check if the result if correct
		if result != test.result {
			t.Errorf("ComputeMorality(%v) = %v, want %v", test.fly, result, test.result)
		}
	}
}

func ReadComputeMoralityTests(directory string) []ComputeMoralityTest {
	// Read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")

	tests := make([]ComputeMoralityTest, 0)
	for _, inputFile := range inputFiles {
		// Read the input file
		tests = append(tests, ReadComputeMoralityTest(directory+"input/"+inputFile.Name()))
	}

	outputFiles := ReadDirectory(directory + "output")

	if len(outputFiles) != len(tests) {
		panic("Number of input and output files do not match")
	}

	for i, outputFile := range outputFiles {
		tests[i].result = ReadBoolFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadComputeMoralityTest(filename string) ComputeMoralityTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var flies []Fly
	scanner := bufio.NewScanner(file)

	// Read each line for Fly data
	for scanner.Scan() {
		line := scanner.Text()
		// Assuming the line is not empty and not the result line
		if line != "" && !strings.HasPrefix(line, "Result:") {
			fly := ReadFly(line)
			flies = append(flies, fly)
		} else {
			break // Break out of the loop when the result line or empty line is encountered
		}
	}

	// The next line is assumed to be the result
	var result bool
	if scanner.Scan() {
		resultLine := scanner.Text()
		var err error
		result, err = strconv.ParseBool(resultLine)
		if err != nil {
			panic(fmt.Sprintf("Error parsing result: %v", err))
		}
	}

	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return ComputeMoralityTest{fly: flies, result: result}
}
func TestComputeMovement(t *testing.T) {
	tests := ReadComputeMovementTests("Tests/ComputeMovement/")

	for _, test := range tests {
		// run the test
		result := ComputeMovement(test.fly)

		// check if the result if correct
		if result != test.result {
			t.Errorf("ComputeMovement(%v) = %v, want %v", test.fly, result, test.result)
		}
	}
}

func ReadComputeMovementTests(directory string) []ComputeMovementTest {
	// Read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")

	tests := make([]ComputeMovementTest, 0)
	for _, inputFile := range inputFiles {
		// Read the input file
		tests = append(tests, ReadComputeMovementTest(directory+"input/"+inputFile.Name()))
	}

	outputFiles := ReadDirectory(directory + "output")

	if len(outputFiles) != len(tests) {
		panic("Number of input and output files do not match")
	}

	for i, outputFile := range outputFiles {
		tests[i].result = ReadBoolFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadComputeMovementTest(filename string) ComputeMovementTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var flies []Fly
	scanner := bufio.NewScanner(file)

	// Read each line for Fly data
	for scanner.Scan() {
		line := scanner.Text()
		// Assuming the line is not empty and not the result line
		if line != "" && !strings.HasPrefix(line, "Result:") {
			fly := ReadFly(line)
			flies = append(flies, fly)
		} else {
			break // Break out of the loop when the result line or empty line is encountered
		}
	}

	// The next line is assumed to be the result
	var result bool
	if scanner.Scan() {
		resultLine := scanner.Text()
		var err error
		result, err = strconv.ParseBool(resultLine)
		if err != nil {
			panic(fmt.Sprintf("Error parsing result: %v", err))
		}
	}

	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return ComputeMovementTest{fly: flies, result: result}
}
func TestDirectedMovement(t *testing.T) {
	tests := ReadDirectedMovementTests("Tests/DirectedMovement/")

	for _, test := range tests {
		// run the test
		result := DirectedMovement(test.fly)

		// check if the result if correct
		if result != test.result {
			t.Errorf("DirectedMovement(%v) = %v, want %v", test.fly, result, test.result)
		}
	}
}

func ReadDirectedMovementTests(directory string) []DirectedMovementTest {
	// Read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")

	tests := make([]DirectedMovementTest, 0)
	for _, inputFile := range inputFiles {
		// Read the input file
		tests = append(tests, ReadDirectedMovementTest(directory+"input/"+inputFile.Name()))
	}

	outputFiles := ReadDirectory(directory + "output")

	if len(outputFiles) != len(tests) {
		panic("Number of input and output files do not match")
	}

	for i, outputFile := range outputFiles {
		tests[i].result = ReadBoolFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadDirectedMovementTest(filename string) DirectedMovementTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var flies []Fly
	scanner := bufio.NewScanner(file)

	// Read each line for Fly data
	for scanner.Scan() {
		line := scanner.Text()
		// Assuming the line is not empty and not the result line
		if line != "" && !strings.HasPrefix(line, "Result:") {
			fly := ReadFly(line)
			flies = append(flies, fly)
		} else {
			break // Break out of the loop when the result line or empty line is encountered
		}
	}

	// The next line is assumed to be the result
	var result bool
	if scanner.Scan() {
		resultLine := scanner.Text()
		var err error
		result, err = strconv.ParseBool(resultLine)
		if err != nil {
			panic(fmt.Sprintf("Error parsing result: %v", err))
		}
	}

	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return DirectedMovementTest{fly: flies, result: result}
}
func TestFindNearestTree(t *testing.T) {
	tests := ReadFindNearestTreeTests("Tests/FindNearestTree/")

	for _, test := range tests {
		// run the test
		result := FindNearestTree(test.fly)

		// check if the result if correct
		if result != test.result {
			t.Errorf("FindNearestTree(%v) = %v, want %v", test.fly, result, test.result)
		}
	}
}

func ReadFindNearestTreeTests(directory string) []FindNearestTreeTest {
	// Read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")

	tests := make([]FindNearestTreeTest, 0)
	for _, inputFile := range inputFiles {
		// Read the input file
		tests = append(tests, ReadFindNearestTreeTest(directory+"input/"+inputFile.Name()))
	}

	outputFiles := ReadDirectory(directory + "output")

	if len(outputFiles) != len(tests) {
		panic("Number of input and output files do not match")
	}

	for i, outputFile := range outputFiles {
		tests[i].result = ReadBoolFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadFindNearestTreeTest(filename string) FindNearestTreeTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var flies []Fly
	scanner := bufio.NewScanner(file)

	// Read each line for Fly data
	for scanner.Scan() {
		line := scanner.Text()
		// Assuming the line is not empty and not the result line
		if line != "" && !strings.HasPrefix(line, "Result:") {
			fly := ReadFly(line)
			flies = append(flies, fly)
		} else {
			break // Break out of the loop when the result line or empty line is encountered
		}
	}

	// The next line is assumed to be the result
	var result bool
	if scanner.Scan() {
		resultLine := scanner.Text()
		var err error
		result, err = strconv.ParseBool(resultLine)
		if err != nil {
			panic(fmt.Sprintf("Error parsing result: %v", err))
		}
	}

	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return FindNearestTreeTest{fly: flies, result: result}
}
func TestGetBaseTemp(t *testing.T) {
	tests := ReadGetBaseTempTests("Tests/GetBaseTemp/")

	for _, test := range tests {
		// run the test
		result := GetBaseTemp(test.fly)

		// check if the result if correct
		if result != test.result {
			t.Errorf("GetBaseTemp(%v) = %v, want %v", test.fly, result, test.result)
		}
	}
}

func ReadGetBaseTempTests(directory string) []GetBaseTempTest {
	// Read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")

	tests := make([]GetBaseTempTest, 0)
	for _, inputFile := range inputFiles {
		// Read the input file
		tests = append(tests, ReadGetBaseTempTest(directory+"input/"+inputFile.Name()))
	}

	outputFiles := ReadDirectory(directory + "output")

	if len(outputFiles) != len(tests) {
		panic("Number of input and output files do not match")
	}

	for i, outputFile := range outputFiles {
		tests[i].result = ReadBoolFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadGetBaseTempTest(filename string) GetBaseTempTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var flies []Fly
	scanner := bufio.NewScanner(file)

	// Read each line for Fly data
	for scanner.Scan() {
		line := scanner.Text()
		// Assuming the line is not empty and not the result line
		if line != "" && !strings.HasPrefix(line, "Result:") {
			fly := ReadFly(line)
			flies = append(flies, fly)
		} else {
			break // Break out of the loop when the result line or empty line is encountered
		}
	}

	// The next line is assumed to be the result
	var result bool
	if scanner.Scan() {
		resultLine := scanner.Text()
		var err error
		result, err = strconv.ParseBool(resultLine)
		if err != nil {
			panic(fmt.Sprintf("Error parsing result: %v", err))
		}
	}

	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return GetBaseTempTest{fly: flies, result: result}
}
func TestGetTreePositions(t *testing.T) {
	tests := ReadGetTreePositionsTests("Tests/GetTreePositions/")

	for _, test := range tests {
		// run the test
		result := GetTreePositions(test.fly)

		// check if the result if correct
		if result != test.result {
			t.Errorf("GetTreePositions(%v) = %v, want %v", test.fly, result, test.result)
		}
	}
}

func ReadGetTreePositionsTests(directory string) []GetTreePositionsTest {
	// Read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")

	tests := make([]GetTreePositionsTest, 0)
	for _, inputFile := range inputFiles {
		// Read the input file
		tests = append(tests, ReadGetTreePositionsTest(directory+"input/"+inputFile.Name()))
	}

	outputFiles := ReadDirectory(directory + "output")

	if len(outputFiles) != len(tests) {
		panic("Number of input and output files do not match")
	}

	for i, outputFile := range outputFiles {
		tests[i].result = ReadBoolFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadGetTreePositionsTest(filename string) GetTreePositionsTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var flies []Fly
	scanner := bufio.NewScanner(file)

	// Read each line for Fly data
	for scanner.Scan() {
		line := scanner.Text()
		// Assuming the line is not empty and not the result line
		if line != "" && !strings.HasPrefix(line, "Result:") {
			fly := ReadFly(line)
			flies = append(flies, fly)
		} else {
			break // Break out of the loop when the result line or empty line is encountered
		}
	}

	// The next line is assumed to be the result
	var result bool
	if scanner.Scan() {
		resultLine := scanner.Text()
		var err error
		result, err = strconv.ParseBool(resultLine)
		if err != nil {
			panic(fmt.Sprintf("Error parsing result: %v", err))
		}
	}

	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return GetTreePositionsTest{fly: flies, result: result}
}
func TestRandomMovement(t *testing.T) {
	tests := ReadRandomMovementTests("Tests/RandomMovement/")

	for _, test := range tests {
		// run the test
		result := RandomMovement(test.fly)

		// check if the result if correct
		if result != test.result {
			t.Errorf("RandomMovement(%v) = %v, want %v", test.fly, result, test.result)
		}
	}
}

func ReadRandomMovementTests(directory string) []RandomMovementTest {
	// Read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")

	tests := make([]RandomMovementTest, 0)
	for _, inputFile := range inputFiles {
		// Read the input file
		tests = append(tests, ReadRandomMovementTest(directory+"input/"+inputFile.Name()))
	}

	outputFiles := ReadDirectory(directory + "output")

	if len(outputFiles) != len(tests) {
		panic("Number of input and output files do not match")
	}

	for i, outputFile := range outputFiles {
		tests[i].result = ReadBoolFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadRandomMovementTest(filename string) RandomMovementTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var flies []Fly
	scanner := bufio.NewScanner(file)

	// Read each line for Fly data
	for scanner.Scan() {
		line := scanner.Text()
		// Assuming the line is not empty and not the result line
		if line != "" && !strings.HasPrefix(line, "Result:") {
			fly := ReadFly(line)
			flies = append(flies, fly)
		} else {
			break // Break out of the loop when the result line or empty line is encountered
		}
	}

	// The next line is assumed to be the result
	var result bool
	if scanner.Scan() {
		resultLine := scanner.Text()
		var err error
		result, err = strconv.ParseBool(resultLine)
		if err != nil {
			panic(fmt.Sprintf("Error parsing result: %v", err))
		}
	}

	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return RandomMovementTest{fly: flies, result: result}
}
func TestUpdateLifeStage(t *testing.T) {
	tests := ReadUpdateLifeStageTests("Tests/UpdateLifeStage/")

	for _, test := range tests {
		// run the test
		result := UpdateLifeStage(test.fly)

		// check if the result if correct
		if result != test.result {
			t.Errorf("UpdateLifeStage(%v) = %v, want %v", test.fly, result, test.result)
		}
	}
}

func ReadUpdateLifeStageTests(directory string) []UpdateLifeStageTest {
	// Read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")

	tests := make([]UpdateLifeStageTest, 0)
	for _, inputFile := range inputFiles {
		// Read the input file
		tests = append(tests, ReadUpdateLifeStageTest(directory+"input/"+inputFile.Name()))
	}

	outputFiles := ReadDirectory(directory + "output")

	if len(outputFiles) != len(tests) {
		panic("Number of input and output files do not match")
	}

	for i, outputFile := range outputFiles {
		tests[i].result = ReadBoolFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadUpdateLifeStageTest(filename string) UpdateLifeStageTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var flies []Fly
	scanner := bufio.NewScanner(file)

	// Read each line for Fly data
	for scanner.Scan() {
		line := scanner.Text()
		// Assuming the line is not empty and not the result line
		if line != "" && !strings.HasPrefix(line, "Result:") {
			fly := ReadFly(line)
			flies = append(flies, fly)
		} else {
			break // Break out of the loop when the result line or empty line is encountered
		}
	}

	// The next line is assumed to be the result
	var result bool
	if scanner.Scan() {
		resultLine := scanner.Text()
		var err error
		result, err = strconv.ParseBool(resultLine)
		if err != nil {
			panic(fmt.Sprintf("Error parsing result: %v", err))
		}
	}

	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return UpdateLifeStageTest{fly: flies, result: result}
}
