//List functions needed from Lantern Fly Simulation
//Author: Leila Michal, Emma Bouchard, Tiffany Ku, and Thu Pham

package main

import (
	"math"
	"math/rand"
	"runtime"
)

// SimulateMigration simulates the migration of flies across the country over multiple years.
// The drawPoints array keeps track of the state of the country at the end of each day.
// In each year, the flies go through their lifecycle, with adults laying eggs and other stages changing over time.
// During the winter months, all flies die, except for the eggs laid by the adults.
// The simulation updates the state of the country and appends it to the drawPoints array for each day.
// Finally, the drawPoints array is returned, representing the state of the country at the end of each day for all years.
func SimulateMigration(initialCountry Country, numYears int, weather Weather) []Country {
	drawPoints := make([]Country, 0)
	drawPoints = append(drawPoints, initialCountry)
	currentCountry := CopyCountry(initialCountry)
	finalState := currentCountry
	for year := 0; year <= numYears; year++ {
		var totalEggs []Fly
		for i := 1; i <= 365; i++ {
			var eggs []Fly
			finalState = UpdateCountry(currentCountry, weather)
			// if adult, lay eggs
			// collect all eggs
			if finalState.flies[i].stage == 5 {
				eggs = ComputeFecundity(finalState.flies[i])
				totalEggs = append(totalEggs, eggs...)
			}

			drawPoints = append(drawPoints, finalState)

			if i > 214 { // Days from May to Dec
				for i := range finalState.flies {
					if finalState.flies[i].isAlive == true && finalState.flies[i].stage != 0 {
						finalState.flies[i].isAlive = false
					}
				}
			}

		}
		// collect all eggs
		currentCountry = finalState

		if CheckDead(finalState.flies) {
			finalState.flies = totalEggs
		}
	}
	return drawPoints
}

// UpdateCountry takes a current country and weather data as parameters,
// creates a new copy of the country, updates the fly population in parallel based on the weather data,
// and returns the updated country.
func UpdateCountry(currentCountry Country, weather Weather) Country {
	newcountry := CopyCountry(currentCountry) //copy current country

	numProcs := runtime.NumCPU() //get number of CPUs

	// update flies
	UpdateFlyMultiProcs(newcountry.flies, weather, newcountry.trees, numProcs)

	return newcountry
}

// UpdateFlyMultiProcs updates the flies in parallel
// takes a slice of flies and a number of processors.
// It divides the slice of flies into approximately equal parts, and sends each part to a separate goroutine for processing.
// It uses a finished channel to wait for all the goroutines to finish.
func UpdateFlyMultiProcs(fly []Fly, weather Weather, trees []Tree, numProcs int) {
	numFlies := len(fly)

	finished := make(chan bool)

	for i := 0; i < numProcs; i++ {
		// each processor getting ~ numParticles/numProcs particles

		startIndex := i * numFlies / numProcs
		endIndex := (i + 1) * numFlies / numProcs

		go UpdateFlySingleProc(fly[startIndex:endIndex], weather, trees, finished)
	}

	for i := 0; i < numProcs; i++ {
		<-finished
	}

}

// UpdateFlySingleProc takes a slice of Fly instances, a Weather instance, a slice of Tree instances, and a finished channel as input.
// The function iterates over the fly slice using a for loop and range function.
// Inside the loop, it calls the UpdateFly function with the current Fly instance, Weather, and Tree slices as arguments.
// After the loop, the function sends a value through the finished channel to signal that the update process is finished.
func UpdateFlySingleProc(fly []Fly, weather Weather, trees []Tree, finished chan bool) {
	for i := range fly {
		fly[i] = UpdateFly(fly[i], weather, trees)
	}
	finished <- true
}

// UpdateFly takes a fly, weather, and trees as parameters.
// It updates the fly's energy, position, life stage, and determines if the fly is alive or not.
// The updated fly is then returned.
func UpdateFly(fly Fly, weather Weather, trees []Tree) Fly {
	// Compute degree-day to affect fly's energy
	fly.energy = ComputeDegreeDay(&fly, weather)

	// Compute movement based on fly's energy and tree locations
	fly.position = ComputeMovement(&fly, trees)

	// Update fly's life stage based on age and conditions
	fly.stage = UpdateLifeStage(&fly)

	// Check if fly has died based on its current condition
	fly.isAlive = ComputeMortality(&fly)

	return fly
}

