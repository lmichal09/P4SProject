package main

import (
	"math/rand"
	"time"
)

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

const (
	instar1To2Threshold     float64 = 166.6
	instar2To3Threshold     float64 = 208.7
	instar3To4Threshold     float64 = 410.5
	instar4ToAdultThreshold float64 = 620
	// TODO: adult to die
	//lowerTempThreshold float64 =

	// survival rate: 1: 0.6488, 2: 0.9087, 3: 0.8948, 4: 0.822
	sRI1 float64 = 0.6488
	sRI2 float64 = 0.9087
	sRI3 float64 = 0.8948
	sRI4 float64 = 0.822
	sRA  float64 = 0.5837

	earthRadius float64 = 6371 // km

	minLat float64 = 24.396308  // Southernmost point in the US
	maxLat float64 = 49.384358  // Northernmost point in the contiguous US
	minLon float64 = -125.00165 // Westernmost point in the contiguous US
	maxLon float64 = -66.93457  // Easternmost point in the contiguous US
)



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

	minLat float64 = 24.396308  // Southernmost point in the US
	maxLat float64 = 49.384358  // Northernmost point in the contiguous US
	minLon float64 = -125.00165 // Westernmost point in the contiguous US
	maxLon float64 = -66.93457  // Easternmost point in the contiguous US

func width() float64 {
	return maxLon - minLon	
}

func InitializeCountry(maxLon, minLon float64, numOfFlies, numOfTrees) Country {
	country := Country{} // Create the initial country.
	width := maxLon - minLon
	country.width = width
	country.flies = make([]Fly, numOfFlies)
	country.trees = make([]Tree, numOfTrees)
	country.population = numOfFlies

	// Initialize the flies.
	for i := range country.flies {
		country.flies[i].position.x = // Read from file
		country.flies[i].position.y = // Read from file
		country.flies[i].velocity.x = rand.Float64() * 2
		country.flies[i].velocity.y = rand.Float64() * 5
		country.flies[i].acceleration.x = rand.Float64() * rand.Float64() * 2
		country.flies[i].acceleration.y = rand.Float64() * rand.Float64() * 5
		country.flies[i].stage = 0

		country.flies[i].isAlive = true
		country.flies[i].locationID = 0
	}

	for i: = range country.

	// Intialize the trees reading from data file
	for i := range country.trees {
		country.trees[i].position.x = 
		country.trees[i].position.y = 
	}

	return country
}


// InitializeQuadrants creates a 5x5 grid of Quadrants
func InitializeQuadrants() Quadrants {
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
	quadrants[0].temp += weatherData["OR"].x
	quadrants[1].temp += weatherData["CT"].x
	quadrants[2].temp += weatherData["VT"].x
	quadrants[3].temp += weatherData["MA"].x
	quadrants[4].temp += weatherData["ME"].x

	// Second Row
	quadrants[5].temp += weatherData["NJ"].x
	quadrants[6].temp += weatherData["NY"].x
	quadrants[7].temp += weatherData["PA"].x
	quadrants[8].temp += weatherData["MD"].x
	quadrants[9].temp += weatherData["DE"].x

	// Third Row:
	quadrants[10].temp += weatherData["WV"].x
	quadrants[11].temp += weatherData["VA"].x
	quadrants[12].temp += weatherData["NC"].x
	quadrants[13].temp += weatherData["DC"].x
	quadrants[14].temp += weatherData["SC"].x

	//Fourth Row:
	quadrants[15].temp += weatherData["NM"].x
	quadrants[16].temp += weatherData["MO"].x
	quadrants[17].temp += weatherData["IN"].x
	quadrants[18].temp += weatherData["OH"].x
	quadrants[19].temp += weatherData["MI"].x

	// Southwest to Southeast (Bottom Row)
	quadrants[20].temp += weatherData["AZ"].x
	quadrants[21].temp += weatherData["UT"].x
	quadrants[22].temp += weatherData["KY"].x
	quadrants[23].temp += weatherData["RI"].x
	quadrants[24].temp += weatherData["KS"].x

	return Quadrants{
		x:         minLon,
		y:         minLat,
		Quadrants: quadrants,
	}
}




// func InitializeCountry(width float64, population int) Country {

// 	country := Country{} // Create the initial country.
// 	country.width = width
// 	country.flies = make([]Fly, population)
// 	country.population = population

// 	// Read Weather data and store in a map
// 	MayJuneWeather := make(map[string]OrderedPair)

