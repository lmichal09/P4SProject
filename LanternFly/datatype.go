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

type Fly struct {
	position, velocity, acceleration OrderedPair
	stage                            int     // 0 = egg, 1 = larva, 2 = pupa, 3 = adult
	energy                           float64
	isAlive                          bool
type Fly struct {
	position, velocity, acceleration OrderedPair
	PercentConsumed                  float64
	stage                            int
}

type Predator struct {
	position, velocity, acceleration OrderedPair
	PercentEaten                     float64
}

type Country struct {
	width  float64
	flies  []Fly
	states []Quadrant
}

type OrderedPair struct {
	x float64
	y float64
}

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
