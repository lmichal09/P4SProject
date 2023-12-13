package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type HaversineTest struct {
	position1 OrderedPair
	position2 OrderedPair
	result    float64
}

type CheckDeadTest struct {
	fly    []Fly
	result bool
}

type DivideCountryTest struct {
	country Country
	result  []Quadrant
}

type GetQuadrantTest struct {
	fly       Fly
	quadrants []Quadrant
	result    int
}

type GetTemperatureTest struct {
	quadrantID int
	quadrants  []Quadrant
	result     float64
}

type InBoundsTest struct {
	fly       Fly
	quadrants []Quadrant
	result    bool
}

type CopyCountryTest struct {
	original Country
	result   Country
}

type CopyFlyTest struct {
	fly    Fly
	result Fly
}

type CopyTreeTest struct {
	tree   Tree
	result Tree
}

type CopyOrderedPairTest struct {
	original OrderedPair
	result   OrderedPair
}
type FindHostDirectionTest struct {
	fly         Fly
	nearestTree OrderedPair
	result      float64
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

func TestInBounds(t *testing.T) {
	tests := ReadInBoundsTests("Tests/InBounds/")

	for _, test := range tests {
		// run the test
		result := InBounds(&test.fly, test.quadrants)

		// check if the result is correct
		if result != test.result {
			t.Errorf("InBounds(%v, %v) = %v, want %v", test.fly, test.quadrants, result, test.result)
		}
	}

	for _, test := range tests {
		fmt.Printf("Running test with fly: %v, quadrants: %v\n", test.fly, test.quadrants)

		// run the test
		result := InBounds(&test.fly, test.quadrants)

		fmt.Printf("Result: %v, Expected: %v\n", result, test.result)

		// check if the result is correct
		if result != test.result {
			t.Errorf("InBounds(%v, %v) = %v, want %v", test.fly, test.quadrants, result, test.result)
		}
	}

}

func ReadInBoundsTests(directory string) []InBoundsTest {
	// Read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")

	tests := make([]InBoundsTest, 0)
	for _, inputFile := range inputFiles {
		// Read the input file
		tests = append(tests, ReadInBoundsTest(directory+"input/"+inputFile.Name()))
	}

	return tests
}

func ReadInBoundsTest(filename string) InBoundsTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var fly Fly
	var quadrants []Quadrant
	var result bool

	scanner := bufio.NewScanner(file)

	// Read each line for Fly data
	for scanner.Scan() {
		line := scanner.Text()

		// Assuming the line is not empty
		if line != "" {
			if strings.HasPrefix(line, "Fly Position:") {
				fly = ReadFlyFromPositionLine(line)
			} else if strings.HasPrefix(line, "Quadrants:") {
				quadrants = ReadQuadrants(line)
			} else if strings.HasPrefix(line, "Result:") {
				// Parse the result line
				var err error
				result, err = strconv.ParseBool(strings.TrimPrefix(line, "Result: "))
				if err != nil {
					panic(fmt.Sprintf("Error parsing result: %v", err))
				}
			}
		}
	}

	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return InBoundsTest{fly: fly, quadrants: quadrants, result: result}
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

func TestHaversine(t *testing.T) {
	tests := ReadHaversineTests("Tests/Haversine/")

	for _, test := range tests {
		// run the test
		result := Haversine(test.position1, test.position2)

		// check if the result is correct
		if result != test.result {
			t.Errorf("Haversine(%v, %v) = %v, want %v", test.position1, test.position2, result, test.result)
		}
	}
}

func ReadHaversineTests(directory string) []HaversineTest {
	inputFiles := ReadDirectory(directory + "input")
	tests := make([]HaversineTest, 0)

	for _, inputFile := range inputFiles {
		tests = append(tests, ReadHaversineTest(directory+"input/"+inputFile.Name()))
	}

	outputFiles := ReadDirectory(directory + "output")
	if len(outputFiles) != len(tests) {
		panic("Number of input and output files do not match")
	}

	for i, outputFile := range outputFiles {
		tests[i].result = ReadFloatFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadHaversineTest(filename string) HaversineTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var test HaversineTest
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		test.position1 = ParseOrderedPair(parts[0])
		test.position2 = ParseOrderedPair(parts[1])
		break
	}

	if scanner.Scan() {
		resultLine := scanner.Text()
		var err error
		test.result, err = strconv.ParseFloat(resultLine, 64)
		if err != nil {
			panic(fmt.Sprintf("Error parsing result: %v", err))
		}
	}

	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	return test
}
func TestDivideCountry(t *testing.T) {
	tests := ReadDivideCountryTests("Tests/DivideCountry/")

	for _, test := range tests {
		// run the test
		result := DivideCountry(test.country)

		// check if the result is correct
		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("DivideCountry(%v) = %v, want %v", test.country, result, test.result)
		}
	}
}
func TestGetQuadrant(t *testing.T) {
	tests := ReadGetQuadrantTests("Tests/GetQuadrant/")

	for _, test := range tests {
		// run the test
		result := GetQuadrant(&test.fly, test.quadrants)

		// check if the result is correct
		if result != test.result {
			t.Errorf("GetQuadrant(%v, %v) = %v, want %v", test.fly, test.quadrants, result, test.result)
		}
	}
}

