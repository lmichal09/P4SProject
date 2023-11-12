//List functions needed from Lantern Fly Simulation
//Author: Leila Michal, Emma Bouchard, Tiffany Ku, and Thu Pham

package main

import (
	"math"
)

// SimulateMigration
// Growth simulation to show how they invade the map
// Simulate pattern of lantern flies in Pittsburgh
func SimulateMigration(initialCountry Country, numGens int, time float64) []Country {
	timePoints := make([]Country, numGens+1)
	timePoints[0] = initialCountry

	//range over num of generations and set the i-th country equal to updating the (i-1)th Country
	for i := 1; i <= len(timePoints); i++ {
		timePoints[i] = UpdateCountry(timePoints[i-1], time)
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

// Simulate Predator-Prey interaction
// Update population sizes based on the consumption rates and predation rules
func PredatorPreyBehavior(size int) int {

	return size

}

//Inpute a Country and Time
//Returns a new Country onject corrpespoding to updating the force of gravity on the objects in a given Country, with a tiem interveral in secs

func UpdateCountry(currentCountry Country, time float64) Country {
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
	accel.x = force.x / f.mass
	accel.y = force.y / f.mass

	return accel
}

// ComputeNetForce thake Country object and Fly b
// Return a new force (due to gravity) acting on b by all other objecys in given universe
func ComputeNetForce(currentCountry Country, f Fly) OrderedPair {
	var NetForce OrderedPair

	//range over all flies other than b and pass
	// computing the force of gravity to subroutine and the nadd components to net force
	for i := range currentCountry.flies {
		// only compute force if current Fly is not b
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

// takes a slice u of Country objects as input
// returns a slice of float64 variables having the same length as u.flies
func AverageSpeed(u []Fly) []float64 {
	//Get length of flies in universe
	numBodies := len(u)

	//Initalize average speed
	//Make array to store average speed of each Fly
	AvgSpeed := make([]float64, numBodies)

	//Iterate over each Fly
	for i := 0; i < numBodies; i++ {
		CombinedSpeed := Speed(u[i].velocity)

		// Calculate the average speed for i-th Fly
		AvgSpeed[i] = CombinedSpeed
	}

	//return slice
	return AvgSpeed
}

func Speed(velocity OrderedPair) float64 {
	return math.Sqrt(velocity.x*velocity.x + velocity.y*velocity.y)

}
