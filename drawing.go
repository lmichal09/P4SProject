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
// func AnimateSystem(timePoints []Country, canvasWidth, drawingFrequency int) []image.Image {
// 	images := make([]image.Image, 0)

// 	for i := range timePoints {
// 		if i%drawingFrequency == 0 {
// 			images = append(images, DrawToCanvas(timePoints[i], canvasWidth))
// 		}
// 	}
// 	return images
// }

func AnimateSystem(dailyTimePoints []Country, canvasWidth int, frequency int) []image.Image {
	images := make([]image.Image, 0)

	for i := range dailyTimePoints {
		if i%frequency == 0 {
			images = append(images, DrawToCanvas(dailyTimePoints[i], canvasWidth))
		}
	}
	return images
}

// GetFlyColor returns the color for a fly based on its stage
func GetFlyColor(fly Fly) color.Color {
	if !fly.isAlive {
		return canvas.MakeColor(10, 20, 10) // Black for dead flies
	}

	switch fly.stage {
	case 0:
		return canvas.MakeColor(255, 0, 0) // Red for egg
	case 1:
		return canvas.MakeColor(255, 165, 0) // Orange for instar1
	case 2:
		return canvas.MakeColor(255, 255, 0) // Yellow for instar2
	case 3:
		return canvas.MakeColor(0, 128, 0) // Green for instar3
	case 4:
		return canvas.MakeColor(0, 0, 255) // Blue for instar4
	case 5:
		return canvas.MakeColor(128, 0, 128) // Purple for adult
	default:
		return canvas.MakeColor(255, 255, 255) // Black for unknown stage
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
	for _, fly := range country.flies {
		// Get the color based on the fly's stage and alive status
		color := GetFlyColor(fly)

		// Set the fly color
		c.SetFillColor(color)

		cx := (fly.position.x / float64(country.width)) * float64(canvasWidth)
		cy := (fly.position.y / float64(country.width)) * float64(canvasWidth)
		r := 5
		c.Circle(cx, cy, float64(r))
		c.Fill()
	}

	for _, tree := range country.trees {
		c.SetFillColor(canvas.MakeColor(0, 175, 0))
		cx := (tree.position.x / float64(country.width)) * float64(canvasWidth)
		cy := (tree.position.y / float64(country.width)) * float64(canvasWidth)
		r := 10
		c.Circle(cx, cy, float64(r))
		c.Fill()
	}

	// we want to return an image!
	return c.GetImage()
}
