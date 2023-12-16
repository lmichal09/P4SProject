package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// parseInteger takes a string input, converts it to an integer using the strconv.Atoi function, and returns the integer value.
// If the conversion fails, an error message is returned.
func parseInteger(s string) (int, error) {
	// Convert string to integer.
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("error parsing integer: %v", err)
	}
	return val, nil
}

// parseFloat is designed to convert a string to a float64 data type.
// The strconv.ParseFloat() function is used for this conversion.
// If an error occurs during the conversion, the error is returned along with a float64 value of 0.
// Otherwise, the converted float64 value and a nil error are returned.
func parseFloat(s string) (float64, error) {
	// Convert string to float64.
	val, err := strconv.ParseFloat(s, 64)

	if err != nil {
		return 0, fmt.Errorf("error parsing float: %v", err)
	}

	return val, nil
}

// parseBool takes a string s as input and returns a boolean value and an error.
// It uses the strconv.ParseBool function from the strconv package to parse the string.
// If the parsing fails, it returns a boolean value of false and an error.
// Otherwise, it returns the parsed boolean value and a nil error.
func parseBool(s string) (bool, error) {
	// ParseBool converts a string into a boolean value.
	val, err := strconv.ParseBool(s)
	if err != nil {
		return false, fmt.Errorf("error parsing bool: %v", err)
	}
	return val, nil
}

// parseString removes the leading and trailing spaces from the given string s and returns the resulting string.
func parseString(s string) string {
	// Remove leading and trailing spaces
	return strings.TrimSpace(s)
}

// ReadSampleDatafromFile takes in a filename and returns a SampleData datatype and error
func ReadSampleDataFromFile(filename string) ([]SampleData, error) {
	var sampleData []SampleData

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	isFirstRow := true
	for scanner.Scan() {
		if isFirstRow {
			isFirstRow = false
			continue
		}

		line := scanner.Text()
		fields := strings.Split(line, "\t")

		if len(fields) < 14 { // Ensure each line has sufficient data
			continue
		}

		if fields[2] != "2021" {
			continue
		}

		if fields[6] == "NA" || fields[6] == "" || fields[7] == "NA" || fields[7] == "" {
			continue // Skip the line if either field is "NA" or empty
		}
		year, err := parseInteger(fields[1])
		if err != nil {
			fmt.Printf("Error parsing year in row: %v\n", err)
			continue
		}

		bioYear, err := parseInteger(fields[2])
		if err != nil {
			fmt.Printf("Error parsing bioYear in row: %v\n", err)
			continue
		}

		latitude, err := parseFloat(fields[3])
		if err != nil {
			fmt.Printf("Error parsing latitude in row: %v\n", err)
			continue
		}

		longitude, err := parseFloat(fields[4])
		if err != nil {
			fmt.Printf("Error parsing longitude in row: %v\n", err)
			continue
		}

		lydePresent, err := parseBool(fields[6])
		if err != nil {
			fmt.Printf("Error parsing lydePresent in row: %v\n", err)
		}
		if lydePresent == false {
			continue
		}

		lydeEstablished, err := parseBool(fields[7])
		if err != nil {
			fmt.Printf("Error parsing lydeEstablished in row: %v\n", err)
			continue
		}

		roundedLongitude, err := parseFloat(fields[12])
		if err != nil {
			fmt.Printf("Error parsing roundedLongitude in row: %v\n", err)
			continue
		}

		roundedLatitude, err := parseFloat(fields[13])
		if err != nil {
			fmt.Printf("Error parsing roundedLatitude in row: %v\n", err)
			continue
		}

		data := SampleData{
			Source:           fields[0],
			Year:             year,
			BioYear:          bioYear,
			Latitude:         latitude,
			Longitude:        longitude,
			State:            fields[5],
			LydePresent:      lydePresent,
			LydeEstablished:  lydeEstablished,
			LydeDensity:      fields[8],
			SourceAgency:     fields[9],
			CollectionMethod: fields[10],
			PointID:          fields[11],
			RoundedLongitude: roundedLongitude,
			RoundedLatitude:  roundedLatitude,
		}

		sampleData = append(sampleData, data)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return sampleData, nil
}

// ProcessFile function reads a CSV file and extracts specific data from it.
// It returns the extracted data as an OrderedPair and an error if there are any.
func ProcessFile(filePath string) (OrderedPair, error) {
	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return OrderedPair{}, err
	}
	defer file.Close()

	// Read file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return OrderedPair{}, err
	}

	// Filter records
	var processed []string
	for i, record := range records {
		if i == 5 || i == 12 {
			if len(record) > 1 {
				processed = append(processed, record[1])
			}
		}
	}

	// Process results
	var result OrderedPair
	if len(processed) >= 2 {
		maxTemp, err := strconv.ParseFloat(processed[0], 64)
		if err != nil {
			return OrderedPair{}, err
		}
		minTemp, err := strconv.ParseFloat(processed[1], 64)
		if err != nil {
			return OrderedPair{}, err
		}

		result.x = maxTemp
		result.y = minTemp
	}

	return result, nil
}

