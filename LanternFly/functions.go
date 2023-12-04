//List functions needed from Lantern Fly Simulation
//Author: Leila Michal, Emma Bouchard, Tiffany Ku, and Thu Pham

//NOTE: "TODO": all the things that need to be done

package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

const G = 6.67408e-11

// SimulateMigration
func SimulateMigration(initialCountry Country, numYears int) []Country {
	timePoints := make([]Country, numYears+1)
	timePoints[0] = initialCountry

	//range over num of generations and set the i-th country equal to updating the (i-1)th Country
	for i := 1; i < len(timePoints); i++ {
		timePoints[i] = UpdateCountry(timePoints[i-1])
		//TODO: PopulationSize()
	}

	finaltimePoints := SorttheTimePoints(timePoints)

	return finaltimePoints
}

// Population calculates and returns the population of flies in each state.
func Population(country *Country) int {
	totalFlies := 0

	for i := 0; i < len(country.flies); i++ {
		if country.flies[i].isAlive {
			totalFlies++
		}
	}

	return totalFlies
}

func SorttheTimePoints(timepoints []Country) []Country {
	n := len(timepoints)

	for i := 0; i <= n; i++ {
		timepoints[i] = SorttheCountry(timepoints[i])
	}
	return timepoints
}

func SorttheCountry(country Country) Country {
	n := len(country.flies)

	for i := 0; i <= n; i++ {
		country.flies[i] = SorttheFlies(country.flies[i])

	}
	return country
}

func SorttheFlies(fly Fly) Fly {

	if fly.stage == 0 {
		fly.color.red = 237
		fly.color.blue = 7
		fly.color.green = 15
	} else if fly.stage == 1 {
		fly.color.red = 252
		fly.color.blue = 20
		fly.color.green = 144
	} else if fly.stage == 2 {
		fly.color.red = 252
		fly.color.blue = 20
		fly.color.green = 249
	} else if fly.stage == 3 {
		fly.color.red = 43
		fly.color.blue = 20
		fly.color.green = 252
	} else if fly.stage == 4 {
		fly.color.red = 20
		fly.color.blue = 252
		fly.color.green = 214
	} else if fly.stage == 5 {
		fly.color.red = 175
		fly.color.blue = 252
		fly.color.green = 20
	} else if fly.stage == 6 {
		fly.color.red = 0
		fly.color.blue = 0
		fly.color.green = 0
	}
	return fly
}

// lanternfly only has one generation per year
// in each generation, the flies go through 5 stages
// by the end of the year, there should only be eggs
// need to record the number of adults
// also keep track on the position of the eggs, which will then be used for next year's simulation
func UpdateCountry(currCountry Country) Country {
	newCountry := CopyCountry(currCountry) // Copy for all flies and attributes associated with each fly

	var totalNewEggs []Fly

	// weather
	var weather []Quadrant

	// loop through days
	for i := 0; ; i++ {
		// keep looping until all flies are dead, except for eggs
		if CheckDead(currCountry.flies) {
			break
		}

		// loop through flies
		for j := 0; j < len(currCountry.flies); j++ {
			// compute degree days
			newCountry.flies[j].energy += ComputeDegreeDay(&currCountry.flies[j], weather)

			// update life stage
			newCountry.flies[j].stage = UpdateLifeStage(&newCountry.flies[j])

			// compute mortality
			newCountry.flies[j].isAlive = ComputeMortality(&newCountry.flies[j])

			// compute movement
			newCountry.flies[j].position = ComputeMovement(&newCountry.flies[j])

			// lay eggs
			if newCountry.flies[j].stage == 5 {
				newEggs := ComputeFecundity(newCountry.flies[j])
				totalNewEggs = append(totalNewEggs, newEggs...)
			}
		}
	}

	// remove all dead flies
	for i := 0; i < len(newCountry.flies); i++ {
		if !newCountry.flies[i].isAlive {
			newCountry.flies = append(newCountry.flies[:i], newCountry.flies[i+1:]...)
		}
	}

	if len(totalNewEggs) > 0 {
		panic("something's wrong")
	}

	// add new eggs to the country
	newCountry.flies = append(newCountry.flies, totalNewEggs...)

	return newCountry
}