// phiE calculates the calculates the difference between two terms of a function, returning the absolute value of their difference.
// It defines a constant value 'k' and calculates the two terms of the function using 'math.Exp'.
// Finally, it returns the absolute difference between the two terms.
func phiE(d float64) float64 {
	// Define the constant value
	const k = 0.012

	// Calculate each term of the function
	term1 := 1 / (k*math.Exp(d-1) - 1)
	term2 := 1 / (k*math.Exp(d) - 1)

	// Return the absolute difference
	return math.Abs(term1 - term2)
}

// ComputeFecundity computes the fecundity of a fly.
// The function takes a fly as an argument and returns a slice of new flies that result from the laying of eggs.
// The code generates a new egg fly at the location of the adult fly.
// The number of egg masses, number of eggs in each mass, and the total number of eggs are randomly determined.
// The probability of laying eggs is determined by the fly's energy level.
// females lay one or two egg masses, each containing 30-60 eggs
func ComputeFecundity(fly Fly) []Fly {
	newFly := make([]Fly, 0)

	probToLayEggs := phiE(fly.energy)

	if rand.Float64() > probToLayEggs {
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
	}

	return newFly
}

// GetTreePositions iterates through each tree in the provided country.
// It adds each tree to a new slice and returns it.
func GetTreePositions(country Country) []Tree {
	var trees []Tree

	// Loop through country's trees
	for i := 0; i < len(country.trees); i++ {
		trees = append(trees, country.trees[i])
	}

	return trees
}

// ComputeDegreeDay calculates the degree days for a single day.
// takes a pointer to a fly and a weather struct as inputs.
// It computes the degree days for the fly by considering the fly's stage and the temperature of the quadrant it is in.
// The result is then returned.
func ComputeDegreeDay(fly *Fly, weather Weather) float64 {
	// get the quadrant of the fly to determine the temperature
	quadrantID := GetQuadrant(fly, weather.Quadrants)

	// get the temperature of the quadrant
	temperature := GetTemperature(quadrantID, weather.Quadrants)

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
	// 1. Default stage base temp is 0.0
	switch stage {
	case 1:
		return 13.00 // 2. First stage base temp is 13.00
	case 2:
		return 12.43 // 3. Second stage base temp is 12.43
	case 3:
		return 8.48 // 4. Third stage base temp is 8.48
	case 4:
		return 6.29 // 5. Fourth stage base temp is 6.29
	case 5:
		return 5.0 // 6. Fifth stage base temp is 5.0
	default:
		return 0.0 // 7. If stage is out of range, return default temp
	}
}

// UpdateLifeStage() updates the life stage of flies based on the cumulative degree-days (CDD)
func UpdateLifeStage(fly *Fly) int {
	// Update fly's life stage based on energy levels
	if fly.energy >= 0 && fly.energy < instar1To2Threshold {
		return 1 // Instar 1
	} else if fly.energy >= instar1To2Threshold && fly.energy < instar2To3Threshold {
		return 2 // Instar 2
	} else if fly.energy >= instar2To3Threshold && fly.energy < instar3To4Threshold {
		return 3 // Instar 3
	} else if fly.energy >= instar3To4Threshold && fly.energy < instar4ToAdultThreshold {
		return 4 // Instar 4
	} else if fly.energy >= instar4ToAdultThreshold && fly.energy < adultToDieThreshold {
		return 5 // Adult
	} else {
		return 6 // Dead
	}
}

