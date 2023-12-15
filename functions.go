//List functions needed from Lantern Fly Simulation
//Author: Leila Michal, Emma Bouchard, Tiffany Ku, and Thu Pham

//NOTE: "TODO": all the things that need to be done

package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
)

// SimulateMigration

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

		fmt.Println("total eggs", len(totalEggs))

		currentCountry = finalState

		if CheckDead(finalState.flies) {
			fmt.Println("All flies are dead!!!!!!!!!!")
			finalState.flies = totalEggs
		}
	}
	return drawPoints
}

func UpdateCountry(currentCountry Country, weather Weather) Country {
	newcountry := CopyCountry(currentCountry)

	numProcs := runtime.NumCPU()

	// update flies
	UpdateFlyMultiProcs(newcountry.flies, weather, newcountry.trees, numProcs)

	// n := len(newcountry.flies)

	// for i := 0; i < n; i++ {
	// 	newcountry.flies[i] = UpdateFly(newcountry.flies[i], weather, newcountry.trees)
	// }
	return newcountry
}

// UpdateFlyMultiProcs updates the flies in parallel
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

func UpdateFlySingleProc(fly []Fly, weather Weather, trees []Tree, finished chan bool) {
	for i := range fly {
		fly[i] = UpdateFly(fly[i], weather, trees)
	}
	finished <- true
}

func UpdateFly(fly Fly, weather Weather, trees []Tree) Fly {
	fly.energy = ComputeDegreeDay(&fly, weather)
	fly.position = ComputeMovement(&fly, trees)
	fly.stage = UpdateLifeStage(&fly)
	fly.isAlive = ComputeMortality(&fly)
	//fmt.Println("dd", fly.energy, "stage", fly.stage)
	return fly
}

// phiE calculates the given function
func phiE(d float64) float64 {
	// Define the constant value
	const k = 0.012

	// Calculate each term of the function
	term1 := 1 / (k*math.Exp(d-1) - 1)
	term2 := 1 / (k*math.Exp(d) - 1)

	// Return the absolute difference
	return math.Abs(term1 - term2)
}