// ComputeDegreeDay calculates the degree days for a single day.
func ComputeDegreeDay(fly *Fly, quadrants []Quadrant) float64 {
	// get the quadrant of the fly to determine the temperature
	quadrantID := GetQuadrant(fly, quadrants)

	// get the temperature of the quadrant
	temperature := GetTemperature(quadrantID, quadrants) // TODO: this is the max temp?

	// get the base temperature base on the fly's stage
	baseTemp := GetBaseTemp(fly.stage)

	// calculate the degree days
	degreeDays := (temperature+baseTemp)/2 - baseTemp

	// if degreeDays is negative, set it to 0
	if degreeDays < 0 {
		degreeDays = 0
	}

	return degreeDays
}

// GetQuadrant returns the quadrant of the fly.
// Base temperature used for calculating GDD for nymph. 1: 13.00°C, 2: 12.43°C, 3: 8.48°C, 4: 6.29°C
func GetBaseTemp(stage int) float64 {
	temp := 0.0

	if stage == 1 {
		temp = 13.00
	} else if stage == 2 {
		temp = 12.43
	} else if stage == 3 {
		temp = 8.48
	} else if stage == 4 {
		temp = 6.29
	}

	return temp
}
func ReadTrees() []OrderedPair {
	var habitats []OrderedPair
	for i, record := range records {
		if i == 0 { // Skip header
			continue
		}

		// Parse longitude
		longitude, err := strconv.ParseFloat(record[0], 64) // Assuming longitude is in the first column
		if err != nil {
			fmt.Printf("Error parsing longitude in row %d: %v\n", i+1, err)
			continue
		}

		// Parse latitude
		latitude, err := strconv.ParseFloat(record[1], 64) // Assuming latitude is in the second column
		if err != nil {
			fmt.Printf("Error parsing latitude in row %d: %v\n", i+1, err)
			continue
		}

		// Append the habitat to the slice
		habitats = append(habitats, OrderedPair{x: longitude, y: latitude})
	}
	return habitats
}

// UpdateLifeStage() updates the life stage of flies based on the cumulative degree-days (CDD)
func UpdateLifeStage(fly *Fly) int {
	stage := fly.stage

	if fly.stage < 0 || fly.stage > 5 {
		panic("invalid life stage")
	}

	// TODO: hatch
	// TODO: must survive according to a factor of egg viability & must accumulate a certain number of degree-days to molt to the next stage

	// Required GDD for each stage. 1: 166.6, 208.7, 410.5, 620
	if fly.energy >= 0 && fly.energy < 166.6 {
		stage = 1
	} else if fly.energy >= instar1To2Threshold && fly.energy < instar2To3Threshold {
		stage = 2
	} else if fly.energy >= instar2To3Threshold && fly.energy < instar3To4Threshold {
		stage = 3
	} else if fly.energy >= instar3To4Threshold && fly.energy < instar4ToAdultThreshold {
		stage = 4
	} else if fly.energy >= instar4ToAdultThreshold {
		stage = 5
	} else { // TODO: adult to die
		stage = 6 // dead
	}

	return stage
}

// ComputeMortality updates the mortality status of flies.
// survival rate: 1: 0.6488, 2: 0.9087, 3: 0.8948, 4: 0.822
func ComputeMortality(fly *Fly) bool {
	if fly.stage == 1 {
		if rand.Float64() > sRI1 {
			fly.isAlive = false
		}
	} else if fly.stage == 2 {
		if rand.Float64() > sRI2 {
			fly.isAlive = false
		}
	} else if fly.stage == 3 {
		if rand.Float64() > sRI3 {
			fly.isAlive = false
		}
	} else if fly.stage == 4 {
		if rand.Float64() > sRI4 {
			fly.isAlive = false
		}
	} else if fly.stage == 5 {
		if rand.Float64() > sRA {
			fly.isAlive = false
		}
	}

	return fly.isAlive
}

