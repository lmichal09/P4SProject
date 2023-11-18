//List functions needed from Lantern Fly Simulation
//Author: Leila Michal, Emma Bouchard, Tiffany Ku, and Thu Pham

//NOTE: ctrl + f "TODO" to find all the things that need to be done
// viability and fecundity are related to temperature

package main

import (
	"math"
)

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

<<<<<<< HEAD
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

// UpdateLifeStage
func UpdateLifeStage(f []Fly, time float64) []Fly {
	for i := 0; i < len(flies); i++ {
		// Accumulate degree-days based on daily temperature
		degreeDays := ComputeDD(dailyTemp) //TODO: datatype of dailyTemp

		// Below the lower temperature threshold, the insect does not accumulate degree-days
		if degreeDays < lowerTempThreshold {
			flies[i].DegreeDays += 0
		} else {
			flies[i].DegreeDays += degreeDays
		}

		// Check and update life stage based on thermal thresholds
		switch flies[i].Stage {
		case 0: // Egg
			if flies[i].DegreeDays >= EggToLarvaThreshold {
				flies[i].Stage = 1      // Change to larva
				flies[i].DegreeDays = 0 // Reset degree days
			}
		case 1: // Larva
			if flies[i].DegreeDays >= LarvaToPupaThreshold {
				flies[i].Stage = 2      // Change to pupa
				flies[i].DegreeDays = 0 // Reset degree days
			}
		case 2: // Pupa
			if flies[i].DegreeDays >= PupaToAdultThreshold {
				flies[i].Stage = 3      // Change to adult
				flies[i].DegreeDays = 0 // Reset degree days
			}
			// No case for adults as they do not transform further
		}
	}
	return flies
}

// The energy (DD) accumulated during 1 day is calculated by subtracting the value of the lower temperature threshold from the value of the daily temperature
func ComputeDD(dailyTemperature float64) float64 {
	return dailyTemperature - lowerTempThreshold
}

// ComputeMortality updates the mortality status of flies.
func ComputeMortality(flies []Fly) []Fly {
    for i := range flies {
        if flies[i].Stage == 3 { // Check for adults
            // For adults, mortality is based on energy threshold
            flies[i].Energy += // TODO: Increment energy based on some logic

            if flies[i].Energy >= AdultEnergyThreshold {  
                flies[i].IsAlive = false
            }
        } else {
            // For immature insects, mortality can be based on a daily probability
            mortalityProbability := // TODO: Calculate based on temperature or other factors
            if someRandomCondition(mortalityProbability) { // TODO: someRandomCondition to simulate probability
                flies[i].IsAlive = false
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
        if flies[i].Stage == 3 { // only adults can lay eggs
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
        if flies[i].Stage == 3 { // adult
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




=======
>>>>>>> parent of d911a0a (main)
// Simulate Predator-Prey interaction
// Update population sizes based on the consumption rates and predation rules
// Track the population of lantern flies and predators over time
// Monitor population growth and decline based on the predation or other factors
func PredatorPreyBehavior(size int) int {

	return size

}

/*
InitializeHabitat()
CreateGrid() // make grid representing habitat
GetWeatherData
LatternflyParameters
InitializePopulation

InitializePreyPredatorModel()
PredatorPopulation
Set Parameter of Prey and Predator (eg: Growth Rate (a), Death Rate(b), Prey caught/ Predator/ unit time)



UpdatePreyPopulation
	Slices of timePoints
timePoints[0] - InitializePreyPredatorModel()
ComputePreyPopulation()

UpdatePredatorPopulation ()
Slices of timePoints
timePoints[0] - InitializePreyPredatorModel()
ComputePredatorPopulation()

UpdateHabitat()
UpdateLifeStage() // handle their development
ComputeMortality()
ComputeFecundity()
ComputeAdultMovement()
ComputeLarvalMovement()
KillSeniors() // should change the name here :)))
UpdateGridOccupancy()

UpdateLifeStage()
	for each fly
		Accumulate degree-day
		If thermal threshold met
			Advance to next stage

ComputeMotality()

ComputeViability()
	Get eggs laid based on temperature
	Add eggs to grids

ComputeFecundity() bool {
	fly.energy (t)
	If fly.energy == setFecundityEnergy {
		Return true
}
}

ComputeAdultMovement()
S = CalculateDistance(Adult, Tree)
ProbabilityOfMoving


ComputeLarvalMovement()


KillSeniors() {
	If senior.CompuLifeSpan {
		Population --
}
}

ComputeLifeSpan () Bool {
	If fly.energy >= setEnergy && fly.time >= setLifeSpan {
		Return true }
}
}

ComputePreyPopulation (current prey population, current predator population, Prey growth rate,  Prey death rate)
Prey Population += ( growth rate * current prey population  - death rate * current prey population * current predator population)


ComputePredatorPopulation (current prey population, current predator population, Predator growth rate, Predator death rate)
Predator Population += (- death rate*current predator population + growth rate * current prey population * current predator population)
*/

//Inpute a Country and Time
//Returns a new Country onject corrpespoding to updating the force of gravity on the objects in a given Country, with a tiem interveral in secs

func UpdateHabitat(currentCountry Country, time float64) Country {
	newCountry := CopyCountry(currentCountry) // Copy for all flies and attributes associated with each fly

	// range over all flies in universe and update accel, vel, and position
	for i := range newCountry.flies {
		newCountry.flies[i].acceleration = UpdateAccel(currentCountry, newCountry.flies[i])
		newCountry.flies[i].velocity = UpdateVelocity(newCountry.flies[i], time)
		newCountry.flies[i].position = UpdatePosition(newCountry.flies[i], time)

	}
	return newCountry
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

	force := ComputeNetForce(currentCountry, f)
	//Now compute accel

	//Split acceleration based on force
	accel.x = force.x
	accel.y = force.y

	return accel
}

// ComputeNetForce thake Country object and Fly b
// Return a new force (due to gravity) acting on b by all other objecys in given universe
func ComputeNetForce(currentCountry Country, f Fly) OrderedPair {
	var NetForce OrderedPair

	//range over all flies other than f and pass
	// computing the force of gravity to subroutine and then add components to net force
	for i := range currentCountry.flies {
		// only compute force if current Fly is not f
		if currentCountry.flies[i] != f {
			force := ComputeForce(f, currentCountry.flies[i])
			// add componets of force to NetForce
			NetForce.x += force.x
			NetForce.y += force.y
		}

	}
	return NetForce
}

//Takes teo Fly objects
//returns the orderedpair corresponding to the compnents of a force vector to the force of gravity of f2 acting on fl

func ComputeForce(fl, f2 Fly) OrderedPair {
	var force OrderedPair

	// apply formula
	// F= G *b.mass*f2.mass/(d*d)

	// Compute magnitude
	d := Distance(fl.position, f2.position)
	F := G * fl.mass * f2.mass / (d * d)

	// Then split into components
	dx := f2.position.x - fl.position.x // f2 is pulling on fl position
	dy := f2.position.y - fl.position.y
	force.x = F * (dx / d)
	force.y = F * (dy / d)
	return force
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