func ReadGetQuadrantTests(directory string) []GetQuadrantTest {
	inputFiles := ReadDirectory(directory + "input")
	tests := make([]GetQuadrantTest, 0)

	for _, inputFile := range inputFiles {
		tests = append(tests, ReadGetQuadrantTest(directory+"input/"+inputFile.Name()))
	}

	outputFiles := ReadDirectory(directory + "output")
	if len(outputFiles) != len(tests) {
		panic("Number of input and output files do not match")
	}

	for i, outputFile := range outputFiles {
		tests[i].result = ReadIntFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}
func TestGetTemperature(t *testing.T) {
	tests := ReadGetTemperatureTests("Tests/GetTemperature/")

	for _, test := range tests {
		// run the test
		result := GetTemperature(test.quadrantID, test.quadrants)

		// check if the result is correct
		if result != test.result {
			t.Errorf("GetTemperature(%v, %v) = %v, want %v", test.quadrantID, test.quadrants, result, test.result)
		}
	}
}

func ReadGetTemperatureTests(directory string) []GetTemperatureTest {
	inputFiles := ReadDirectory(directory + "input")
	tests := make([]GetTemperatureTest, 0)

	for _, inputFile := range inputFiles {
		tests = append(tests, ReadGetTemperatureTest(directory+"input/"+inputFile.Name()))
	}

	outputFiles := ReadDirectory(directory + "output")
	if len(outputFiles) != len(tests) {
		panic("Number of input and output files do not match")
	}

	for i, outputFile := range outputFiles {
		tests[i].result = ReadFloatFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadGetTemperatureTest(filename string) GetTemperatureTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var test GetTemperatureTest
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		test.quadrantID, _ = strconv.Atoi(parts[0])
		test.quadrants = ReadQuadrants(strings.Join(parts[1:], " "))
		break
	}

	return test
}

func ReadGetQuadrantTest(filename string) GetQuadrantTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var fly Fly
	var quadrants []Quadrant
	var result int
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Fly Position:") {
			fly.position = ParseOrderedPair(strings.TrimPrefix(line, "Fly Position: "))
		} else if strings.HasPrefix(line, "Quadrants:") {
			quadrants = ReadQuadrants(strings.TrimPrefix(line, "Quadrants: "))
		}
	}

	return GetQuadrantTest{fly: fly, quadrants: quadrants, result: result}
}

func ReadDivideCountryTests(directory string) []DivideCountryTest {
	inputFiles := ReadDirectory(directory + "input")
	tests := make([]DivideCountryTest, 0)

	for _, inputFile := range inputFiles {
		tests = append(tests, ReadDivideCountryTest(directory+"input/"+inputFile.Name()))
	}

	outputFiles := ReadDirectory(directory + "output")
	if len(outputFiles) != len(tests) {
		panic("Number of input and output files do not match")
	}

	for i, outputFile := range outputFiles {
		tests[i].result = ReadQuadrantsFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadDivideCountryTest(filename string) DivideCountryTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var country Country
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		country = Country{width: ParseFloat(line)}
		break
	}

	return DivideCountryTest{country: country}
}

func TestFindHostDirection(t *testing.T) {
	tests := ReadFindHostDirectionTests("Tests/FindNearestHostDirection/")

	for _, test := range tests {
		// run the test
		result := FindHostDirection(test.fly.position, test.nearestTree)

		// check if the result is correct
		if result != test.result {
			t.Errorf("FindNearestHostDirection(%v, %v) = %v, want %v", test.fly, test.nearestTree, result, test.result)
		}
	}
}

