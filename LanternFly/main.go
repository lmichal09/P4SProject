package main

import (
	"fmt"
	"gifhelper"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
)

// BUG: everything in io.go is commented out

func main() {
	fmt.Println("Lantern Flies simulation!")
	// step 1: reading input from a single file.

	fmt.Println("Success! Now we are ready to do something cool with our data.")

	outputFile := "output/output.gif" // Define the output file path and name

	initialCountry := InitializeCountry()
	fmt.Println("Country initialized.")

	quadrants := InitializeQuadrants()
	fmt.Println("Quadrants initialized.")

	timePoints := SimulateMigration(initialCountry, 3, quadrants)
	fmt.Println("Migration simulated.")

	canvasWidth := 1000
	imageFrequency := 1

	// Call your AnimateSystem function to generate images
	images := AnimateSystem(timePoints, canvasWidth, imageFrequency) //error

	fmt.Println("Images drawn!")

	fmt.Println("Generating an animated GIF.")

	// Save the images as an animated GIF
	gifhelper.ImagesToGIF(images, outputFile)

	fmt.Println("GIF drawn!")

	fmt.Println("Simulation complete!")
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
