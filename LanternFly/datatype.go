package main

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
	id    int
	temp  float64
}

// Coordinate represents geographical coordinates with latitude and longitude.
type Coordinate struct {
	Latitude  float64
	Longitude float64
}

const (
	instar1To2Threshold     float64 = 166.6
	instar2To3Threshold     float64 = 208.7
	instar3To4Threshold     float64 = 410.5
	instar4ToAdultThreshold float64 = 620
	// TODO: adult to die
	//lowerTempThreshold float64 =

	// survival rate: 1: 0.6488, 2: 0.9087, 3: 0.8948, 4: 0.822
	sRI1 float64 = 0.6488
	sRI2 float64 = 0.9087
	sRI3 float64 = 0.8948
	sRI4 float64 = 0.822
	sRA  float64 = 0.5837

	earthRadius float64 = 6371 // km

	minLat float64 = 24.396308  // Southernmost point in the US
	maxLat float64 = 49.384358  // Northernmost point in the contiguous US
	minLon float64 = -125.00165 // Westernmost point in the contiguous US
	maxLon float64 = -66.93457  // Easternmost point in the contiguous US
)
