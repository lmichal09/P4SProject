package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
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
	// TODO: Implement logic to initialize the habitat based on coordinates
	// You may use the InitializeHabitat function from your existing code

	// Placeholder - Replace with actual logic
	width := 100.0
	numFlies := 100
	numPredators := 10
	return InitializeHabitat(width, numFlies, numPredators)
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
