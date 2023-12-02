package main

import (
	"canvas"
	"image"
)

//place your drawing code here.

// let's place our drawing functions here.

//AnimateSystem takes a slice of Sky objects along with a canvas width
//parameter and generates a slice of images corresponding to drawing each Sky
//on a canvasWidth x canvasWidth canvas

// Drawing Sky slice if it s divisible by drawing frequency
func AnimateSystem(timePoints []Country, canvasWidth, drawingFrequency int) []image.Image {
	images := make([]image.Image, 0)

	for i := range timePoints {
		if i%drawingFrequency == 0 {
			images = append(images, DrawToCanvas(timePoints[i], canvasWidth))
		}
	}
	return images
}

// GetFlyColor returns the color for a fly based on its stage
func GetFlyColor(stage int, isAlive bool) canvas.Color {
	if !isAlive {
		return canvas.MakeColor(0, 0, 0) // Black for dead flies
	}

	// Add color based on the fly's stage
	switch stage {
	case 0:
		return canvas.MakeColor(255, 0, 0) // Red for egg
	case 1:
		return canvas.MakeColor(255, 165, 0) // Orange for instar1
	case 2:
		return canvas.MakeColor(255, 255, 0) // Yellow for instar2
	case 3:
		return canvas.MakeColor(0, 255, 0) // Green for instar3
	case 4:
		return canvas.MakeColor(0, 0, 255) // Blue for instar4
	case 5:
		return canvas.MakeColor(128, 0, 128) // Purple for adult
	default:
		return canvas.MakeColor(255, 255, 255) // White for unknown stage
	}
}

// DrawToCanvas generates the image corresponding to a canvas after drawing a Quadrant
// object's flies on a square canvas that is canvasWidth pixels x canvasWidth pixels
func DrawToCanvas(country Country, canvasWidth int) image.Image {
	// set a new square canvas
	c := canvas.CreateNewCanvas(canvasWidth, canvasWidth)

	// create a black background
	c.SetFillColor(canvas.MakeColor(0, 0, 0))
	c.ClearRect(0, 0, canvasWidth, canvasWidth)
	c.Fill()

	// range over all the flies and draw them.
	for _, fly := range country.Fly {
		// Get the color based on the fly's stage and alive status
		color := GetFlyColor(fly.Stage, fly.isAlive)

		// Set the fly color
		c.SetFillColor(color)

		cx := (fly.Position.X / country.Width) * float64(canvasWidth)
		cy := (fly.Position.Y / country.Width) * float64(canvasWidth)
		r := 5
		c.Circle(cx, cy, float64(r))
		c.Fill()
	}

	// we want to return an image!
	return c.GetImage()
}