// LoadWeatherData loads weather data from specific CSV files. It iterates through each file in the specified directory and checks if it is a file that should be processed.
// If so, it removes the ".csv" extension from the file name and processes the file using the ProcessFile function.
// The processed data is then stored in a map with the modified file name as the key.
// The function returns the map and any error that occurred during the process.
func LoadWeatherData(folder string) (map[string]OrderedPair, error) {
	directory := folder

	//List of specific files to process
	specificFiles := map[string]bool{
		"AZ.csv": true, "CT.csv": true, "DC.csv": true, "DE.csv": true, "IN.csv": true,
		"KS.csv": true, "KY.csv": true, "MA.csv": true, "MD.csv": true, "ME.csv": true,
		"MI.csv": true, "MO.csv": true, "NC.csv": true, "NJ.csv": true, "NM.csv": true,
		"NY.csv": true, "OH.csv": true, "OR.csv": true, "PA.csv": true, "RI.csv": true,
		"SC.csv": true, "UT.csv": true, "VA.csv": true, "VT.csv": true, "WV.csv": true,
	}

	weatherData := make(map[string]OrderedPair)

	// Read the directory
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err // Return error
	}

	// Process each file
	for _, file := range files {
		if specificFiles[file.Name()] {
			// Remove .csv extension from the file name
			fileNameWithoutExtension := strings.TrimSuffix(file.Name(), ".csv")

			// Process the file
			filePath := filepath.Join(directory, file.Name())
			processed, err := ProcessFile(filePath)
			if err != nil {
				fmt.Println("Error processing file", file.Name(), ":", err)
				continue
			}

			// Store the processed data with the modified file name
			weatherData[fileNameWithoutExtension] = processed
		}
	}

	return weatherData, nil // Return map and nil error
}

