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
