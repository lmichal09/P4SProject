package main

import (
	"encoding/csv"
	"fmt"
	"gifhelper"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"math/rand"
	"os"
	"strconv"
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
	var habitats []OrderedPair
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
		habitats = append(habitats, OrderedPair{Longitude: longitude, Latitude: latitude})
	}
	fmt.Println("Success! Now we are ready to do something cool with our data.")

	// Call your AnimateSystem function to generate images
	images := AnimateSystem(timePoints, canvasWidth, imageFrequency) //error

	fmt.Println("Images drawn!")

	fmt.Println("Generating an animated GIF.")

	// Save the images as an animated GIF
	gifhelper.ImagesToGIF(images, outputFile)

	fmt.Println("GIF drawn!")

	fmt.Println("Simulation complete!")
}

// CreateInitialHabitat initializes a Country with flies based on the provided coordinates.
func CreateInitialHabitat(coordinates []OrderedPair) Country {
	country := Country{}  // Create the initial country.
	country.width = 100.0 // Set a default width, you can adjust this value as needed

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
