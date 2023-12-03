//List functions needed from Lantern Fly Simulation
//Author: Leila Michal, Emma Bouchard, Tiffany Ku, and Thu Pham

//NOTE: ctrl + f "TODO" to find all the things that need to be done
// viability and fecundity are related to temperature

package main

import (
	"math"
	"math/rand"
)

const G = 6.67408e-11

// SimulateMigration
func SimulateMigration(initialCountry Country, numYears int) []Country {
	timePoints := make([]Country, numYears+1)
	timePoints[0] = initialCountry

	//range over num of generations and set the i-th country equal to updating the (i-1)th Country
	for i := 1; i <= len(timePoints); i++ {
		timePoints[i] = UpdateCountry(timePoints[i-1])
		//PopulationSize()
	}
	return timePoints
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

// lanternfly only has one generation per year
// in each generation, the flies go through 5 stages
// by the end of the year, there should only be eggs
// need to record the number of adults
// also keep track on the position of the eggs, which will then be used for next year's simulation
func UpdateCountry(currCountry Country) Country {
	newCountry := CopyCountry(currCountry) // Copy for all flies and attributes associated with each fly

	// loop through days
	for i := 0; ; i++ {
		// keep looping until all flies are dead, except for eggs
		if CheckDead(currCountry.flies) {
			break
		}

		// loop through flies
		for j := 0; j < len(currCountry.flies); j++ {
			// compute degree days
			newCountry.flies[j].energy += ComputeDegreeDay(&currCountry.flies[j])

			// update life stage
			newCountry.flies[j].stage = UpdateLifeStage(&newCountry.flies[j])

			// compute mortality
			newCountry.flies[j].isAlive = ComputeMortality(&newCountry.flies[j])

			// compute movement
			newCountry.flies[j].position = ComputeMovement(&newCountry.flies[j])

			// lay eggs
			if newCountry.flies[j].stage == 5 {
				newEggs := ComputeFecundity(newCountry.flies[j])
				newCountry.flies = append(newCountry.flies, newEggs...)
			}
		}
	}

	// remove all dead flies
	for i := 0; i < len(newCountry.flies); i++ {
		if !newCountry.flies[i].isAlive {
			newCountry.flies = append(newCountry.flies[:i], newCountry.flies[i+1:]...)
		}
	}

	return newCountry
}

// ComputeDegreeDay calculates the degree days for a single day.
func ComputeDegreeDay(fly *Fly) float64 {
	// get the quadrant of the fly to determine the temperature
	quadrant := GetQuadrant(fly)

	// get the temperature of the quadrant
	temperature := GetTemperature(quadrant) // TODO: this is the max temp?

	// get the base temperature base on the fly's stage
	baseTemp := GetBaseTemp(fly.stage)

	// calculate the degree days
	degreeDays := (temperature-baseTemp)/2 - baseTemp

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

// UpdateLifeStage() updates the life stage of flies based on the cumulative degree-days (CDD)
func UpdateLifeStage(fly *Fly) int {
	stage := fly.stage

	// TODO: hatch
	// TODO: must survive according to a factor of egg viability & must accumulate a certain number of degree-days to molt to the next stage

	// Required GDD for each stage. 1: 166.6, 208.7, 410.5, 620
	if fly.energy >= 0 && fly.energy < 166.6 {
		stage = 1
	} else if fly.energy >= instar1To2Threshold {
		stage = 2
	} else if fly.energy >= instar2To3Threshold {
		stage = 3
	} else if fly.energy >= instar3To4Threshold {
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
func ComputeFecundity(fly Fly, temperature float64) []Fly {
	newFly := make([]Fly, 0)

	// TODO: probability of laying eggs

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
	// determine the proportion of random vs. directed movement
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

	// randomly choose a distance
	distance := rand.Float64() * maxDistance

	// randomly choose a direction
	angle := rand.Float64() * 2 * math.Pi

	// calculate the new position
	newX := fly.position.x + distance*math.Cos(angle)
	newY := fly.position.y + distance*math.Sin(angle)

	return OrderedPair{newX, newY}
}

// DirectedMovement updates the position of adult flies based on directed movement
func DirectedMovement(fly *Fly) OrderedPair {
	// identify the nearest host tree or a direction with higher concentration of host trees
	direction := FindHostDirection(fly.position, hostMaps)

	// move towards the direction, assume a simpler linear movement
	newX := fly.position.x + direction.x
	newY := fly.position.y + direction.y

	return OrderedPair{newX, newY}
}

func FindHostDirection(position OrderedPair, hostMaps []HostMap) OrderedPair {
	minDistance := math.MaxFloat64
	var nearestTree HostMap

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
	lon1Rad := position1 * math.Pi / 180
	lat1Rad := position1 * math.Pi / 180
	lon2Rad := position2 * math.Pi / 180
	lat2Rad := position2 * math.Pi / 180

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

// UpdateAccel takes a universe object and a Fly int hat universe.
// Returns the net acceleration due to the force of gravity of the Fly (in componets) computed overall flies in the universe.
func UpdateAccel(currentCountry Country, f Fly) OrderedPair {
	var accel OrderedPair
	// randomly

	//Split acceleration based on force
	accel.x = force.x
	accel.y = force.y

	return accel
}

//UpdateVelocity take a Fly and a float time
//Uses components in that Fly estimates over time secs

func UpdateVelocity(f Fly, time float64) OrderedPair {
	var vel OrderedPair

	vel.x = f.velocity.x + f.acceleration.x*time
	vel.y = f.velocity.y + f.acceleration.y*time

	return vel
}

// Updatepositon take a Fly and a float time
// Uses components in that Fly estimates over time secs
func UpdatePosition(f Fly, time float64) OrderedPair {
	var pos OrderedPair

	pos.x = f.position.x + f.velocity.x*time + 0.5*f.acceleration.x*time*time
	pos.y = f.position.y + f.velocity.y*time + 0.5*f.acceleration.y*time*time

	return pos
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

func GetQuadrant(fly Fly) int {
	var quadrant int

	return quadrant // Placeholder
}

func GetTemperature(quadrant int) float64 {
	return 0 // Placeholder
}
