package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ProcessFile processes the CSV file, removing specific columns and keeping certain rows.
// output of 1 state and its maxTemp and minTemp

func ProcessFile(filePath string) (OrderedPair, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return OrderedPair{}, err
	}
	defer file.Close()

	// Read the file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return OrderedPair{}, err
	}

	// Process the file
	var processed []string
	for i, record := range records {
		if i == 5 || i == 12 { // Keeping only rows 5 and 12
			if len(record) > 1 {
				processed = append(processed, record[1]) // Keeping only the second column
			}
		}
	}

	// Store the processed data
	var result OrderedPair
	if len(processed) >= 2 {
		// Parse the string values into float64
		maxTemp, err := strconv.ParseFloat(processed[0], 64)
		if err != nil {
			return OrderedPair{}, err
		}
		// Parse the string values into float64
		minTemp, err := strconv.ParseFloat(processed[1], 64)
		if err != nil {
			return OrderedPair{}, err
		}

		// Store the processed data into result OrderPair
		result.x = maxTemp
		result.y = minTemp
	}

	return result, nil
}

// LoadWeatherData loads the weather data from the specified directory.
// output: map of state and its maxTemp and minTemp
func LoadWeatherData(directory string) (map[string]OrderedPair, error) {
	directory = directory

	//List of specific files to process
	specificFiles := map[string]bool{
		"AZ.csv": true, "CT.csv": true, "DC.csv": true, "DE.csv": true, "IN.csv": true, "KS.csv": true, "KY.csv": true, "MA.csv": true, "MD.csv": true, "ME.csv": true, "MI.csv": true, "MO.csv": true, "NC.csv": true, "NJ.csv": true, "NM.csv": true, "NY.csv": true, "OH.csv": true, "OR.csv": true, "PA.csv": true, "RI.csv": true, "SC.csv": true, "UT.csv": true, "VA.csv": true, "VT.csv": true, "WV.csv": true,
	}

	// Create a map to store the processed data
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
