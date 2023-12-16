package main

import (
	"canvas"
	"image"
	"image/color"
)

// AnimateSystem takes a slice of Country objects along with a canvas width
// parameter and generates a slice of images corresponding to drawing each Country
// It then creates an empty slice called images to store the resulting animation frames.
// The function iterates through the daily time points, and only appends a new frame to the images slice if the current index is divisible by the frequency.
// Finally, the function returns the slice of images representing the animated system.
func AnimateSystem(dailyTimePoints []Country, canvasWidth, canvasHeight int, frequency int) []image.Image {
	images := make([]image.Image, 0)

	for i := range dailyTimePoints {
		if i%frequency == 0 {
			images = append(images, DrawToCanvas(dailyTimePoints[i], canvasWidth, canvasHeight))
		}
	}
	return images
}

// GetFlyColor returns the color for a fly based on its stage
// takes a Fly object as an argument and returns a color.Color object.
// The color of the fly depends on its stage and whether it is alive or not.
// The switch statement handles this by checking the stage of the fly and setting the color accordingly.
// If the fly is dead, it returns black.
func GetFlyColor(fly Fly) color.Color {
	if !fly.isAlive {
		return canvas.MakeColor(0, 0, 0) // Black for dead flies
	}

	switch fly.stage {
	case 0:
		return canvas.MakeColor(255, 0, 0) // Red for egg
	case 1:
		return canvas.MakeColor(255, 165, 0) // Orange for instar1
	case 2:
		return canvas.MakeColor(255, 255, 0) // Yellow for instar2
	case 3:
		return canvas.MakeColor(255, 255, 255) // White for instar3
	case 4:
		return canvas.MakeColor(0, 0, 255) // Blue for instar4
	case 5:
		return canvas.MakeColor(128, 0, 128) // Purple for adult
	default:
		return canvas.MakeColor(0, 0, 2) // Black for unknown stage
	}
}

// DrawToCanvas generates the image corresponding to a canvas after drawing a Quadrant
// object's flies on a square canvas that is canvasWidth pixels x canvasWidth pixels
// takes a Country, canvas width, and canvas height as input and returns an image.Image.
// It creates a new canvas and sets a black background.
// It then draws the flies and trees in the country.
// The function returns the drawn image.
func DrawToCanvas(country Country, canvasWidth, canvasHeight int) image.Image {
	// set a new square canvas
	c := canvas.CreateNewCanvas(canvasWidth, canvasHeight)

	// create a black background
	c.SetFillColor(canvas.MakeColor(0, 0, 0))
	c.ClearRect(0, 0, canvasWidth, canvasHeight)
	c.Fill()

	// range over all the flies and draw them.
	for _, fly := range country.flies {
		// Get the color based on the fly's stage and alive status
		color := GetFlyColor(fly)

		// Set the fly color
		c.SetFillColor(color)

		cx := (fly.position.x / float64(country.width)) * float64(canvasHeight)
		cy := (fly.position.y / float64(country.width)) * float64(canvasHeight)
		r := 10
		c.Circle(cx, cy, float64(r))
		c.Fill()
	}

	//draw trees
	for _, tree := range country.trees {
		c.SetFillColor(canvas.MakeColor(0, 175, 0))
		cx := (tree.position.x / float64(country.width)) * float64(canvasHeight)
		cy := (tree.position.y / float64(country.width)) * float64(canvasHeight)
		r := 100
		c.Circle(cx, cy, float64(r))
		c.Fill()
	}

	// we want to return an image!
	return c.GetImage()
}
