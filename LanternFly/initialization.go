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

	weatherData, err := LoadWeatherData("data")
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

	return Quadrants{
		x:         minLon,
		y:         minLat,
		Quadrants: quadrants,
	}
}

func InitializeCountry(numberOfFlies, numberOfTree int) Country {
	var country Country
	country.width = maxLon - minLon
	country.flies = make([]Fly, numberOfFlies)
	country.trees = make([]Tree, numberOfTree)

	// Initialize flies
	for i := 0; i < numberOfFlies; i++ {
		country.flies[i].position.x = 
		country.flies[i].position.y = 
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

	// Initialize trees	
	for i := 0; i < numberOfTree; i++ {
		country.trees[i].position.x = 
		country.trees[i].position.y = 
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
