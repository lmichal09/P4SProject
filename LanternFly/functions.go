//List functions needed from Lantern Fly Simulation
//Author: Leila Michal, Emma Bouchard, Tiffany Ku, and Thu Pham

package main

import (
	"math"
)

// SimulateMigration
// Growth simulation to show how they invade the map
// Simulate pattern of lantern flies in Pittsburgh
func SimulateMigration(lanternFly [][]int) {

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

//Inpute a uNIVERSE AND TIME
//Returns a new universe onject corrpespoding to updating the force of gravity on the objects in a given universe, with a tiem interveral in secs

func UpdateUniverse(currentUniverse Universe, time float64) Universe {
	newUniverse := CopyUniverse(currentUniverse) // Copy for all bodies and attributes associated with each body

	// range over all bodies in universe and update accel, vel, and position
	for i := range newUniverse.bodies {
		newUniverse.bodies[i].acceleration = UpdateAccel(currentUniverse, newUniverse.bodies[i])
		newUniverse.bodies[i].velocity = UpdateVelocity(newUniverse.bodies[i], time)
		newUniverse.bodies[i].position = UpdatePosition(newUniverse.bodies[i], time)

	}
	return newUniverse
}

// CopyUniverse  takes a universe and return a copy of all bodies in this universe with fields copied over.
func CopyUniverse(currentUniverse Universe) Universe {
	var newUniverse Universe
	newUniverse.width = currentUniverse.width

	// copy  bodies over
	numBodies := len(currentUniverse.bodies)
	newUniverse.bodies = make([]Body, numBodies)

	//copy every body's field in the new universe
	for i := range newUniverse.bodies {

		newUniverse.bodies[i] = CopyBody(currentUniverse.bodies[i])

	}
	return newUniverse

}

// CopybODY TAKE bODY object and return an a body with all field of input object
func CopyBody(oldBody Body) Body {
	var newBody Body

	newBody.name = oldBody.name
	newBody.mass = oldBody.mass
	newBody.red = oldBody.red
	newBody.green = oldBody.green
	newBody.blue = oldBody.blue

	//copy ordered pair
	newBody.position.x = oldBody.position.x
	newBody.position.y = oldBody.position.y
	newBody.velocity.y = oldBody.velocity.y
	newBody.acceleration.y = oldBody.acceleration.y

	return newBody

}

// UpdateAccel takes a universe object and a body int hat universe.
// Returns the net acceleration due to the force of gravity of the body (in componets) computed overall bodies in the universe.
func UpdateAccel(currentUniverse Universe, b Body) OrderedPair {
	var accel OrderedPair

	force := ComputeNetForce(currentUniverse, b)
	//Now compute accel

	//Split acceleration based on force
	accel.x = force.x / b.mass
	accel.y = force.y / b.mass

	return accel
}

// ComputeNetForce thake Universe object and body b
// Return a new force (due to gravity) acting on b by all other objecys in given universe
func ComputeNetForce(currentUniverse Universe, b Body) OrderedPair {
	var NetForce OrderedPair

	//range over all bodies other than b and pass
	// computing the force of gravity to subroutine and the nadd components to net force
	for i := range currentUniverse.bodies {
		// only compute force if current body is not b
		if currentUniverse.bodies[i] != b {
			force := ComputeForce(b, currentUniverse.bodies[i])
			// add componets of force to NetForce
			NetForce.x += force.x
			NetForce.y += force.y
		}

	}
	return NetForce
}

//Takes teo body objects
//returns the orderedpair corresponding to the compnents of a force vector to the force of gravity of b2 acting on b1

func ComputeForce(b1, b2 Body) OrderedPair {
	var force OrderedPair

	// apply formula
	// F= G *b.mass*b2.mass/(d*d)

	// Compute magnitude
	d := Distance(b1.position, b2.position)
	F := G * b1.mass * b2.mass / (d * d)

	// Then split into components
	dx := b2.position.x - b1.position.x // b2 is pulling on b1 position
	dy := b2.position.y - b1.position.y
	force.x = F * (dx / d)
	force.y = F * (dy / d)
	return force
}

//UpdateVelocity take a Body and a flost time/
//Uses components in that bODY ESTIAMTED OVER TIME SECONDS

func UpdateVelocity(b Body, time float64) OrderedPair {
	var vel OrderedPair

	vel.x = b.velocity.x + b.acceleration.x*time
	vel.y = b.velocity.y + b.acceleration.y*time

	return vel
}

// //Updatepositon take a Body and a float time
// Uses components in that bODY ESTIAMTED OVER TIME SECONDS
func UpdatePosition(b Body, time float64) OrderedPair {
	var pos OrderedPair

	pos.x = b.position.x + b.velocity.x*time + 0.5*b.acceleration.x*time*time
	pos.y = b.position.y + b.velocity.y*time + 0.5*b.acceleration.y*time*time

	return pos
}

// Distance takes two position ordered pairs and it returns the distance between these two points in 2-D space.
func Distance(p1, p2 OrderedPair) float64 {
	// this is the distance formula from days of precalculus long ago ...
	deltaX := p1.x - p2.x
	deltaY := p1.y - p2.y
	return math.Sqrt(deltaX*deltaX + deltaY*deltaY)
}

func MaxDistance(u Universe) float64 {
	//Initial where to store max distance
	maxDist := 0.0

	//Iterare though all the bodies and calculate distance for all the bodies in universe
	for i := 0; i < len(u.bodies); i++ {
		for j := i; j < len(u.bodies); j++ {
			dist := Distance(u.bodies[i].position, u.bodies[j].position)

			//Check if dist is the max distance between the two points
			if dist > maxDist {
				//make dist the new max distance
				maxDist = dist
			}

		}

	}

	return maxDist

}

// takes a slice u of Universe objects as input
// returns a slice of float64 variables having the same length as u.bodies
func AverageSpeed(u []Body) []float64 {
	//Get length of bodies in universe
	numBodies := len(u)

	//Initalize average speed
	//Make array to store average speed of each body
	AvgSpeed := make([]float64, numBodies)

	//Iterate over each body
	for i := 0; i < numBodies; i++ {
		CombinedSpeed := Speed(u[i].velocity)

		// Calculate the average speed for i-th body
		AvgSpeed[i] = CombinedSpeed
	}

	//return slice
	return AvgSpeed
}

func Speed(velocity OrderedPair) float64 {
	return math.Sqrt(velocity.x*velocity.x + velocity.y*velocity.y)

}
