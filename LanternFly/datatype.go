package main

// TODO: need to define the following constants
/*
// thermal constants: the number of degree-days required for a development change to occur
// lower temperature threshold: the lowest temperature in which an insect can develop. Below the lower temperature threshold, the insect does not accumulate degree-days; therefore, it may die or enter diapause.
const (
    eggToLarvaThreshold float64 = // ... degree-days
    larvaToPupaThreshold float64 = // ... degree-days
    pupaToAdultThreshold float64 = // ... degree-days
	adultToSenileThreshold float64 = // ... degree-days
	lowerTempThreshold float64 = // ... degree-days
)
*/

type Country struct {
	width      float64
	flies      []Fly
	population int
}

type Fly struct {
	position, velocity, acceleration OrderedPair
	stage                            int     // 0 = egg, 1 = instar1, 2 = instar2, 3 = instar3, 4 = instar4, 5 = adult
	energy                           float64 // Degree-days
	isAlive                          bool
	locationID                       int
}

type stage struct {
	a, b       float64
	Tmin, Tmax float64
}

type OrderedPair struct {
	x float64
	y float64
}
<<<<<<< HEAD
=======

// QuadTree simply contains a pointer to the root.
// Another way of doing this would be type QuadTree *Node
type QuadTree struct {
	root *Node
}

// Node object contains a slice of children (this could just as easily be an array of length 4).
// A node refers to a star. Sometimes, the star will be a "dummy" star, sometimes it is a star in the
// universe, and sometimes it is nil. Every internal node points to a dummy star.
type Node struct {
	children []*Node
	fly      *Fly
	sector   Quadrant
}

// Quadrant is an object representing a sub-square within a larger universe.
type Quadrant struct {
	x     float64 //bottom left corner x coordinate
	y     float64 //bottom right corner y coordinate
	width float64
}

// Coordinate represents geographical coordinates with latitude and longitude.
type Coordinate struct {
	Latitude  float64
	Longitude float64
}
>>>>>>> 4402fb46c8e5a6a832321fcb264bc955ecd48dbd