// ComputeMortality updates the mortality status of flies.
// survival rate: 1: 0.6488, 2: 0.9087, 3: 0.8948, 4: 0.822
// calculates the mortality of a fly based on its stage.
// The stage is a number from 1 to 5, and the mortality rate for each stage is determined by a set of survival rates (sRI1 to sRI5).
// The function uses a switch statement to select the appropriate survival rate based on the fly's stage, and then generates a random float between 0 and 1.
// If the random float is less than or equal to the survival rate, the function returns true, indicating that the fly has died.
// If the fly's stage is invalid or the random float is greater than the survival rate, the function returns false, indicating that the fly is still alive.
func ComputeMortality(fly *Fly) bool {
	// Compute mortality based on stage and survival rates
	switch fly.stage {
	case 1:
		return rand.Float64() <= sRI1 // stage 1 survival rate
	case 2:
		return rand.Float64() <= sRI2 // stage 2 survival rate
	case 3:
		return rand.Float64() <= sRI3 // stage 3 survival rate
	case 4:
		return rand.Float64() <= sRI4 // stage 4 survival rate
	case 5:
		return rand.Float64() <= sRA // stage 5 adult survival rate
	default:
		// Handle invalid stages, consider them dead
		return false
	}
}

// ComputeMovement updates the position of adult flies
// determines the movement of a Fly instance.
// It has a 70% chance of executing RandomMovement and a 30% chance of executing DirectedMovement.
func ComputeMovement(fly *Fly, trees []Tree) OrderedPair {
	// Randomly decide between random movement and directed movement
	if rand.Float64() < 0.7 {

		// Random movement: flies move randomly within a certain distance
		return RandomMovement(fly)
	} else {

		// Directed movement: flies move towards the nearest host tree
		return DirectedMovement(fly, trees)
	}
}

// RandomMovement updates the position of adult flies based on random movement
// takes a fly object and updates its position randomly within a maximum distance.
// It uses trigonometry to calculate the new position based on the angle and maximum distance.
func RandomMovement(fly *Fly) OrderedPair {
	var maxDistance float64
	if rand.Float64() < 0.7 {
		maxDistance = 90.0
	} else {
		maxDistance = 10.0
	}

	angle := rand.Float64() * 2 * math.Pi // Random angle between 0 and 2*Pi radians

	// Calculate the new position based on random movement
	newX := fly.position.x + maxDistance*math.Cos(angle)
	newY := fly.position.y + maxDistance*math.Sin(angle)

	return OrderedPair{newX, newY}
}

// DirectedMovement updates the position of adult flies based on directed movement
// implements directed movement for a fly.
// The fly is attracted to the nearest host tree, and its position is updated based on the displacement vector between the fly and the tree.
// A small random jitter is added to the new position to create more realistic movement.
func DirectedMovement(fly *Fly, trees []Tree) OrderedPair {
	// Find the nearest host tree
	nearestTree := FindNearestTree(fly.position, trees)

	dx := nearestTree.x - fly.position.x
	dy := nearestTree.y - fly.position.y

	// Calculate the new position based on directed movement towards the nearest tree
	newX := fly.position.x + rand.Float64()*dx
	newY := fly.position.y + rand.Float64()*dy

	return OrderedPair{newX, newY}
}

// FindNearestTree finds the nearest host tree to a given fly's position.
// hat takes two parameters: a fly's position (flyPosition) and a slice of trees (trees).
// It returns the position of the nearest tree.
// The function starts by checking if there are any trees in the slice.
// If there are no trees, it returns a default value.
// If there are trees, the function initializes variables to store the nearest tree and its distance from the fly.
// It then iterates through the list of trees to find the nearest one.
// If a tree is closer to the fly than the current nearest tree, the function updates the nearest tree and its distance.
// Finally, the function returns the position of the nearest tree.
func FindNearestTree(flyPosition OrderedPair, trees []Tree) OrderedPair {
	if len(trees) == 0 {
		// Handle the case where there are no trees
		return OrderedPair{0, 0} // You can change this default value
	}

	// Initialize variables to store the nearest tree and its distance
	nearestTree := trees[0]
	minDistance := distance(flyPosition, nearestTree.position)

	// Iterate through the list of trees to find the nearest one
	for _, tree := range trees {
		d := distance(flyPosition, tree.position)
		if d < minDistance {
			minDistance = d
			nearestTree = tree
		}
	}

	return nearestTree.position
}

// GetQuadrant checks the position of a Fly and determines which Quadrant it belongs to by iterating through the quadrants slice
// and comparing the Fly's x and y coordinates with the coordinates of each Quadrant.
// If a match is found, the function returns the id of the corresponding Quadrant.
// If no match is found, the function returns -1.
func GetQuadrant(fly *Fly, quadrants []Quadrant) int {

	// Iterate through quadrants
	for _, q := range quadrants {
		// Check if fly's x position lies within the quadrant
		if fly.position.x >= q.x && fly.position.x <= q.x+q.width {
			// Check if fly's y position lies within the quadrant
			if fly.position.y >= q.y && fly.position.y <= q.y+q.width {
				// If both conditions are met, return the quadrant id
				return q.id
			}
		}
	}

	// If no matching quadrant is found, return -1
	return -1
}