// 	// Load from subfolder of Data
// 	MayJuneWeather = LoadWeatherData("Data/Hatch_May-Jun")

// 	// Initialize the flies.
// 	for i := range country.flies {
// 		country.flies[i].position = LoadLatternFlyPosition("latternfly.csv")

// 		// Velocity and acceleration are random since no data is available
// 		country.flies[i].velocity.x = rand.Float64() * 2
// 		country.flies[i].velocity.y = rand.Float64() * 5
// 		country.flies[i].acceleration.x = rand.Float64() * rand.Float64() * 2
// 		country.flies[i].acceleration.y = rand.Float64() * rand.Float64() * 5

// 		//latternfly's stage consider as 0 : egg state
// 		country.flies[i].stage = 0

// 		// Set Energy to every fly
// 		// Latternfly's energy is random from its state and weather data from May to June
// 		// random from minTemp of the state to maxTemp of the state
// 		// States are sorted by alphabetical order
// 		// AZ, CT, DC, DE, IN, KS, KY, MA, MD, ME, MI, MO, NC, NJ, NM, NY, OH, OR, PA, RI, SC, UT, VA, VT, WV

// 		if country.flies[i].state == "AZ" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["AZ"].Y, MayJuneWeather["AZ"].X)
// 		} else if country.flies[i].state == "CT" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["CT"].Y, MayJuneWeather["CT"].X)
// 		} else if country.flies[i].state == "DC" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["DC"].Y, MayJuneWeather["DC"].X)
// 		} else if country.flies[i].state == "DE" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["DE"].Y, MayJuneWeather["DE"].X)
// 		} else if country.flies[i].state == "IN" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["IN"].Y, MayJuneWeather["IN"].X)
// 		} else if country.flies[i].state == "KS" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["KS"].Y, MayJuneWeather["KS"].X)
// 		} else if country.flies[i].state == "KY" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["KY"].Y, MayJuneWeather["KY"].X)
// 		} else if country.flies[i].state == "MA" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["MA"].Y, MayJuneWeather["MA"].X)
// 		} else if country.flies[i].state == "MD" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["MD"].Y, MayJuneWeather["MD"].X)
// 		} else if country.flies[i].state == "ME" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["ME"].Y, MayJuneWeather["ME"].X)
// 		} else if country.flies[i].state == "MI" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["MI"].Y, MayJuneWeather["MI"].X)
// 		} else if country.flies[i].state == "MO" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["MO"].Y, MayJuneWeather["MO"].X)
// 		} else if country.flies[i].state == "NC" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["NC"].Y, MayJuneWeather["NC"].X)
// 		} else if country.flies[i].state == "NJ" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["NJ"].Y, MayJuneWeather["NJ"].X)
// 		} else if country.flies[i].state == "NM" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["NM"].Y, MayJuneWeather["NM"].X)
// 		} else if country.flies[i].state == "NY" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["NY"].Y, MayJuneWeather["NY"].X)
// 		} else if country.flies[i].state == "OH" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["OH"].Y, MayJuneWeather["OH"].X)
// 		} else if country.flies[i].state == "OR" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["OR"].Y, MayJuneWeather["OR"].X)
// 		} else if country.flies[i].state == "PA" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["PA"].Y, MayJuneWeather["PA"].X)
// 		} else if country.flies[i].state == "RI" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["RI"].Y, MayJuneWeather["RI"].X)
// 		} else if country.flies[i].state == "SC" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["SC"].Y, MayJuneWeather["SC"].X)
// 		} else if country.flies[i].state == "UT" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["UT"].Y, MayJuneWeather["UT"].X)
// 		} else if country.flies[i].state == "VA" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["VA"].Y, MayJuneWeather["VA"].X)
// 		} else if country.flies[i].state == "VT" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["VT"].Y, MayJuneWeather["VT"].X)
// 		} else if country.flies[i].state == "WV" {
// 			country.flies[i].energy += randomInRange(MayJuneWeather["WV"].Y, MayJuneWeather["WV"].X)
// 		}
// 	}
// 	//randomly choose 10% of the flies to be alive
// 	for i := 0; i < population/10; i++ {
// 		country.flies[rand.Intn(population)].isAlive = true
// 	}
// 	return country
// }

// func randomInRange(min, max float64) float64 {
// 	rand.Seed(time.Now().UnixNano()) // Initialize the random number generator
// 	return min + rand.Float64()*(max-min)
// }


