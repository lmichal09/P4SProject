package main

import (
	"math/rand"
	"time"
)

// InitializeCountry creates and initializes a Country with 10 flies and 5 trees with different positions for testing.
func InitializeCountry() Country {
	// Define the width of the country (you can change this value)
	countryWidth := 100.0

	// Create 10 sample flies with different positions
	flies := make([]Fly, 10)
	for i := 0; i < 10; i++ {
		fly := Fly{
			position:   OrderedPair{x: rand.Float64() * 100.0, y: rand.Float64() * 100.0}, // Different positions for each fly
			stage:      0,                                                                 // Initialize the stage as needed
			energy:     0,                                                                 // Initialize the energy as needed
			isAlive:    true,                                                              // Set initial status
			locationID: 0,                                                                 // Initialize the location ID as needed
			color:      Color{red: 255, green: 0, blue: 0},                                // Initialize the color as needed
		}
		flies[i] = fly
	}

	// Create 5 sample trees with different positions
	trees := make([]Tree, 5)
	for i := 0; i < 5; i++ {
		tree := Tree{
			position: OrderedPair{x: rand.Float64() * 100.0, y: rand.Float64() * 100.0}, // Different positions for each tree
		}
		trees[i] = tree
	}

	// Create a Country with the defined width, flies, and trees
	country := Country{
		width: countryWidth,
		flies: flies,
		trees: trees,
	}

	return country
}

// InitializeWeather creates and initializes 25 different weather conditions for the simulation.
func InitializeWeather() Weather {
	// Seed the random number generator to ensure different weather conditions each time
	rand.Seed(time.Now().UnixNano())

	// Define the temperature ranges for each quadrant
	minTemperature := 20.0 // Minimum temperature in degrees Celsius
	maxTemperature := 30.0 // Maximum temperature in degrees Celsius

	// Create a list of 25 quadrants with random temperature values
	quadrants := make([]Quadrant, 25)
	for i := 0; i < 25; i++ {
		// Generate a random temperature within the specified range
		temperature := minTemperature + rand.Float64()*(maxTemperature)

		// Calculate x and y coordinates for each quadrant to cover a 5x5 grid
		xCoord := float64(i%5) * 20.0 // 20 units width per quadrant
		yCoord := float64(i/5) * 20.0 // 20 units height per quadrant

		// Create a quadrant with the random temperature
		quadrant := Quadrant{
			x:     xCoord,      // x-coordinate (bottom left point)
			y:     yCoord,      // y-coordinate (bottom left point)
			temp:  temperature, // Random temperature for the quadrant
			width: 20.0,        // Width of each quadrant
			id:    i + 1,       // Quadrant ID
		}
		quadrants[i] = quadrant
	}

	// Create the Weather object with the list of quadrants
	weather := Weather{
		x:         0.0,       // Adjust the x-coordinate as needed
		y:         0.0,       // Adjust the y-coordinate as needed
		Quadrants: quadrants, // List of quadrants with random temperature data
	}

	return weather
}
