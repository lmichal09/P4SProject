package main

import (
	"fmt"
	"gifhelper"
	"math"
	"math/rand"
)

func main() {
	fmt.Println("Lantern Flies simulation!")

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
