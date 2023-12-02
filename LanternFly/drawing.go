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

// DrawToCanvas generates the image corresponding to a canvas after drawing a Universe
// object's bodies on a square canvas that is canvasWidth pixels x canvasWidth pixels
func DrawToCanvas(c Country, canvasWidth int) image.Image {
	// set a new square canvas
	c := canvas.CreateNewCanvas(canvasWidth, canvasWidth)

	// create a black background
	c.SetFillColor(canvas.MakeColor(0, 0, 0))
	c.ClearRect(0, 0, canvasWidth, canvasWidth)
	c.Fill()

	// range over all the bodies and draw them.
	for _, b := range s.Boids {
		c.SetFillColor(canvas.MakeColor(255, 60, 25))
		cx := (b.Position.X / s.Width) * float64(canvasWidth)
		cy := (b.Position.Y / s.Width) * float64(canvasWidth)
		r := 5
		c.Circle(cx, cy, float64(r))
		c.Fill()

	}
	// we want to return an image!
	return c.GetImage()
}
