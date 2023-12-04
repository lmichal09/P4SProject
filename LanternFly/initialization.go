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

type OrderedPair struct {
	x float64
	y float64
}

const (
	minLat float64 = 31.33   // Southernmost point in the US
	maxLat float64 = 45.71   // Northernmost point in the contiguous US
	minLon float64 = -123.27 // Westernmost point in the contiguous US
	maxLon float64 = -68.93  // Easternmost point in the contiguous US
)

type Weather struct {
	x         float64 // Bottom left corner x coordinate (Longitude)
	y         float64 // Bottom left corner y coordinate (Latitude)
	Quadrants []Quadrant
}

type Quadrant struct {
	x     float64 // Bottom left corner x coordinate (Longitude)
	y     float64 // Bottom left corner y coordinate (Latitude)
	width float64
	id    int
	temp  float64
}

type Country struct {
	width float64
	flies []Fly
	trees []Tree
}

type Fly struct {
	position, velocity, acceleration OrderedPair
	stage                            int     // 0 = egg, 1 = instar1, 2 = instar2, 3 = instar3, 4 = instar4, 5 = adult
	energy                           float64 // Degree-days
	isAlive                          bool
	locationID                       int
}

type Tree struct {
	position OrderedPair
}

type SampleData struct {
	Source           string
	Year             int
	BioYear          int
	Latitude         float64
	Longitude        float64
	State            string
	LydePresent      bool
	LydeEstablished  bool
	LydeDensity      string
	SourceAgency     string
	CollectionMethod string
	PointID          string
	RoundedLongitude float64
	RoundedLatitude  float64
}

/*
Set Time: 2021
25 States: "PA" "NJ" "VA" "DE" "MD" "NY" "UT" "MA" "MI" "NC" "WV" "CT" "VT" "OH" "IN" "KY" "DC" "SC" "NM" "AZ" "RI" "OR" "MO" "KS" "ME"

The life cycle of the lanternfly is as follows: 6 stages 0-5, stage 7: Dead
	Time line of the lanternfly is as follows:
	Eggs can be found on any outdoor surface from October through June.
	May-June: Eggs hatch into 1st instar nymphs.
	May-July: 1st instar nymphs molt into 2nd instar nymphs, 2nd instar nymphs molt into 3rd instar nymphs.
	july-September: 3rd instar nymphs molt into 4th instar nymphs.
	July-December: 4th instar nymphs molt into adults.
	September-November: Adults lay eggs.
	Oct-June: Eggs overwinter.
*/

// width and population need to be read from file

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

