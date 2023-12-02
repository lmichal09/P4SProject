//List functions needed from Lantern Fly Simulation
//Author: Leila Michal, Emma Bouchard, Tiffany Ku, and Thu Pham

//NOTE: ctrl + f "TODO" to find all the things that need to be done
// viability and fecundity are related to temperature

package main

import (
	"math"
)

const G = 6.67408e-11

// SimulateMigration
// Growth simulation to show how they invade the map
// Simulate pattern of lantern flies in Pittsburgh
func SimulateMigration(initialHabitat Country, numGens int, time float64) []Country {
	timePoints := make([]Country, numGens+1)
	timePoints[0] = initialHabitat

	//range over num of generations and set the i-th country equal to updating the (i-1)th Country
	for i := 1; i <= len(timePoints); i++ {
		timePoints[i] = UpdateHabitat(timePoints[i-1], time)
		//PopulationSize()
	}
	return timePoints
}

// Track the population of lantern flies and predators over time
// Monitor population growth and decline based on the predation or other factors
func PopulationSize(size float64) float64 {
	var PopulationSize float64
	PopulationSize = 0.1 * size

	return PopulationSize
}

func UpdateHabitat(currHabitat Habitat, time float64) Habitat {
	newHabitat := CopyHabitat(currentHabitat) // Copy for all flies and attributes associated with each fly

	currHabitat.flies = UpdateLifeStage(currHabitat.flies, time)
	currHabitat.flies = ComputeMortality(currHabitat.flies)
	currHabitat.flies = ComputeFecundity(currHabitat.flies)
	currHabitat.flies = ComputeAdultMovement(currHabitat.flies)
	currHabitat.flies = ComputeLarvalMovement(currHabitat.flies)
	currHabitat.flies = RemoveSeniors(currHabitat.flies)

	return newHabitat
}


// UpdateLifeStage() updates the life stage of flies based on the cumulative degree-days (CDD)
func UpdateLifeStage(flies []Fly, cdd float64) []Fly {
    for i := range flies {

		// TODO: Assuming each fly has a location ID or some identifier to link it to weather data?
		locationID := flies[i].locationID
		tempData := GetWeatherData[locationID]

		// compute degree dat for this fly
		dd := ComputeDegreeDay(tempData.maxTemp, tempData.minTemp, baseTemp, upperTemp)

		// update cumulative degree-days for this fly
		flies[i].energy += dd

		// update life stage based on cumulative degree-days
        switch {
        case cdd < 270:
            // Egg stage
            flies[i].stage = 0

        case cdd >= 270 && cdd < 465:
            // First instar stage
            flies[i].stage = 1

        case cdd >= 465 && cdd < 645:
            // Second instar stage
            flies[i].stage = 2

        case cdd >= 645 && cdd < 825:
            // Third instar stage
            flies[i].stage = 3

        case cdd >= 825 && cdd < 1112:
            // Fourth instar stage
            flies[i].stage = 4

        case cdd >= 1112 && cdd < 1825:
            // Adult stage
            flies[i].stage = 5

        case cdd >= 1825:
            // After reaching the egg laying stage, the stage remains adult
            flies[i].stage = 5
        }
    }

    return flies
}


// ComputeDegreeDay calculates the degree days for a single day.
// maxTemp and minTemp are the day's maximum and minimum temperatures.
// baseTemp is the base (threshold) temperature for development.
// upperTemp is an optional upper threshold temperature; use a negative value if not needed.
func ComputeDegreeDay(maxTemp, minTemp, baseTemp, upperTemp float64) float64 {
    // Calculate the mean temperature for the day
    meanTemp := (maxTemp + minTemp) / 2

    // If mean temperature is below the base temperature, return 0
    if meanTemp < baseTemp {
        return 0
    }

    // If an upper threshold is specified and the mean temperature is above it, cap the mean temperature
    if upperTemp >= 0 && meanTemp > upperTemp {
        meanTemp = upperTemp
    }

    // Calculate and return the degree days
    return meanTemp - baseTemp
}


// ComputeMortality updates the mortality status of flies.
func ComputeMortality(flies []Fly) []Fly {
    for i := range flies {
        if flies[i].stage == 3 { // Check for adults
            // For adults, mortality is based on energy threshold
            flies[i].Energy += // TODO: Increment energy based on some logic

            if flies[i].Energy >= AdultEnergyThreshold {  
                flies[i].isAlive = false
            }
        } else {
            // For immature insects, mortality can be based on a daily probability
            mortalityProbability := // TODO: Calculate based on temperature or other factors
            if someRandomCondition(mortalityProbability) { // TODO: someRandomCondition to simulate probability
                flies[i].isAlive = false
            }
        }
    }
    return flies
}