// InitializeQuadrants creates a 5x5 grid of Quadrants
// initializes the quadrants based on the maximum and minimum longitude and latitude values.
// It calculates the width and height of each quadrant.
// Then, it creates a list of quadrants, with each quadrant's properties such as its coordinates, width, and temperature.
// The code loads the weather data for each state and calculates the average temperature for each quadrant based on the state data.
// Finally, it returns the initialized weather data.
func InitializeQuadrants() Weather {
	totalWidth := maxLon - minLon
	totalHeight := maxLat - minLat

	quadrantWidth := totalWidth / 5
	quadrantHeight := totalHeight / 5

	var quadrants []Quadrant
	quadrantID := 1

	// First row:
	for i := 0; i < 5; i++ {
		quadrant := Quadrant{
			x:     minLon + float64(i)*quadrantWidth,
			y:     minLat + quadrantHeight*4,
			width: quadrantWidth,
			id:    quadrantID,
			temp:  0.0,
		}
		quadrants = append(quadrants, quadrant)
		quadrantID++
	}

	// Second row:
	for i := 0; i < 5; i++ {
		quadrant := Quadrant{
			x:     minLon + float64(i)*quadrantWidth,
			y:     minLat + quadrantHeight*3,
			width: quadrantWidth,
			id:    quadrantID,
			temp:  0.0,
		}
		quadrants = append(quadrants, quadrant)
		quadrantID++
	}

	// Third row:
	for i := 0; i < 5; i++ {
		quadrant := Quadrant{
			x:     minLon + float64(i)*quadrantWidth,
			y:     minLat + quadrantHeight*2,
			width: quadrantWidth,
			id:    quadrantID,
			temp:  0.0,
		}
		quadrants = append(quadrants, quadrant)
		quadrantID++
	}

	// Fourth row:
	for i := 0; i < 5; i++ {
		quadrant := Quadrant{
			x:     minLon + float64(i)*quadrantWidth,
			y:     minLat + quadrantHeight,
			width: quadrantWidth,
			id:    quadrantID,
			temp:  0.0,
		}
		quadrants = append(quadrants, quadrant)
		quadrantID++
	}

	// Fifth row:
	for i := 0; i < 5; i++ {
		quadrant := Quadrant{
			x:     minLon + float64(i)*quadrantWidth,
			y:     minLat,
			width: quadrantWidth,
			id:    quadrantID,
			temp:  0.0,
		}
		quadrants = append(quadrants, quadrant)
		quadrantID++
	}

	// Quadrants are ordered from:

	// Northwest to Northeast (Top Row): OR, CT, VT, MA, ME
	// Second Row: NJ, NY, PA, MD, DE
	// Third Row (Middle Row): WV, VA, NC, DC, SC
	// Fourth Row: NM, MO, IN, OH, MI
	// Southwest to Southeast (Bottom Row): AZ, UT, KY, RI, KS

	weatherData, err := LoadWeatherData("Data/Hatch_May-Jun")
	if err != nil {
		fmt.Println("Error loading weather data:", err)

	}

	// 	Northwest to Northeast (Top Row)
	quadrants[0].temp += FareinheitToCelsius(weatherData["OR"].x)
	quadrants[1].temp += FareinheitToCelsius(weatherData["CT"].x)
	quadrants[2].temp += FareinheitToCelsius(weatherData["VT"].x)
	quadrants[3].temp += FareinheitToCelsius(weatherData["MA"].x)
	quadrants[4].temp += FareinheitToCelsius(weatherData["ME"].x)

	// Second Row
	quadrants[5].temp += FareinheitToCelsius(weatherData["NJ"].x)
	quadrants[6].temp += FareinheitToCelsius(weatherData["NY"].x)
	quadrants[7].temp += FareinheitToCelsius(weatherData["PA"].x)
	quadrants[8].temp += FareinheitToCelsius(weatherData["MD"].x)
	quadrants[9].temp += FareinheitToCelsius(weatherData["DE"].x)

	// Third Row (Middle Row)
	quadrants[10].temp += FareinheitToCelsius(weatherData["WV"].x)
	quadrants[11].temp += FareinheitToCelsius(weatherData["VA"].x)
	quadrants[12].temp += FareinheitToCelsius(weatherData["NC"].x)
	quadrants[13].temp += FareinheitToCelsius(weatherData["DC"].x)
	quadrants[14].temp += FareinheitToCelsius(weatherData["SC"].x)

	// Fourth Row
	quadrants[15].temp += FareinheitToCelsius(weatherData["NM"].x)
	quadrants[16].temp += FareinheitToCelsius(weatherData["MO"].x)
	quadrants[17].temp += FareinheitToCelsius(weatherData["IN"].x)
	quadrants[18].temp += FareinheitToCelsius(weatherData["OH"].x)
	quadrants[19].temp += FareinheitToCelsius(weatherData["MI"].x)

	// Southwest to Southeast (Bottom Row)
	quadrants[20].temp += FareinheitToCelsius(weatherData["AZ"].x)
	quadrants[21].temp += FareinheitToCelsius(weatherData["UT"].x)
	quadrants[22].temp += FareinheitToCelsius(weatherData["KY"].x)
	quadrants[23].temp += FareinheitToCelsius(weatherData["RI"].x)
	quadrants[24].temp += FareinheitToCelsius(weatherData["KS"].x)

	return Weather{
		x:         minLon,
		y:         minLat,
		Quadrants: quadrants,
	}
}