func LoadWeatherData(directory string) (map[string]OrderedPair, error) {
	directory = directory

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
	totalHeight := maxLon - minLon

	quadrantWidth := totalWidth / 5
	quadrantHeight := totalHeight / 5

	var quadrants []Quadrant
	quadrantID := 1

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			quadrant := Quadrant{
				x:     minLon + float64(j)*quadrantWidth,
				y:     minLat + float64(i)*quadrantHeight,
				width: quadrantWidth,
				id:    quadrantID,
				temp:  0.0,
			}
			quadrants = append(quadrants, quadrant)
			quadrantID++
		}
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

	tree, err := ReadTrees("Data/processed_data.csv")
	if err != nil {
		fmt.Println("Error loading weather data:", err)

	}

	numberOfTree := len(tree)
	country.trees = make([]Tree, numberOfTree)

	// Initialize trees
	for i := 0; i < numberOfTree; i++ {
		country.trees[i].position.x = tree[i].x
		country.trees[i].position.y = tree[i].y
	}

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
		country.flies[i].velocity.x = rand.Float64() * 0.1
		country.flies[i].velocity.y = rand.Float64() * 0.2
		country.flies[i].acceleration.x = rand.Float64() * rand.Float64() * 0.1
		country.flies[i].acceleration.y = rand.Float64() * rand.Float64() * 0.2
		country.flies[i].stage = 0
		country.flies[i].isAlive = true

	}

	// Add Location ID:
	for i := 0; i < numberOfFlies; i++ {
		if country.flies[i].position.x < -123.270000 && country.flies[i].position.y < 31.330000 {
			country.flies[i].locationID = 1
		} else if country.flies[i].position.x < -112.402000 && country.flies[i].position.y < 31.330000 {
			country.flies[i].locationID = 2
		} else if country.flies[i].position.x < -101.534000 && country.flies[i].position.y < 31.330000 {
			country.flies[i].locationID = 3
		} else if country.flies[i].position.x < -90.666000 && country.flies[i].position.y < 31.330000 {
			country.flies[i].locationID = 4
		} else if country.flies[i].position.x < -79.798000 && country.flies[i].position.y < 31.330000 {
			country.flies[i].locationID = 5
		} else if country.flies[i].position.x < -123.270000 && country.flies[i].position.y < 42.198000 {
			country.flies[i].locationID = 6
		} else if country.flies[i].position.x < -112.402000 && country.flies[i].position.y < 42.198000 {
			country.flies[i].locationID = 7
		} else if country.flies[i].position.x < -101.534000 && country.flies[i].position.y < 42.198000 {
			country.flies[i].locationID = 8
		} else if country.flies[i].position.x < -90.666000 && country.flies[i].position.y < 42.198000 {
			country.flies[i].locationID = 9
		} else if country.flies[i].position.x < -79.798000 && country.flies[i].position.y < 42.198000 {
			country.flies[i].locationID = 10
		} else if country.flies[i].position.x < -123.270000 && country.flies[i].position.y < 53.066000 {
			country.flies[i].locationID = 11
		} else if country.flies[i].position.x < -112.402000 && country.flies[i].position.y < 53.066000 {
			country.flies[i].locationID = 12
		} else if country.flies[i].position.x < -101.534000 && country.flies[i].position.y < 53.066000 {
			country.flies[i].locationID = 13
		} else if country.flies[i].position.x < -90.666000 && country.flies[i].position.y < 53.066000 {
			country.flies[i].locationID = 14
		} else if country.flies[i].position.x < -79.798000 && country.flies[i].position.y < 53.066000 {
			country.flies[i].locationID = 15
		} else if country.flies[i].position.x < -123.270000 && country.flies[i].position.y < 63.934000 {
			country.flies[i].locationID = 16
		} else if country.flies[i].position.x < -112.402000 && country.flies[i].position.y < 63.934000 {
			country.flies[i].locationID = 17
		} else if country.flies[i].position.x < -101.534000 && country.flies[i].position.y < 63.934000 {
			country.flies[i].locationID = 18
		} else if country.flies[i].position.x < -90.666000 && country.flies[i].position.y < 63.934000 {
			country.flies[i].locationID = 19
		} else if country.flies[i].position.x < -79.798000 && country.flies[i].position.y < 63.934000 {
			country.flies[i].locationID = 20
		} else if country.flies[i].position.x < -123.270000 && country.flies[i].position.y < 74.802000 {
			country.flies[i].locationID = 21
		} else if country.flies[i].position.x < -112.402000 && country.flies[i].position.y < 74.802000 {
			country.flies[i].locationID = 22
		} else if country.flies[i].position.x < -101.534000 && country.flies[i].position.y < 74.802000 {
			country.flies[i].locationID = 23
		} else if country.flies[i].position.x < -90.666000 && country.flies[i].position.y < 74.802000 {
			country.flies[i].locationID = 24
		} else if country.flies[i].position.x < -79.798000 && country.flies[i].position.y < 74.802000 {
			country.flies[i].locationID = 25
		}
	}

	weatherData, err := LoadWeatherData("data")

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

func randomInRange(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano()) // Initialize the random number generator
	return min + rand.Float64()*(max-min)
}

func FareinheitToCelsius(f float64) float64 {
	return (f - 32) * 5 / 9
}

func main() {
	quadrants := InitializeQuadrants()
	for _, quadrant := range quadrants.Quadrants {
		fmt.Printf("Quadrant ID: %d, Bottom Left (X, Y): (%.6f, %.6f), Width: %.6f, Temp: %.6f\n",
			quadrant.id, quadrant.x, quadrant.y, quadrant.width, quadrant.temp)
	}

	country := InitializeCountry()
	for _, fly := range country.flies {
		if fly.locationID == 1 {
			fmt.Printf("Location ID: %d, Energy: %.6f\n, Position X: %.6f, PositionY: %.6f\n",
				fly.locationID, fly.energy, fly.position.x, fly.position.y)
		}
	}

}