// GetTemperature returns the temperature of the quadrant based on its ID.
// If the quadrant ID is -1, the function returns 0 as the temperature.
// If the quadrant ID matches a quadrant in the quadrant slice, the function retrieves and returns the temperature of that quadrant.
// If the quadrant ID does not match any quadrant, the function returns 0 as the temperature.
func GetTemperature(quadrantID int, quadrant []Quadrant) float64 {
	temp := 0.0

	if quadrantID == -1 {
		return temp
	}

	// loop through the quadrant slice and find the quadrant with the matching id
	for _, q := range quadrant {
		if q.id == quadrantID {
			temp = q.temp
		}
	}

	return temp
}

// distance calculates the Euclidean distance between two points (2D).
// It calculates the differences in the x and y coordinates and then uses the Pythagorean theorem to find the distance.
func distance(point1 OrderedPair, point2 OrderedPair) float64 {
	dx := point2.x - point1.x
	dy := point2.y - point1.y
	return math.Sqrt(dx*dx + dy*dy)
}

// CheckDead takes a slice of Fly and checks if all the flies in the given list are dead.
// If even one fly is alive, the function returns false immediately.
// If the loop finishes without finding any alive flies, the function returns true.
func CheckDead(flies []Fly) bool {
	// Check each fly in the list
	for i := range flies {
		// Return false if fly is alive
		if flies[i].isAlive {
			return false
		}
	}
	// Return true if all flies are dead
	return true
}

// InBounds checks if the fly is within the simulation bounds.
// checks if a fly's position is within the bounds of a given set of quadrants.
// It returns true if the fly's position is within a quadrant and false otherwise.
func InBounds(fly *Fly, quadrants []Quadrant) bool {
	// Check if fly position is within quadrants
	for _, q := range quadrants {
		if fly.position.x >= q.x && fly.position.x <= q.x+q.width &&
			fly.position.y >= q.y && fly.position.y <= q.y+q.width {
			return true
		}
	}
	return false
}

// CopyCountry creates a new copy of a given Country object.
// It creates a new Country instance and initializes it with the values from the original country.
// The flies and trees arrays are also copied using the CopyFly and CopyTree functions respectively.
// The new country is then returned.
func CopyCountry(original Country) Country {
	// Create a new Country instance
	copyCountry := Country{
		width: original.width,
		flies: make([]Fly, len(original.flies)),
		trees: make([]Tree, len(original.trees)),
	}

	// Deep copy flies
	for i, fly := range original.flies {
		copyCountry.flies[i] = CopyFly(fly)
	}

	// Deep copy trees
	for i, tree := range original.trees {
		copyCountry.trees[i] = CopyTree(tree)
	}

	return copyCountry
}

// CopyFly creates a copy of the original Fly instance by using a new Fly instance (copyFly) with the same values for all fields as the original Fly instance.
// It then returns the copied Fly instance.
func CopyFly(original Fly) Fly {
	// Create a new Fly instance
	copyFly := Fly{
		position: CopyOrderedPair(original.position),
		stage:    original.stage,
		energy:   original.energy,
		isAlive:  original.isAlive,
		color:    original.color,
	}

	return copyFly
}

// CopyTree creates a copy of a tree by copying the root position of the original tree.
// It uses the CopyOrderedPair function to create a copy of the ordered pair.
func CopyTree(original Tree) Tree {
	// Create a new Tree instance
	copyTree := Tree{
		position: CopyOrderedPair(original.position),
	}

	return copyTree
}

// CopyOrderedPair an OrderedPair as an argument and creates a copy of it.
// The function then returns the copied OrderedPair.
func CopyOrderedPair(original OrderedPair) OrderedPair {
	// Create a new OrderedPair instance
	copyPair := OrderedPair{
		x: original.x,
		y: original.y,
	}

	return copyPair
}