// ConvertFlySliceToOrderedPairSlice converts a slice of Fly to a slice of OrderedPair
func ConvertFlySliceToOrderedPairSlice(flies []Fly) []OrderedPair {
	positions := make([]OrderedPair, len(flies))
	for i, fly := range flies {
		positions[i] = fly.position
	}
	return positions
}

func ReadFindHostDirectionTests(directory string) []FindHostDirectionTest {
	inputFiles := ReadDirectory(directory + "input")
	tests := make([]FindHostDirectionTest, 0)

	for _, inputFile := range inputFiles {
		tests = append(tests, ReadFindHostDirectionTest(directory+"input/"+inputFile.Name()))
	}

	outputFiles := ReadDirectory(directory + "output")
	if len(outputFiles) != len(tests) {
		panic("Number of input and output files do not match")
	}

	for i, outputFile := range outputFiles {
		tests[i].result = ReadFloatFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadFindHostDirectionTest(filename string) FindHostDirectionTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var fly Fly
	var nearestTree OrderedPair

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Fly Position:") {
			fly = ReadFlyFromPositionLine(line)
		} else if strings.HasPrefix(line, "Nearest Tree:") {
			nearestTree = ParseOrderedPair(strings.TrimPrefix(line, "Nearest Tree:"))
		}
	}

	return FindHostDirectionTest{fly: fly, nearestTree: nearestTree}
}

func TestCopyCountry(t *testing.T) {
	tests := ReadCopyCountryTests("Tests/CopyCountry/")

	for _, test := range tests {
		// run the test
		result := CopyCountry(test.original)

		// check if the result is correct
		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("CopyCountry(%v) = %v, want %v", test.original, result, test.result)
		}
	}
}

func ReadCopyCountryTests(directory string) []CopyCountryTest {
	inputFiles := ReadDirectory(directory + "input")
	tests := make([]CopyCountryTest, 0)

	for _, inputFile := range inputFiles {
		tests = append(tests, ReadCopyCountryTest(directory+"input/"+inputFile.Name()))
	}

	outputFiles := ReadDirectory(directory + "output")
	if len(outputFiles) != len(tests) {
		panic("Number of input and output files do not match")
	}

	for i, outputFile := range outputFiles {
		tests[i].result = ReadCountryFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadCopyCountryTest(filename string) CopyCountryTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var original Country

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		original = Country{width: ParseFloat(line)}
		break
	}

	return CopyCountryTest{original: original}
}

func ReadCountryFromFile(filename string) Country {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var country Country
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		country = Country{width: ParseFloat(line)}
		break
	}

	return country
}

func TestCopyTree(t *testing.T) {
	tests := ReadCopyTreeTests("Tests/CopyTree/")

	for _, test := range tests {
		// run the test
		result := CopyTree(test.tree)

		// check if the result is correct
		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("CopyTree(%v) = %v, want %v", test.tree, result, test.result)
		}
	}
}

func ReadCopyTreeTests(directory string) []CopyTreeTest {
	inputFiles := ReadDirectory(directory + "input")
	tests := make([]CopyTreeTest, 0)

	for _, inputFile := range inputFiles {
		tests = append(tests, ReadCopyTreeTest(directory+"input/"+inputFile.Name()))
	}

	outputFiles := ReadDirectory(directory + "output")
	if len(outputFiles) != len(tests) {
		panic("Number of input and output files do not match")
	}

	for i, outputFile := range outputFiles {
		tests[i].result = ReadTreeFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadCopyTreeTest(filename string) CopyTreeTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var tree Tree
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		tree = Tree{position: ParseOrderedPair(line)}
		break
	}

	return CopyTreeTest{tree: tree}
}

func TestCopyOrderedPair(t *testing.T) {
	tests := ReadCopyOrderedPairTests("Tests/CopyOrderedPair/")

	for _, test := range tests {
		// run the test
		result := CopyOrderedPair(test.original)

		// check if the result is correct
		if !reflect.DeepEqual(result, test.result) {
			t.Errorf("CopyOrderedPair(%v) = %v, want %v", test.original, result, test.result)
		}
	}
}

