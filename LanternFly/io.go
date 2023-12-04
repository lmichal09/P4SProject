package main

// import (
// 	"bufio"
// 	"encoding/csv"
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
// 	"strconv"
// 	"strings"
// )

// // ReadSamplesFromDirectory reads a collection of files from a given directory
// // and returns a map where the keys are sample names and the values are slices of SampleData.
// func ReadSamplesFromDirectory(directory string) map[string][]SampleData {
// 	allData := make(map[string][]SampleData)

// 	dirContents, err := ioutil.ReadDir(directory)
// 	if err != nil {
// 		panic("Error reading directory!")
// 	}

// 	for _, fileData := range dirContents {
// 		// Construct the full file path
// 		filePath := filepath.Join(directory, fileData.Name())

// 		// Read data from file
// 		data := ReadSampleDataFromFile(filePath)

// 		// Use PointID as the key for the map
// 		for _, sample := range data {
// 			key := sample.PointID
// 			allData[key] = append(allData[key], sample)
// 		}
// 	}

// 	return allData
// }

// // ReadSampleDataFromFile reads data from a file and returns a slice of SampleData.
// func ReadSampleDataFromFile(filename string) []SampleData {
// 	var sampleData []SampleData

// 	file, err := os.Open(filename)
// 	defer file.Close()
// 	if err != nil {
// 		panic(err)
// 	}

// 	scanner := bufio.NewScanner(file)

// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		fields := strings.Split(line, "\t")

// 		// Convert fields to appropriate types
// 		year := parseInteger(fields[1])
// 		bioYear := parseInteger(fields[2])
// 		latitude := parseFloat(fields[3])
// 		longitude := parseFloat(fields[4])
// 		lydePresent := parseBool(fields[6])
// 		lydeEstablished := parseBool(fields[7])
// 		roundedLongitude := parseFloat(fields[12])
// 		roundedLatitude := parseFloat(fields[13])

// 		// Create a SampleData instance and append to the slice
// 		data := SampleData{
// 			Source:           fields[0],
// 			Year:             year,
// 			BioYear:          bioYear,
// 			Latitude:         latitude,
// 			Longitude:        longitude,
// 			State:            fields[5],
// 			LydePresent:      lydePresent,
// 			LydeEstablished:  lydeEstablished,
// 			LydeDensity:      fields[8],
// 			SourceAgency:     fields[9],
// 			CollectionMethod: fields[10],
// 			PointID:          fields[11],
// 			RoundedLongitude: roundedLongitude,
// 			RoundedLatitude:  roundedLatitude,
// 		}

// 		sampleData = append(sampleData, data)
// 	}

// 	err = scanner.Err()
// 	if err != nil {
// 		panic(err)
// 	}

// 	return sampleData
// }

// // WriteToFile writes the sample data to a CSV file.
// func WriteToFile(filename string, data []SampleData) error {
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()

// 	// Writing header
// 	header := []string{
// 		"Source", "Year", "BioYear", "Latitude", "Longitude",
// 		"State", "LydePresent", "LydeEstablished", "LydeDensity",
// 		"SourceAgency", "CollectionMethod", "PointID",
// 		"RoundedLongitude", "RoundedLatitude",
// 	}
// 	if err := writer.Write(header); err != nil {
// 		return err
// 	}

// 	// Writing data
// 	for _, sample := range data {
// 		record := []string{
// 			sample.Source, strconv.Itoa(sample.Year), strconv.Itoa(sample.BioYear),
// 			strconv.FormatFloat(sample.Latitude, 'f', -1, 64),
// 			strconv.FormatFloat(sample.Longitude, 'f', -1, 64),
// 			sample.State, strconv.FormatBool(sample.LydePresent),
// 			strconv.FormatBool(sample.LydeEstablished), sample.LydeDensity,
// 			sample.SourceAgency, sample.CollectionMethod, sample.PointID,
// 			strconv.FormatFloat(sample.RoundedLongitude, 'f', -1, 64),
// 			strconv.FormatFloat(sample.RoundedLatitude, 'f', -1, 64),
// 		}

// 		if err := writer.Write(record); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// // Helper function to convert string to integer
// func parseInteger(s string) int {
// 	val, err := strconv.Atoi(s)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return val
// }

// // Helper function to convert string to float64
// func parseFloat(s string) float64 {
// 	val, err := strconv.ParseFloat(s, 64)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return val
// }

// // Helper function to convert string to bool
// func parseBool(s string) bool {
// 	val, err := strconv.ParseBool(s)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return val
// }
