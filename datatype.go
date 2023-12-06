package main

type Country struct {
	width float64
	flies []Fly
	trees []Tree
}

type Tree struct {
	position OrderedPair
}

type Fly struct {
	position   OrderedPair
	stage      int     // 0 = egg, 1 = instar1, 2 = instar2, 3 = instar3, 4 = instar4, 5 = adult, 6= dead
	energy     float64 // Degree-days
	isAlive    bool
	locationID int
	color      Color // color to show on scatter plot (red, orange, yellow, green, blue, purple, black) neon colors

}

type Weather struct {
	x         float64 // Bottom left corner x coordinate (Longitude)
	y         float64 // Bottom left corner y coordinate (Latitude)
	Quadrants []Quadrant
}

type Color struct {
	red   uint8
	blue  uint8
	green uint8
}

type stage struct {
	a, b       float64
	Tmin, Tmax float64
}

type OrderedPair struct {
	x float64
	y float64
}

// Quadrant is an object representing a sub-square within a larger universe.
type Quadrant struct {
	x     float64 //bottom left corner x coordinate
	y     float64 //bottom right corner y coordinate
	width float64
	id    int
	temp  float64
}

// SampleData represents the structure of the data in the file
type SampleData struct {
	Source           string
	Year             int
	BioYear          int
	Latitude         float64
	Longitude        float64
	State            string
	LydePresent      bool
	LydeEstablished  bool
	LydeDensity      string
	SourceAgency     string
	CollectionMethod string
	PointID          string
	RoundedLongitude float64
	RoundedLatitude  float64
}

const (
	instar1To2Threshold     float64 = 6.6
	instar2To3Threshold     float64 = 10.7
	instar3To4Threshold     float64 = 18.5
	instar4ToAdultThreshold float64 = 24.5
	adultToDieThreshold     float64 = 35.0

	// survival rate: 1: 0.6488, 2: 0.9087, 3: 0.8948, 4: 0.822
	sRI1 float64 = 0.6488
	sRI2 float64 = 0.9087
	sRI3 float64 = 0.8948
	sRI4 float64 = 0.822
	sRA  float64 = 0.5837

	earthRadius float64 = 6371 // km

	minLat float64 = 31.33   // Southernmost point in the US
	maxLat float64 = 45.71   // Northernmost point in the contiguous US
	minLon float64 = -123.27 // Westernmost point in the contiguous US
	maxLon float64 = -68.93  // Easternmost point in the contiguous US

)
