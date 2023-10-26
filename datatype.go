package main

type Fly struct {
	position, velocity, acceleration OrderedPair
	stage                            int
}

type Predator struct {
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
