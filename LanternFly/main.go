package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"math/rand"
	"os"
)

func main() {
	// Read longitude and latitude coordinates from the text file
	coordinates, err := ReadCoordinates("lydetext.txt")
	if err != nil {
		fmt.Printf("Error reading coordinates: %v\n", err)
		return
	}

	// Create an initial habitat based on the coordinates
	initialHabitat := CreateInitialHabitat(coordinates)

	// Simulate lantern fly migration over multiple generations
	numGenerations := 100
	timePerGeneration := 1.0 // Time for each generation
	timePoints := SimulateMigration(initialHabitat, numGenerations, timePerGeneration)

	// Create images for each time point in the simulation
	canvasWidth := 500
	drawingFrequency := 1
	images := AnimateSystem(timePoints, canvasWidth, drawingFrequency)

	// Save the images as a GIF
	err = SaveGIF(images, "lantern_migration.gif")
	if err != nil {
		fmt.Printf("Error saving GIF: %v\n", err)
		return
	}

	fmt.Println("Simulation completed. GIF saved as 'lantern_migration.gif'")
}

// ReadCoordinates reads longitude and latitude coordinates from a text file
func ReadCoordinates(filename string) ([]Coordinate, error) {
	// Use the provided code to read sample data from the text file
	sampleData := ReadSampleDataFromFile(filename)

	// Convert SampleData to Coordinate
	coordinates := make([]Coordinate, len(sampleData))
	for i, data := range sampleData {
		coordinates[i] = Coordinate{
			Latitude:  data.Latitude,
			Longitude: data.Longitude,
		}
	}

	return coordinates, nil
}

// CreateInitialHabitat creates the initial habitat based on the given coordinates
func CreateInitialHabitat(coordinates []Coordinate) Country {
	country := Country{}  // Create the initial country.
	country.width = 100.0 // Set a default width, you can adjust this value as needed

	// Initialize the states
	//country.states = make([]Quadrant, 25)
	//country.states = LoadStates("states.csv")

	// Initialize the flies based on the provided coordinates
	numFlies := len(coordinates)
	country.flies = make([]Fly, numFlies)

	for i, coord := range coordinates {
		// Use coordinates to set the initial position of flies
		country.flies[i].position.x = coord.Longitude
		country.flies[i].position.y = coord.Latitude

		// Velocity and acceleration are random since no data is available
		country.flies[i].velocity.x = rand.Float64() * 2
		country.flies[i].velocity.y = rand.Float64() * 5
		country.flies[i].acceleration.x = rand.Float64() * rand.Float64() * 2
		country.flies[i].acceleration.y = rand.Float64() * rand.Float64() * 5

		// Lantern fly's stage random from 1-4
		country.flies[i].stage = rand.Intn(4) + 1

		// Initialize the energy of flies based on their stage
		switch country.flies[i].stage {
		case 1: // egg
			country.flies[i].energy = rand.Float64() * 39.5
		case 2: // nymph1
			country.flies[i].energy = rand.Float64() * 250
		case 3: // nymph2
			country.flies[i].energy = rand.Float64() * 108.7
		case 4: // adult
			country.flies[i].energy = rand.Float64() * 180
		}

		// When initialized, consider all flies are alive
		country.flies[i].isAlive = true

		// LocationID is random from 0-24
		country.flies[i].locationID = rand.Intn(25) // Get data from file, can be changed later
	}

	// Initialize the predators.
	numPredators := 10
	country.predators = make([]Predator, numPredators)

	for i := range country.predators {
		country.predators[i].position.x = rand.Float64() * country.width
		country.predators[i].position.y = rand.Float64() * country.width
		country.predators[i].velocity.x = rand.Float64() * 2
		country.predators[i].velocity.y = rand.Float64() * 5
		country.predators[i].acceleration.x = rand.Float64() * rand.Float64() * 2
		country.predators[i].acceleration.y = rand.Float64() * rand.Float64() * 5
		country.predators[i].PercentEaten = rand.Float64() * 100
	}

	return country
}

// SaveGIF saves a sequence of images as a GIF file
func SaveGIF(images []image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Set up the GIF header
	gifWriter := gif.GIF{}
	for _, img := range images {
		// Convert each image to a paletted image for GIF encoding
		palettedImg := image.NewPaletted(img.Bounds(), color.Palette{color.White, color.Black})
		draw.Draw(palettedImg, palettedImg.Bounds(), img, image.Point{}, draw.Over)

		gifWriter.Image = append(gifWriter.Image, palettedImg)
		gifWriter.Delay = append(gifWriter.Delay, 0) // No delay between frames (change if needed)
	}

	// Encode the GIF and write to the file
	err = gif.EncodeAll(file, &gifWriter)
	if err != nil {
		return err
	}

	return nil
}
