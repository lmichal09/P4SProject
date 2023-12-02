package main

import (
	"fmt"
)

func main() {
	fmt.Println("Lantern Flies simulation!")

	// step 1: reading input from a single file.

	filename := "Data/lydetext.txt"
	allData := ReadSamplesFromDirectory(filename)

	//step 2: reading input from a directory

	for sampleName, data := range allData {
		csvFilename := sampleName + ".csv"
		err := WriteToFile(csvFilename, data)
		if err != nil {
			fmt.Printf("Error writing to file %s: %v\n", csvFilename, err)
		} else {
			fmt.Printf("Data written to %s\n", csvFilename)
		}
	}

	fmt.Println("Success! Now we are ready to do something cool with our data.")
}
package main

import (
	"fmt"
	"gifhelper"
	"math"
	"math/rand"
	"os"
	"strconv"
)

// Define your Boid and Sky types here

func main() {
	fmt.Println("Sky Boids simulation!")

	// Declaring Sky and setting its fields.

	// now we need to implement the system

	//let's take command line arguments (CLAs) from the user
	//CLAs get stored in an ARRAY of strings called os.Args
	//this array has length equal to number of arguments given by the user + 1

	//os.Args[0] is the name of the program (./boids)
	fmt.Println(os.Args[0])

	//let's take CLAs: numGens, time, output path?, width of canvas

	numBoids, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	if numBoids < 0 {
		panic("Error: Invalid numBoids argument")
	}

	skyWidth, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		panic(err)
	}

	if skyWidth < 0 {
		panic("Error: Invalid skyWidth argument")
	}

	initialSpeed, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		panic(err)
	}

	if initialSpeed < 0 {
		panic("Error: Invalid initialSpeed argument")
	}

	maxBoidSpeed, err := strconv.ParseFloat(os.Args[4], 64)
	if err != nil {
		panic(err)
	}

	if maxBoidSpeed < 0 {
		panic("Error: Invalid maxBoidSpeed argument")
	}

	numGens, err := strconv.Atoi(os.Args[5])
	if err != nil {
		panic(err)
	}

	if numGens < 0 {
		panic("Error: Invalid numGens argument")
	}

	proximity, err := strconv.ParseFloat(os.Args[6], 64)
	if err != nil {
		panic(err)
	}

	if proximity < 0 {
		panic("Error: Invalid proximity argument")
	}

	separationFactor, err := strconv.ParseFloat(os.Args[7], 64)
	if err != nil {
		panic(err)
	}

	if separationFactor < 0 {
		panic("Error: Invalid separationForce argument")
	}

	alignmentFactor, err := strconv.ParseFloat(os.Args[8], 64)
	if err != nil {
		panic(err)
	}

	if alignmentFactor < 0 {
		panic("Error: Invalid alignmentFactor argument")
	}

	cohesionFactor, err := strconv.ParseFloat(os.Args[9], 64)
	if err != nil {
		panic(err)
	}

	if cohesionFactor < 0 {
		panic("Error: Invalid cohesionFactor argument")
	}

	timeStep, err := strconv.ParseFloat(os.Args[10], 64)
	if err != nil {
		fmt.Println("Error: Invalid timeStep argument")
		return
	}

	canvasWidth, err := strconv.Atoi(os.Args[11])
	if err != nil {
		fmt.Println("Error: Invalid canvasWidth argument")
		return
	}

	imageFrequency, err := strconv.Atoi(os.Args[12])
	if err != nil {
		fmt.Println("Error: Invalid imageFrequency argument")
		return
	}

	outputFile := "output/output.gif" // Define the output file path and name

	fmt.Println("CLAs read!")

	fmt.Println("Now, simulating boids.")

	// Declare all Sky objects
	var initialSky Sky
	initialSky.Width = skyWidth
	initialSky.Boids = make([]Boid, 0)
	initialSky.MaxBoidSpeed = maxBoidSpeed
	initialSky.Proximity = proximity
	initialSky.SeparationFactor = separationFactor
	initialSky.AlignmentFactor = alignmentFactor
	initialSky.CohesionFactor = cohesionFactor

	// Call the initializeBoids function to create the initial boids
	initialSky.Boids = initializeBoids(numBoids, initialSpeed, initialSky.Width)

	// Call your SimulateBoids function to perform the simulation
	timePoints := SimulateBoids(initialSky, numGens, timeStep) //error

	fmt.Println("Simulation complete!")

	fmt.Println("Drawing Skies.")

	// Call your AnimateSystem function to generate images
	images := AnimateSystem(timePoints, canvasWidth, imageFrequency) //error

	fmt.Println("Images drawn!")

	fmt.Println("Generating an animated GIF.")

	// Save the images as an animated GIF
	gifhelper.ImagesToGIF(images, outputFile)

	fmt.Println("GIF drawn!")

	fmt.Println("Simulation complete!")
}

// Initialize and generate randomly direction
func initializeBoids(numBoids int, initialSpeed float64, skyWidth float64) []Boid {
	boids := make([]Boid, numBoids)

	for i := range boids {
		// Generate a random unit vector for direction
		direction := randomUnitVector()

		// Calculate the initial velocity by scaling the direction vector
		boids[i].Velocity.X = direction.X * initialSpeed
		boids[i].Velocity.Y = direction.Y * initialSpeed

		// Set the initial position within the sky's width
		boids[i].Position.X = rand.Float64() * skyWidth
		boids[i].Position.Y = rand.Float64() * skyWidth
	}

	return boids
}

// Vector represents a 2D vector with X and Y components.
type Vector struct {
	X, Y float64
}

// randomUnitVector generates a random unit vector.
func randomUnitVector() Vector {
	// Generate a random angle between 0 and 2π (360 degrees)
	angle := rand.Float64() * 2 * math.Pi

	// Calculate the X and Y components of the unit vector
	x := math.Cos(angle)
	y := math.Sin(angle)

	// Create and return the unit vector
	return Vector{X: x, Y: y}
}