// ComputeFecundity
// females lay one or two egg masses, each containing 30-60 eggs
func ComputeFecundity(fly Fly) []Fly {
	newFly := make([]Fly, 0)

	// TODO: probability of laying eggs z
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

func GetTreePositions(country Country) []Tree {
	var trees []Tree

	for i := 0; i < len(country.trees); i++ {
		trees = append(trees, country.trees[i])
	}

	return trees
}

// ComputeDegreeDay calculates the degree days for a single day.
func ComputeDegreeDay(fly *Fly, weather Weather) float64 {
	// get the quadrant of the fly to determine the temperature
	quadrantID := GetQuadrant(fly, weather.Quadrants)
	// fmt.Println("weather", weather.Quadrants)
	// fmt.Println("quadrant", quadrantID)

	// get the temperature of the quadrant
	temperature := GetTemperature(quadrantID, weather.Quadrants)
	// fmt.Println("temp", temperature)

	// get the base temperature base on the fly's stage
	baseTemp := GetBaseTemp(fly.stage)
	// fmt.Println("base temp", baseTemp)

	// calculate the degree days
	degreeDays := (temperature+baseTemp)/2 - baseTemp
	// fmt.Println("dd", degreeDays)

	// if degreeDays is negative, set it to 0
	if degreeDays < 0 {
		degreeDays = 0
	}

	return degreeDays
}

// GetQuadrant returns the quadrant of the fly.
// Base temperature used for calculating GDD for nymph. 1: 13.00°C, 2: 12.43°C, 3: 8.48°C, 4: 6.29°C
func GetBaseTemp(stage int) float64 {
	switch stage {
	case 1:
		return 13.00
	case 2:
		return 12.43
	case 3:
		return 8.48
	case 4:
		return 6.29
	case 5:
		return 5.0
	default:
		return 0.0
	}
}

// UpdateLifeStage() updates the life stage of flies based on the cumulative degree-days (CDD)
func UpdateLifeStage(fly *Fly) int {
	if fly.energy >= 0 && fly.energy < instar1To2Threshold {
		return 1
	} else if fly.energy >= instar1To2Threshold && fly.energy < instar2To3Threshold {
		return 2
	} else if fly.energy >= instar2To3Threshold && fly.energy < instar3To4Threshold {
		return 3
	} else if fly.energy >= instar3To4Threshold && fly.energy < instar4ToAdultThreshold {
		return 4
	} else if fly.energy >= instar4ToAdultThreshold && fly.energy < adultToDieThreshold {
		return 5
	} else {
		return 6 // Dead
	}
}

// ComputeMortality updates the mortality status of flies.
// survival rate: 1: 0.6488, 2: 0.9087, 3: 0.8948, 4: 0.822
func ComputeMortality(fly *Fly) bool {
	switch fly.stage {
	case 1:
		return rand.Float64() <= sRI1
	case 2:
		return rand.Float64() <= sRI2
	case 3:
		return rand.Float64() <= sRI3
	case 4:
		return rand.Float64() <= sRI4
	case 5:
		return rand.Float64() <= sRA
	default:
		// Handle invalid stages, consider them dead
		return false
	}
}

// ComputeMovement updates the position of adult flies
func ComputeMovement(fly *Fly, trees []Tree) OrderedPair {
	// Randomly decide between random movement and directed movement
	if rand.Float64() < 0.7 {
		// fmt.Println("random movement")
		// Random movement: flies move randomly within a certain distance
		return RandomMovement(fly)
	} else {
		// fmt.Println("directed movement")
		// Directed movement: flies move towards the nearest host tree
		return DirectedMovement(fly, trees)
	}
}

// RandomMovement updates the position of adult flies based on random movement
func RandomMovement(fly *Fly) OrderedPair {
	var maxDistance float64
	if rand.Float64() < 0.7 {
		maxDistance = 90.0
	} else {
		maxDistance = 10.0
	}

	angle := rand.Float64() * 2 * math.Pi // Random angle between 0 and 2*Pi radians

	// fmt.Println("original", fly.position.x, fly.position.y)
	// Calculate the new position based on random movement
	newX := fly.position.x + maxDistance*math.Cos(angle)
	newY := fly.position.y + maxDistance*math.Sin(angle)
	// fmt.Println("new", newX, newY)

	return OrderedPair{newX, newY}
}

// DirectedMovement updates the position of adult flies based on directed movement
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

func GetQuadrant(fly *Fly, quadrants []Quadrant) int {
	// var quadrant int
	// fmt.Println("fly", fly.position.x, fly.position.y)
	// fmt.Println("quadrant", quadrants)

	// loop through the quadrant slice and find the quadrant with the matching id
	for _, q := range quadrants {
		if fly.position.x >= q.x &&
			fly.position.x <= q.x+q.width &&
			fly.position.y >= q.y &&
			fly.position.y <= q.y+q.width {
			return q.id
		}
	}

	return -1 // Placeholder
}

// GetTemperature returns the temperature of the quadrant.
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

	// if quadrantID < 1 || quadrantID > len(quadrant) {
	// 	panic("invalid quadrant ID")
	// }

	return temp
}

// distance calculates the Euclidean distance between two points (2D).
func distance(point1 OrderedPair, point2 OrderedPair) float64 {
	dx := point2.x - point1.x
	dy := point2.y - point1.y
	return math.Sqrt(dx*dx + dy*dy)
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

// InBounds checks if the fly is within the simulation bounds.
func InBounds(fly *Fly, quadrants []Quadrant) bool {
	for _, q := range quadrants {
		if fly.position.x >= q.x && fly.position.x <= q.x+q.width &&
			fly.position.y >= q.y && fly.position.y <= q.y+q.width {
			return true
		}
	}
	return false
}

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

func CopyTree(original Tree) Tree {
	// Create a new Tree instance
	copyTree := Tree{
		position: CopyOrderedPair(original.position),
	}

	return copyTree
}

func CopyOrderedPair(original OrderedPair) OrderedPair {
	// Create a new OrderedPair instance
	copyPair := OrderedPair{
		x: original.x,
		y: original.y,
	}

	return copyPair
}