// ComputeFecundity
// females lay one or two egg masses, each containing 30-60 eggs
func ComputeFecundity(fly Fly) []Fly {
	newFly := make([]Fly, 0)

	// TODO: probability of laying eggs 10.5%

	// randomly choose the number of egg masses
	numEggMasses := rand.Intn(2) + 1

	// randomly choose the number of eggs in each egg mass
	numEggs := rand.Intn(30) + 30

	totalEggs := numEggMasses * numEggs

	// location of the eggs is the location of the adult
	for i := 0; i < totalEggs; i++ {
		newEgg := Fly{
			position: OrderedPair{
				x: fly.position.x,
				y: fly.position.y,
			},
			stage:   0,
			energy:  0,
			isAlive: true,
		}
		newFly = append(newFly, newEgg)
	}

	return newFly
}

// ComputeMovement updates the position of adult flies
func ComputeMovement(fly *Fly) OrderedPair {
	// TODO: determine the proportion of random vs. directed movement
	if rand.Float64() < 0.5 {
		// random movement
		return RandomMovement(fly)
	} else {
		// directed movement
		return DirectedMovement(fly)
	}
}

// RandomMovement updates the position of adult flies based on random movement
func RandomMovement(fly *Fly) OrderedPair {
	maxDistance := 0.0 // TODO: max distance

	// possibility of flies carried by human
	if rand.Float64() < 0.1 {
		return LongDistanceMovement(fly)
	}

	// randomly choose a distance
	distance := rand.Float64() * maxDistance

	// randomly choose a direction
	angle := rand.Float64() * 2 * math.Pi

	// calculate the new position
	new := ConvertDistanceToCoordinates(distance, angle, fly.position)

	return new
}

// LongDistanceMovement simulates long-distance movement for a Fly.
func LongDistanceMovement(fly *Fly) OrderedPair {
	maxDistance := 2000.0 // TODO: Define your max long-distance here (in kilometers)

	// Randomly choose a distance within the maximum limit
	distance := rand.Float64() * maxDistance

	// Randomly choose a direction
	angle := rand.Float64() * 2 * math.Pi

	// Calculate the new position using the conversion function
	return ConvertDistanceToCoordinates(distance, angle, fly.position)
}

// DirectedMovement updates the position of adult flies based on directed movement
func DirectedMovement(fly *Fly) OrderedPair {
	// identify the nearest host tree or a direction with higher concentration of host trees
	direction := FindHostDirection(fly.position, habitats)

	// move towards the direction, assume a simpler linear movement
	new := ConvertDistanceToCoordinates(1, direction, fly.position)

	return new
}

// ConvertDistanceToCoordinates converts a given distance (in kilometers) and direction (in radians) into new coordinates.
func ConvertDistanceToCoordinates(distance, direction float64, startingCoordinates OrderedPair) OrderedPair {
	newX := startingCoordinates.x + (distance/earthRadius)*(180.0/math.Pi)*math.Cos(direction)
	newY := startingCoordinates.y + (distance/earthRadius)*(180.0/math.Pi)*math.Sin(direction)
	return OrderedPair{newX, newY}
}

func FindHostDirection(position OrderedPair, hostMaps []Coordinate) OrderedPair {
	minDistance := math.MaxFloat64
	var nearestTree Coordinate

	// Find the nearest host tree
	for _, tree := range hostMaps {
		// calculate distance between position of fly and tree (longtitude, latitude) using Haversine formula
		distance := Haversine(position, tree)
		if distance < minDistance {
			minDistance = distance
			nearestTree = tree
		}
	}

	// TODO: how far can a fly move in one day?

	// Find the direction of the nearest host tree
	direction := OrderedPair{
		x: nearestTree.longitude - position.x,
		y: nearestTree.latitude - position.y,
	}

	return direction
}

