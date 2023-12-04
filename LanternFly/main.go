package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"math/rand"
	"os"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

func main() {
	fmt.Println("Lantern Flies simulation!")

	// Set up an HTTP server with the upload handler
	http.HandleFunc("/", uploadHandler)
	http.ListenAndServe(":8080", nil)

	outputFile := "output/output.gif" // Define the output file path and name

	fmt.Println("CLAs read!")

	fmt.Println("Now, simulating boids.")

	// Declare all Fly objects
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
	// Open the CSV file
    file, err := os.Open("tree.csv")
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()

    // Create a CSV reader
    reader := csv.NewReader(file)

    // Read the file
    records, err := reader.ReadAll()
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }

    // Skip the header row and process the data
    var habitats []OrderPair
    for i, record := range records {
        if i == 0 { // Skip header
            continue
        }

        // Parse longitude
        longitude, err := strconv.ParseFloat(record[0], 64) // Assuming longitude is in the first column
        if err != nil {
            fmt.Printf("Error parsing longitude in row %d: %v\n", i+1, err)
            continue
        }

        // Parse latitude
        latitude, err := strconv.ParseFloat(record[1], 64) // Assuming latitude is in the second column
        if err != nil {
            fmt.Printf("Error parsing latitude in row %d: %v\n", i+1, err)
            continue
        }

        // Append the habitat to the slice
        habitats = append(habitats, Habitat{Longitude: longitude, Latitude: latitude})
    }
	fmt.Println("Success! Now we are ready to do something cool with our data.")
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