// ReadTrees reads data from a CSV file, where each row contains the longitude and latitude of a habitat for a tree.
// It then stores the data as an array of OrderedPair objects, with each object representing a unique pair of coordinates.
// The function returns the array of OrderedPair objects and an error, if any occurred during the process.
func ReadTrees(filePath string) ([]OrderedPair, error) {
	var habitats []OrderedPair

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for i, record := range records {
		if i == 0 { // Skip header
			continue
		}

		// Trim spaces and parse Latitude
		latitudeStr := strings.TrimSpace(record[0])
		latitude, err := strconv.ParseFloat(latitudeStr, 64)
		if err != nil {
			fmt.Printf("Error parsing latitude in row %d: %v\n", i+1, err)
			continue
		}

		// Trim spaces and parse latitude
		longitudeStr := strings.TrimSpace(record[1])
		longitude, err := strconv.ParseFloat(longitudeStr, 64)
		if err != nil {
			fmt.Printf("Error parsing longitude in row %d: %v\n", i+1, err)
			continue
		}

		habitats = append(habitats, OrderedPair{x: longitude, y: latitude})
	}

	return habitats, nil
}

// InitialiCountry  is responsible for setting up and initializing a Country object, representing a geographical region.
// The country has certain attributes, including its width and height, as well as a collection of trees and flies.
func InitializeCountry() Country {
	var country Country
	country.width = maxLat - minLat
	country.height = maxLon - minLon

	// Initialize trees
	tree, err := ReadTrees("Data/processed_data.csv")
	if err != nil {
		fmt.Println("Error loading weather data:", err)

	}

	numberOfTree := len(tree)
	country.trees = make([]Tree, numberOfTree)

	for i := 0; i < numberOfTree; i++ {
		country.trees[i].position.x = tree[i].x
		country.trees[i].position.y = tree[i].y
	}

	// Remove trees out of range
	for i := 0; i < numberOfTree; i++ {
		if country.trees[i].position.y < minLat || country.trees[i].position.y > maxLat ||
			country.trees[i].position.x < minLon || country.trees[i].position.x > maxLon {
			country.trees = append(country.trees[:i], country.trees[i+1:]...)
			i--
			numberOfTree--
		}
	}

	// Initialize flies
	flies, err := ReadSampleDataFromFile("Data/lydetext.txt")
	if err != nil {
		fmt.Println("Error loading weather data:", err)
	}

	numberOfFlies := len(flies)

	country.flies = make([]Fly, numberOfFlies)

	// Initialize flies
	for i := 0; i < numberOfFlies; i++ {
		country.flies[i].position.x = flies[i].RoundedLongitude
		country.flies[i].position.y = flies[i].RoundedLatitude
		country.flies[i].stage = 0
	}

	// Only 10% of the flies are alive

	for i := 0; i < numberOfFlies; i++ {
		if rand.Float64() < 0.1 {
			country.flies[i].isAlive = true
		} else {
			country.flies[i].isAlive = false
		}
	}

	// Remove flies that are not alive from the list
	for i := 0; i < numberOfFlies; i++ {
		if country.flies[i].isAlive == false {
			country.flies = append(country.flies[:i], country.flies[i+1:]...)
			i--
			numberOfFlies--
		}
	}

	// Add Location ID:
	for i := 0; i < numberOfFlies; i++ {
		if country.flies[i].position.y > 42.834 && country.flies[i].position.x > -123.27 && country.flies[i].position.y < -112.402 {
			country.flies[i].locationID = 1
		} else if country.flies[i].position.y > 42.834 && country.flies[i].position.x > -112.402 && country.flies[i].position.x < -101.53399999999999 {
			country.flies[i].locationID = 2
		} else if country.flies[i].position.y > 42.834 && country.flies[i].position.x > -101.53399999999999 && country.flies[i].position.x < -90.666 {
			country.flies[i].locationID = 3
		} else if country.flies[i].position.y > 42.834 && country.flies[i].position.x > -90.666 && country.flies[i].position.x < -79.798 {
			country.flies[i].locationID = 4
		} else if country.flies[i].position.y > 42.834 && country.flies[i].position.x > -79.798 && country.flies[i].position.x < maxLat {
			country.flies[i].locationID = 5
		} else if country.flies[i].position.y > 39.958 && country.flies[i].position.x > -123.27 && country.flies[i].position.y < -112.402 {
			country.flies[i].locationID = 6
		} else if country.flies[i].position.y > 39.958 && country.flies[i].position.x > -112.402 && country.flies[i].position.x < -101.53399999999999 {
			country.flies[i].locationID = 7
		} else if country.flies[i].position.y > 39.958 && country.flies[i].position.x > -101.53399999999999 && country.flies[i].position.x < -90.666 {
			country.flies[i].locationID = 8
		} else if country.flies[i].position.y > 39.958 && country.flies[i].position.x > -90.666 && country.flies[i].position.x < -79.798 {
			country.flies[i].locationID = 9
		} else if country.flies[i].position.y > 39.958 && country.flies[i].position.x > -79.798 && country.flies[i].position.x < maxLat {
			country.flies[i].locationID = 10
		} else if country.flies[i].position.y > 37.082 && country.flies[i].position.x > -123.27 && country.flies[i].position.y < -112.402 {
			country.flies[i].locationID = 11
		} else if country.flies[i].position.y > 37.082 && country.flies[i].position.x > -112.402 && country.flies[i].position.x < -101.53399999999999 {
			country.flies[i].locationID = 12
		} else if country.flies[i].position.y > 37.082 && country.flies[i].position.x > -101.53399999999999 && country.flies[i].position.x < -90.666 {
			country.flies[i].locationID = 13
		} else if country.flies[i].position.y > 37.082 && country.flies[i].position.x > -90.666 && country.flies[i].position.x < -79.798 {
			country.flies[i].locationID = 14
		} else if country.flies[i].position.y > 37.082 && country.flies[i].position.x > -79.798 && country.flies[i].position.x < maxLat {
			country.flies[i].locationID = 15
		} else if country.flies[i].position.y > 34.205999999999996 && country.flies[i].position.x > -123.27 && country.flies[i].position.y < -112.402 {
			country.flies[i].locationID = 16
		} else if country.flies[i].position.y > 34.205999999999996 && country.flies[i].position.x > -112.402 && country.flies[i].position.x < -101.53399999999999 {
			country.flies[i].locationID = 17
		} else if country.flies[i].position.y > 34.205999999999996 && country.flies[i].position.x > -101.53399999999999 && country.flies[i].position.x < -90.666 {
			country.flies[i].locationID = 18
		} else if country.flies[i].position.y > 34.205999999999996 && country.flies[i].position.x > -90.666 && country.flies[i].position.x < -79.798 {
			country.flies[i].locationID = 19
		} else if country.flies[i].position.y > 34.205999999999996 && country.flies[i].position.x > -79.798 && country.flies[i].position.x < maxLat {
			country.flies[i].locationID = 20
		} else if country.flies[i].position.y > 31.33 && country.flies[i].position.x > -123.27 && country.flies[i].position.y < -112.402 {
			country.flies[i].locationID = 21
		} else if country.flies[i].position.y > 31.33 && country.flies[i].position.x > -112.402 && country.flies[i].position.x < -101.53399999999999 {
			country.flies[i].locationID = 22
		} else if country.flies[i].position.y > 31.33 && country.flies[i].position.x > -101.53399999999999 && country.flies[i].position.x < -90.666 {
			country.flies[i].locationID = 23
		} else if country.flies[i].position.y > 31.33 && country.flies[i].position.x > -90.666 && country.flies[i].position.x < -79.798 {
			country.flies[i].locationID = 24
		} else if country.flies[i].position.y > 31.33 && country.flies[i].position.x > -79.798 && country.flies[i].position.x < maxLat {
			country.flies[i].locationID = 25
		}
	}

	weatherData, err := LoadWeatherData("Data/Hatch_May-Jun")

	if err != nil {
		fmt.Println("Error loading weather data:", err)
	}

	// Add the energy of flies:

	for i := 0; i < numberOfFlies; i++ {
		if country.flies[i].locationID == 1 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["OR"].x, weatherData["OR"].y))
		} else if country.flies[i].locationID == 2 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["CT"].x, weatherData["CT"].y))
		} else if country.flies[i].locationID == 3 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["VT"].x, weatherData["VT"].y))
		} else if country.flies[i].locationID == 4 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["MA"].x, weatherData["MA"].y))
		} else if country.flies[i].locationID == 5 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["ME"].x, weatherData["ME"].y))
		} else if country.flies[i].locationID == 6 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["NJ"].x, weatherData["NJ"].y))
		} else if country.flies[i].locationID == 7 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["NY"].x, weatherData["NY"].y))
		} else if country.flies[i].locationID == 8 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["PA"].x, weatherData["PA"].y))
		} else if country.flies[i].locationID == 9 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["MD"].x, weatherData["MD"].y))
		} else if country.flies[i].locationID == 10 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["DE"].x, weatherData["DE"].y))
		} else if country.flies[i].locationID == 11 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["WV"].x, weatherData["WV"].y))
		} else if country.flies[i].locationID == 12 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["VA"].x, weatherData["VA"].y))
		} else if country.flies[i].locationID == 13 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["NC"].x, weatherData["NC"].y))
		} else if country.flies[i].locationID == 14 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["DC"].x, weatherData["DC"].y))
		} else if country.flies[i].locationID == 15 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["SC"].x, weatherData["SC"].y))
		} else if country.flies[i].locationID == 16 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["NM"].x, weatherData["NM"].y))
		} else if country.flies[i].locationID == 17 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["MO"].x, weatherData["MO"].y))
		} else if country.flies[i].locationID == 18 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["IN"].x, weatherData["IN"].y))
		} else if country.flies[i].locationID == 19 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["OH"].x, weatherData["OH"].y))
		} else if country.flies[i].locationID == 20 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["MI"].x, weatherData["MI"].y))
		} else if country.flies[i].locationID == 21 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["AZ"].x, weatherData["AZ"].y))
		} else if country.flies[i].locationID == 22 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["UT"].x, weatherData["UT"].y))
		} else if country.flies[i].locationID == 23 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["KY"].x, weatherData["KY"].y))
		} else if country.flies[i].locationID == 24 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["RI"].x, weatherData["RI"].y))
		} else if country.flies[i].locationID == 25 {
			country.flies[i].energy += FareinheitToCelsius(randomInRange(weatherData["KS"].x, weatherData["KS"].y))
		}
	}

	return country
}

// RandomInRange returns a random float64 in range [min, max).
// uses the rand package to generate a random float64 within the specified range.
// It seeds the random number generator with the current Unix time in nanoseconds to ensure that it generates a new random sequence every time it's called.
func randomInRange(max, min float64) float64 {
	rand.Seed(time.Now().UnixNano()) // Seed random number generator
	return min + rand.Float64()*(max-min)
}

// Fahrenheit to Celsius takes a float64 value f as an argument.
// The function calculates the equivalent temperature in Celsius and returns it.
func FareinheitToCelsius(f float64) float64 {
	// Convert to Celsius
	return (f - 32) * 5 / 9
}