func ReadCopyOrderedPairTests(directory string) []CopyOrderedPairTest {
	inputFiles := ReadDirectory(directory + "input")
	tests := make([]CopyOrderedPairTest, 0)

	for _, inputFile := range inputFiles {
		tests = append(tests, ReadCopyOrderedPairTest(directory+"input/"+inputFile.Name()))
	}

	outputFiles := ReadDirectory(directory + "output")
	if len(outputFiles) != len(tests) {
		panic("Number of input and output files do not match")
	}

	for i, outputFile := range outputFiles {
		tests[i].result = ReadOrderedPairFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadCopyOrderedPairTest(filename string) CopyOrderedPairTest {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var original OrderedPair

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		original = ParseOrderedPair(line)
		break
	}

	return CopyOrderedPairTest{original: original}
}

func ReadOrderedPairFromFile(filename string) OrderedPair {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var pair OrderedPair
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		pair = ParseOrderedPair(line)
		break
	}

	return pair
}

func ReadTreeFromFile(filename string) Tree {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var tree Tree
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		tree = Tree{position: ParseOrderedPair(line)}
		break
	}

	return tree
}

func ReadQuadrantsFromFile(filename string) []Quadrant {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var quadrants []Quadrant
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		quadrants = append(quadrants, ParseQuadrant(line))
	}

	return quadrants
}

func ParseQuadrant(s string) Quadrant {
	parts := strings.Split(s, " ")

	if len(parts) != 5 {
		panic("Invalid number of parts in line")
	}

	x, _ := strconv.ParseFloat(parts[0], 64)
	y, _ := strconv.ParseFloat(parts[1], 64)
	width, _ := strconv.ParseFloat(parts[2], 64)
	id, _ := strconv.Atoi(parts[3])
	temp, _ := strconv.ParseFloat(parts[4], 64)

	return Quadrant{x: x, y: y, width: width, id: id, temp: temp}
}

func ReadFloatFromFile(filename string) float64 {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	f, err := strconv.ParseFloat(scanner.Text(), 64)
	if err != nil {
		panic(err)
	}

	return f
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

	if len(parts) != 6 {
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

	locationID, err := strconv.Atoi(parts[5])
	if err != nil {
		fmt.Printf("Error parsing locationID: %v\n", parts[5])
		panic(err)
	}
	fly.locationID = locationID

	color := ParseColor(parts[6])
	fly.color = color

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
func ParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
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

// ReadQuadrants parses the line containing quadrant information and returns a slice of Quadrant.
func ReadQuadrants(line string) []Quadrant {
	quadrants := make([]Quadrant, 0)

	// Extract the part of the line after "Quadrants:"
	quadrantStr := strings.TrimPrefix(line, "Quadrants: ")

	// Split the string into individual quadrants using comma as the separator
	quadrantTokens := strings.Split(quadrantStr, ", ")

	for _, token := range quadrantTokens {
		// Parse each token to create a Quadrant
		quadrant := ParseQuadrant(token)
		quadrants = append(quadrants, quadrant)
	}

	return quadrants
}

// ReadFlyFromPositionLine parses the line containing fly position information and returns a Fly object.
func ReadFlyFromPositionLine(line string) Fly {
	// Extract the part of the line after "Fly Position:"
	positionStr := strings.TrimPrefix(line, "Fly Position: ")

	// Split the string into individual values using space as the separator
	positionValues := strings.Split(positionStr, " ")

	// Parse each value and assign it to the corresponding variable
	x, _ := strconv.ParseFloat(positionValues[0], 64)
	y, _ := strconv.ParseFloat(positionValues[1], 64)

	// Create and return a Fly object with the parsed position
	return Fly{position: OrderedPair{x: x, y: y}}
}

func ReadIntFromFile(filename string) int {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read the integer from the file
	if scanner.Scan() {
		value, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic(err)
		}
		return value
	}

	panic("Unable to read integer from file")
}
func ReadHosts(line string) []Fly {
	// Parse and return the hosts from the line
	// Assuming the line format is like "Hosts: {host1_data} {host2_data} ..."
	parts := strings.Split(strings.TrimPrefix(line, "Hosts: "), " ")

	hosts := make([]Fly, len(parts))
	for i, part := range parts {
		hosts[i] = ReadFly(part)
	}

	return hosts
}
