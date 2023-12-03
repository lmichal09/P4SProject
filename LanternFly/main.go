package main

import (
	"encoding/csv"
	"fmt"
	"gifhelper"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"

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
	habitats := make([]OrderedPair, 0)
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

// uploadHandler handles the HTTP requests to the root ("/") URL.
// It serves two main purposes:
//  1. If the request method is POST, it processes the uploaded image,
//     resizes it, pastes it onto a white background, and creates a GIF file.
//  2. If the request method is GET, it displays an HTML form allowing users
//     to upload an image.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method == http.MethodPost {
		// Retrieve the uploaded file from the request
		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Error reading file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Decode the uploaded image
		img, _, err := image.Decode(file)
		if err != nil {
			http.Error(w, "Error decoding image", http.StatusInternalServerError)
			return
		}

		// Resize the image to a smaller size (you can adjust the width and height)
		resizedImg := resize.Resize(200, 0, img, resize.Lanczos3)

		// Create a new image with a white background
		gifImg := image.NewRGBA(image.Rect(0, 0, resizedImg.Bounds().Dx(), resizedImg.Bounds().Dy()))
		draw.Draw(gifImg, gifImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

		// Paste the resized image onto the white background
		cutterImg, err := cutter.Crop(resizedImg, cutter.Config{
			Width:  gifImg.Bounds().Dx(),
			Height: gifImg.Bounds().Dy(),
			Mode:   cutter.Centered,
		})
		if err != nil {
			http.Error(w, "Error cropping image", http.StatusInternalServerError)
			return
		}
		draw.Draw(gifImg, gifImg.Bounds(), cutterImg, image.Point{}, draw.Over)

		// Create a GIF file
		outFile, err := os.Create("output.gif")
		if err != nil {
			http.Error(w, "Error creating output file", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		// Encode the GIF
		err = gif.Encode(outFile, gifImg, nil)
		if err != nil {
			http.Error(w, "Error encoding GIF", http.StatusInternalServerError)
			return
		}

		// Respond to the client with a success message
		fmt.Fprintln(w, "GIF created successfully")
	} else {
		// Display the HTML form to upload an image
		form := `<html><body><form action="/" method="post" enctype="multipart/form-data"><input type="file" name="image"><input type="submit" value="Upload"></form></body></html>`
		w.Write([]byte(form))
	}
}
