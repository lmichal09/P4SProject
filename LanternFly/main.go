package main

import (
	"fmt"
	"gifhelper"
)

// BUG: everything in io.go is commented out

func main() {
	fmt.Println("Lantern Flies simulation!")
	// step 1: reading input from a single file.

	fmt.Println("Success! Now we are ready to do something cool with our data.")

	// outputFile := "LanternFly_demo/output.gif" // Define the output file path and name

	initialCountry := InitializeCountry()
	fmt.Println("Country initialized.")

	numYears := 3

	weather := InitializeWeather()
	// fmt.Println(weather)
	fmt.Println("Quadrants initialized.")

	timePoints := SimulateMigration(initialCountry, numYears, weather)
	fmt.Println("Migration simulated.")

	canvasWidth := 1000
	canvasHeight := 1000
	imageFrequency := 10

	// Call your AnimateSystem function to generate images
	images := AnimateSystem(timePoints, canvasWidth, canvasHeight, imageFrequency) //error

	fmt.Println("Images drawn!")

	fmt.Println("Generating an animated GIF.")

	// Save the images as an animated GIF
	gifhelper.ImagesToGIF(images, "flies!")

	fmt.Println("GIF drawn!")

	fmt.Println("Simulation complete!")
}