// ComputeFecundity 
// NOTE: females lay one or two egg masses, each containing 30-60 eggs
// TODO: eggs are laid in the same grid as the adult or the neighboring grid?
func ComputeFecundity(flies []Fly, temperature float64) []Fly {
    for i := range flies {
        if flies[i].stage == 3 { // only adults can lay eggs
            // Calculate fecundity based on temperature
            // φ(T) - Oviposition probability function dependent on temperature
            flies[i].Fecundity = calculateFecundity(temperature)
        }
    }
    return flies
}

// TODO: the calculation is not in the paper, need to check (Garcia et al. 2016).
func CalculateFecundity(temperature float64) float64 {
    // Implement the calculation for fecundity based on temperature
    return 0 // Placeholder
}

// ComputeLarvalMovement updates the position of adult flies 
// TODO: in the paper we referenced, the movement of an adult inside the simulated area at each time step had no preferential direction and was calculated as described in Garcia et al. but i think we should consider the distribution of trees, or other factors
func ComputeAdultMovement(flies []Fly) []Fly {
    for i := range flies {
        if flies[i].stage == 3 { // adult
            // Calculate movement based on a probability function
			// TODO: MovingProbability() ???
            distance := CalculateMovementDistance()
            flies[i].Position = UpdatePosition(flies[i].Position, distance)
        }
    }
    return flies
}

func CalculateMovementDistance() float64 {
    // TODO: Implement logic to calculate movement distance
    return 0 // Placeholder
}

func UpdatePosition(position OrderedPair, distance float64) OrderedPair {
    // Update the position based on the calculated distance
    // This is a simplified placeholder logic
    return OrderedPair{position.X + distance, position.Y}
}

// Each larva had a probability l of moving to an adjacent plant per day, using a Moore neighborhood of radius 1
func ComputeLarvalMovement(flies []Fly, moveProbability float64) []Fly {
	for i := range flies {
		if flies[i].stage == /* larval stage identifier */ {
			if ShouldMove(moveProbability) {
				flies[i].Position = GetNewPosition(flies[i].Position)
			}
		}
	}
	return flies
}

func ShouldMove(moveProb float64) bool {
	return rand.Float64() < moveProb
}

func GetNewPosition(currentPosition OrderedPair) OrderedPair {
	// Moore neighborhood moves: stay, or move to one of the 8 surrounding cells


	return OrderedPair{x: currentPosition.x + dx, y: currentPosition.y + dy}
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

	newFly.PercentConsumed = oldFly.PercentConsumed
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

// Distance takes two position ordered pairs and it returns the distance between these two points in 2-D space.
func Distance(p1, p2 OrderedPair) float64 {
	// this is the distance formula from days of precalculus long ago ...
	deltaX := p1.x - p2.x
	deltaY := p1.y - p2.y
	return math.Sqrt(deltaX*deltaX + deltaY*deltaY)
}

func MaxDistance(c Country) float64 {
	//Initial where to store max distance
	maxDist := 0.0

	//Iterare though all the flies and calculate distance for all the flies in universe
	for i := 0; i < len(c.flies); i++ {
		for j := i; j < len(c.flies); j++ {
			dist := Distance(c.flies[i].position, c.flies[j].position)

			//Check if dist is the max distance between the two points
			if dist > maxDist {
				//make dist the new max distance
				maxDist = dist
			}

		}

	}

	return maxDist

}

// takes a slice f of Country objects as input
// returns a slice of float64 variables having the same length as f.flies
func AverageSpeed(f []Fly) []float64 {
	//Get length of flies in universe
	numFlies := len(f)

	//Initalize average speed
	//Make array to store average speed of each Fly
	AvgSpeed := make([]float64, numFlies)

	//Iterate over each Fly
	for i := 0; i < numFlies; i++ {
		CombinedSpeed := Speed(f[i].velocity)

		// Calculate the average speed for i-th Fly
		AvgSpeed[i] = CombinedSpeed
	}

	//return slice
	return AvgSpeed
}

func Speed(velocity OrderedPair) float64 {
	return math.Sqrt(velocity.x*velocity.x + velocity.y*velocity.y)

}