// Haversine calculates the distance between two points on a sphere
// on the Earth given their longitude and latitude in degrees.
func Haversine(position1, position2 OrderedPair) float64 {
	// Convert latitude and longitude from degrees to radians.
	lon1Rad := position1.x * math.Pi / 180
	lat1Rad := position1.y * math.Pi / 180
	lon2Rad := position2.x * math.Pi / 180
	lat2Rad := position2.y * math.Pi / 180

	// Calculate the differences in latitude and longitude.
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	// Apply the Haversine formula.
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Calculate the distance.
	distance := earthRadius * c

	return distance
}

// CheckDead takes a slice of Fly and return true if all flies are dead
func CheckDead(flies []Fly) bool {
	for i := range flies {
		if flies[i].isAlive {
			return false
		}
	}
	return true
}

// DivideCountry divides the country into 25 sections
func DivideCountry(country Country) []Quadrant {
	const (
		gridRows    int     = 5
		gridColumns int     = 5
		totalWidth  float64 = maxLon - minLon
		totalHeight float64 = maxLat - minLat
		quadWidth   float64 = totalWidth / float64(gridColumns)
		quadHeight  float64 = totalHeight / float64(gridRows)
	)

	// Check for invalid country dimensions
	if country.width <= 0 || country.width <= 0 { //
		panic("invalid country dimensions")
	}

	qts := make([]Quadrant, 25)

	// Id for each quadrant
	id := 1

	for row := 0; row < gridRows; row++ {
		for col := 0; col < gridColumns; col++ {
			q := Quadrant{
				id:    id,
				x:     minLon + (float64(col) * quadWidth),
				y:     minLat + (float64(row) * quadHeight),
				width: quadWidth,
				// temperature will be set later based on weather data
			}
			qts = append(qts, q)
			id++
		}
	}

	return qts
}

// GetQuadrant returns the quadrant of the fly.
// divide the country into 25 sections
func GetQuadrant(fly *Fly, quadrants []Quadrant) int {
	var quadrant int

	// loop through the quadrant slice and find the quadrant with the matching id
	for _, q := range quadrants {
		if fly.position.x >= q.x && fly.position.x <= q.x+q.width && fly.position.y >= q.y && fly.position.y <= q.y+q.width {
			quadrant = q.id
		}
	}

	// if the fly is out of bounds, panic
	// TODO: panic("fly is out of simulation bounds")

	return quadrant // Placeholder
}

// GetTemperature returns the temperature of the quadrant.
func GetTemperature(quadrantID int, quadrant []Quadrant) float64 {
	temp := 0.0
	// loop through the quadrant slice and find the quadrant with the matching id
	for _, q := range quadrant {
		if q.id == quadrantID {
			temp = q.temp
		}
	}

	if quadrantID < 1 || quadrantID > len(quadrant) {
		panic("invalid quadrant ID")
	}

	return temp
}

// CopyCountry takes a Country and return a copy of all flies in this Country with fields copied over.
func CopyCountry(currentCountry Country) Country {
	var newCountry Country
	newCountry.width = currentCountry.width

	// copy flies over
	numFlies := len(currentCountry.flies)
	newCountry.flies = make([]Fly, numFlies)

	//copy every fly's field in the new Country
	for i := range newCountry.flies {

		newCountry.flies[i] = CopyFly(currentCountry.flies[i])

	}
	return newCountry

}

// CopyFly takes Fly object and return an a Fly with all field of input object
func CopyFly(oldFly Fly) Fly {
	var newFly Fly

	//copy ordered pair
	newFly.position.x = oldFly.position.x
	newFly.position.y = oldFly.position.y
	newFly.velocity.y = oldFly.velocity.y
	newFly.acceleration.y = oldFly.acceleration.y
	newFly.stage = oldFly.stage

	return newFly

}
