package main

import (
	"canvas"
	"image"
	"image/color"
)

//AnimateSystem takes a slice of Country objects along with a canvas width
//parameter and generates a slice of images corresponding to drawing each Country
//on a canvasWidth x canvasWidth canvas

// Drawing Country slice if it is divisible by drawing frequency
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
func GetFlyColor(fly Fly) color.Color {
	if !fly.isAlive {
		return canvas.MakeColor(0, 0, 0) // Black for dead flies
	}

	// Use SorttheFlies to get the color based on the fly's stage
	sortedFly := SorttheFlies(fly)

	return canvas.MakeColor(
		sortedFly.color.red,
		sortedFly.color.green,
		sortedFly.color.blue,
	)
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
	for _, fly := range country.flies {
		// Get the color based on the fly's stage and alive status
		color := GetFlyColor(fly)

		// Set the fly color
		c.SetFillColor(color)

		cx := (fly.position.x / country.width) * float64(canvasWidth)
		cy := (fly.position.x / country.width) * float64(canvasWidth)
		r := 5
		c.Circle(cx, cy, float64(r))
		c.Fill()
	}

	// we want to return an image!
	return c.GetImage()
}
