package main

import (
	"fmt"
	"gifhelper"
)

// initializes a system, simulates migration, and generates an animated GIF to visualize the system.
func main() {
	fmt.Println("Lantern Flies simulation!")
	// Reading input

	fmt.Println("Success! Now we are ready to do something cool with our data.")

	// Initialize the system
	initialCountry := InitializeCountry()
	fmt.Println("Country initialized.")

	numYears := 1

	weather := InitializeQuadrants()
	fmt.Println("Quadrants initialized.")

	timePoints := SimulateMigration(initialCountry, numYears, weather)
	fmt.Println("Migration simulated.")

	canvasWidth := 10000
	canvasHeight := 10000
	imageFrequency := 30

	// Animate the system
	images := AnimateSystem(timePoints, int(canvasWidth), int(canvasHeight), imageFrequency)
	fmt.Println("Images drawn!")

	fmt.Println("Generating an animated GIF.")

	// Save the images as an animated GIF
	gifhelper.ImagesToGIF(images, "flies!")

	fmt.Println("GIF drawn!")

	fmt.Println("Simulation complete!")
}
