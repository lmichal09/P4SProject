package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func parseInteger(s string) (int, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("error parsing integer: %v", err)
	}
	return val, nil
}

func parseFloat(s string) (float64, error) {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing float: %v", err)
	}
	return val, nil
}

func parseBool(s string) (bool, error) {
	val, err := strconv.ParseBool(s)
	if err != nil {
		return false, fmt.Errorf("error parsing bool: %v", err)
	}
	return val, nil
}

func parseString(s string) string {
	return strings.TrimSpace(s)
}

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

func ProcessFile(filePath string) (OrderedPair, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return OrderedPair{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return OrderedPair{}, err
	}

	var processed []string
	for i, record := range records {
		if i == 5 || i == 12 { // Keeping only rows 5 and 12
			if len(record) > 1 {
				processed = append(processed, record[1]) // Keeping only the second column
			}
		}
	}

	var result OrderedPair
	if len(processed) >= 2 {
		// Parse the string values into float64
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

func LoadWeatherData(folder string) (map[string]OrderedPair, error) {
	directory := folder

	//List of specific files to process
	specificFiles := map[string]bool{
		"AZ.csv": true, "CT.csv": true, "DC.csv": true, "DE.csv": true, "IN.csv": true, "KS.csv": true, "KY.csv": true, "MA.csv": true, "MD.csv": true, "ME.csv": true, "MI.csv": true, "MO.csv": true, "NC.csv": true, "NJ.csv": true, "NM.csv": true, "NY.csv": true, "OH.csv": true, "OR.csv": true, "PA.csv": true, "RI.csv": true, "SC.csv": true, "UT.csv": true, "VA.csv": true, "VT.csv": true, "WV.csv": true,
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
func InitializeQuadrants() Weather {
	totalWidth := maxLon - minLon
	//totalHeight := maxLat - minLat

	quadrantWidth := totalWidth / 5
	//quadrantHeight := totalHeight/ 5

	var quadrants []Quadrant
	quadrantID := 1

	// First row:
	for i := 0; i < 5; i++ {
		quadrant := Quadrant{
			x:     minLat + (maxLat-minLat)/5*4,
			y:     minLon + float64(i)*quadrantWidth,
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
			x:     minLat + (maxLat-minLat)/5*3,
			y:     minLon + float64(i)*quadrantWidth,
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
			x:     minLat + (maxLat-minLat)/5*2,
			y:     minLon + float64(i)*quadrantWidth,
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
			x:     minLat + (maxLat-minLat)/5*1,
			y:     minLon + float64(i)*quadrantWidth,
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
			x:     minLat + (maxLat-minLat)/5*0,
			y:     minLon + float64(i)*quadrantWidth,
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

		// Trim spaces and parse longitude
		longitudeStr := strings.TrimSpace(record[0])
		longitude, err := strconv.ParseFloat(longitudeStr, 64)
		if err != nil {
			fmt.Printf("Error parsing longitude in row %d: %v\n", i+1, err)
			continue
		}

		// Trim spaces and parse latitude
		latitudeStr := strings.TrimSpace(record[1])
		latitude, err := strconv.ParseFloat(latitudeStr, 64)
		if err != nil {
			fmt.Printf("Error parsing latitude in row %d: %v\n", i+1, err)
			continue
		}

		habitats = append(habitats, OrderedPair{x: longitude, y: latitude})
	}

	return habitats, nil
}

func InitializeCountry() Country {
	var country Country
	country.width = maxLon - minLon
	country.height = maxLat - minLat

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

	// Convert into equirectangular projection
	for i := 0; i < numberOfTree; i++ {
		country.trees[i].position = equirectangularProjection(country.trees[i].position)
	}

	// Remove trees out of range
	for i := 0; i < numberOfTree; i++ {
		if country.trees[i].position.x < minLat || country.trees[i].position.x > maxLat || country.trees[i].position.y < minLon || country.trees[i].position.y > maxLon {
			// append
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
		country.flies[i].position.x = flies[i].RoundedLatitude
		country.flies[i].position.y = flies[i].RoundedLongitude
		country.flies[i].stage = 0
	}

	// Change the positions to equirectangular projection

	for i := 0; i < numberOfFlies; i++ {
		country.flies[i].position = equirectangularProjection(country.flies[i].position)
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
		if country.flies[i].position.x > -6623.397249337799 && country.flies[i].position.y > 3483.737051774025 && country.flies[i].position.y < 3803.533660803776 {
			country.flies[i].locationID = 1
		} else if country.flies[i].position.x > -6623.397249337799 && country.flies[i].position.y > 3803.533660803776 && country.flies[i].position.y < 4123.330269833527 {
			country.flies[i].locationID = 2
		} else if country.flies[i].position.x > -6623.397249337799 && country.flies[i].position.y > 4123.330269833527 && country.flies[i].position.y < 4443.126878863278 {
			country.flies[i].locationID = 3
		} else if country.flies[i].position.x > -6623.397249337799 && country.flies[i].position.y > 4443.126878863278 && country.flies[i].position.y < 4762.923487893029 {
			country.flies[i].locationID = 4
		} else if country.flies[i].position.x > -6623.397249337799 && country.flies[i].position.y > 4762.923487893029 && country.flies[i].position.y < 5082.720096922779 {
			country.flies[i].locationID = 5
		} else if country.flies[i].position.x > -7894.6318914036865 && country.flies[i].position.y > 3483.737051774025 && country.flies[i].position.y < 3803.533660803776 {
			country.flies[i].locationID = 6
		} else if country.flies[i].position.x > -7894.6318914036865 && country.flies[i].position.y > 3803.533660803776 && country.flies[i].position.y < 4123.330269833527 {
			country.flies[i].locationID = 7
		} else if country.flies[i].position.x > -7894.6318914036865 && country.flies[i].position.y > 4123.330269833527 && country.flies[i].position.y < 4443.126878863278 {
			country.flies[i].locationID = 8
		} else if country.flies[i].position.x > -7894.6318914036865 && country.flies[i].position.y > 4443.126878863278 && country.flies[i].position.y < 4762.923487893029 {
			country.flies[i].locationID = 9
		} else if country.flies[i].position.x > -7894.6318914036865 && country.flies[i].position.y > 4762.923487893029 && country.flies[i].position.y < 5082.720096922779 {
			country.flies[i].locationID = 10
		} else if country.flies[i].position.x > -9165.8665334695745 && country.flies[i].position.y > 3483.737051774025 && country.flies[i].position.y < 3803.533660803776 {
			country.flies[i].locationID = 11
		} else if country.flies[i].position.x > -9165.866533469574 && country.flies[i].position.y > 3803.533660803776 && country.flies[i].position.y < 4123.330269833527 {
			country.flies[i].locationID = 12
		} else if country.flies[i].position.x > -9165.866533469574 && country.flies[i].position.y > 4123.330269833527 && country.flies[i].position.y < 4443.126878863278 {
			country.flies[i].locationID = 13
		} else if country.flies[i].position.x > -9165.866533469574 && country.flies[i].position.y > 4443.126878863278 && country.flies[i].position.y < 4762.923487893029 {
			country.flies[i].locationID = 14
		} else if country.flies[i].position.x > -9165.866533469574 && country.flies[i].position.y > 4762.923487893029 && country.flies[i].position.y < 5082.720096922779 {
			country.flies[i].locationID = 15
		} else if country.flies[i].position.x > -10437.101175535463 && country.flies[i].position.y > 3483.737051774025 && country.flies[i].position.y < 3803.533660803776 {
			country.flies[i].locationID = 16
		} else if country.flies[i].position.x > -10437.101175535463 && country.flies[i].position.y > 3803.533660803776 && country.flies[i].position.y < 4123.330269833527 {
			country.flies[i].locationID = 17
		} else if country.flies[i].position.x > -10437.101175535463 && country.flies[i].position.y > 4123.330269833527 && country.flies[i].position.y < 4443.126878863278 {
			country.flies[i].locationID = 18
		} else if country.flies[i].position.x > -10437.1011755354639 && country.flies[i].position.y > 4443.126878863278 && country.flies[i].position.y < 4762.923487893029 {
			country.flies[i].locationID = 19
		} else if country.flies[i].position.x > -10437.101175535463 && country.flies[i].position.y > 4762.923487893029 && country.flies[i].position.y < 5082.720096922779 {
			country.flies[i].locationID = 20
		} else if country.flies[i].position.x > -10437.101175535463 && country.flies[i].position.y > 3483.737051774025 && country.flies[i].position.y < 3803.533660803776 {
			country.flies[i].locationID = 21
		} else if country.flies[i].position.x > -10437.101175535463 && country.flies[i].position.y > 3803.533660803776 && country.flies[i].position.y < 4123.330269833527 {
			country.flies[i].locationID = 22
		} else if country.flies[i].position.x > -10437.101175535463 && country.flies[i].position.y > 4123.330269833527 && country.flies[i].position.y < 4443.126878863278 {
			country.flies[i].locationID = 23
		} else if country.flies[i].position.x > -10437.101175535463 && country.flies[i].position.y > 4443.126878863278 && country.flies[i].position.y < 4762.923487893029 {
			country.flies[i].locationID = 24
		} else if country.flies[i].position.x > -10437.101175535463 && country.flies[i].position.y > 4762.923487893029 && country.flies[i].position.y < 5082.720096922779 {
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

func randomInRange(max, min float64) float64 {
	rand.Seed(time.Now().UnixNano()) // Initialize the random number generator
	return min + rand.Float64()*(max-min)
}

func FareinheitToCelsius(f float64) float64 {
	return (f - 32) * 5 / 9
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func equirectangularProjection(position OrderedPair) OrderedPair {
	latRad := toRadians(position.x)
	lonRad := toRadians(position.y)

	x := earthRadius * lonRad * math.Cos(latRad)
	y := earthRadius * latRad
	return OrderedPair{x: x, y: y}
}
