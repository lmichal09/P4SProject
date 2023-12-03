package main

import (
	"math/rand"
	"time"
)

/* The life cycle of the lanternfly is as follows:
Eggs can be found on any outdoor surface from October through June.
May-June: Eggs hatch into 1st instar nymphs.
May-July: 1st instar nymphs molt into 2nd instar nymphs, 2nd instar nymphs molt into 3rd instar nymphs.
july-September: 3rd instar nymphs molt into 4th instar nymphs.
July-December: 4th instar nymphs molt into adults.
September-November: Adults lay eggs.
Oct-June: Eggs overwinter.

Time: 2021
25 States: "PA" "NJ" "VA" "DE" "MD" "NY" "UT" "MA" "MI" "NC" "WV" "CT" "VT" "OH" "IN" "KY" "DC" "SC" "NM" "AZ" "RI" "OR" "MO" "KS" "ME"

The life cycle of the lanternfly is as follows: 6 stages 0-5, stage 7: Dead
*/

// width and population need to be read from file

func InitializeCountry(width float64, population int) Country {

	country := Country{} // Create the initial country.
	country.width = width
	country.flies = make([]Fly, population)
	country.population = population

	// Read Weather data and store in a map
	MayJuneWeather := make(map[string]OrderedPair)

	// Load from subfolder of Data
	MayJuneWeather = LoadWeatherData("Data/Hatch_May-Jun")

	// Initialize the flies.
	for i := range country.flies {
		country.flies[i].position = LoadLatternFlyPosition("latternfly.csv")

		// Velocity and acceleration are random since no data is available
		country.flies[i].velocity.x = rand.Float64() * 2
		country.flies[i].velocity.y = rand.Float64() * 5
		country.flies[i].acceleration.x = rand.Float64() * rand.Float64() * 2
		country.flies[i].acceleration.y = rand.Float64() * rand.Float64() * 5

		//latternfly's stage consider as 0 : egg state
		country.flies[i].stage = 0

		// Set Energy to every fly
		// Latternfly's energy is random from its state and weather data from May to June
		// random from minTemp of the state to maxTemp of the state
		// States are sorted by alphabetical order
		// AZ, CT, DC, DE, IN, KS, KY, MA, MD, ME, MI, MO, NC, NJ, NM, NY, OH, OR, PA, RI, SC, UT, VA, VT, WV

		if country.flies[i].state == "AZ" {
			country.flies[i].energy += randomInRange(MayJuneWeather["AZ"].Y, MayJuneWeather["AZ"].X)
		} else if country.flies[i].state == "CT" {
			country.flies[i].energy += randomInRange(MayJuneWeather["CT"].Y, MayJuneWeather["CT"].X)
		} else if country.flies[i].state == "DC" {
			country.flies[i].energy += randomInRange(MayJuneWeather["DC"].Y, MayJuneWeather["DC"].X)
		} else if country.flies[i].state == "DE" {
			country.flies[i].energy += randomInRange(MayJuneWeather["DE"].Y, MayJuneWeather["DE"].X)
		} else if country.flies[i].state == "IN" {
			country.flies[i].energy += randomInRange(MayJuneWeather["IN"].Y, MayJuneWeather["IN"].X)
		} else if country.flies[i].state == "KS" {
			country.flies[i].energy += randomInRange(MayJuneWeather["KS"].Y, MayJuneWeather["KS"].X)
		} else if country.flies[i].state == "KY" {
			country.flies[i].energy += randomInRange(MayJuneWeather["KY"].Y, MayJuneWeather["KY"].X)
		} else if country.flies[i].state == "MA" {
			country.flies[i].energy += randomInRange(MayJuneWeather["MA"].Y, MayJuneWeather["MA"].X)
		} else if country.flies[i].state == "MD" {
			country.flies[i].energy += randomInRange(MayJuneWeather["MD"].Y, MayJuneWeather["MD"].X)
		} else if country.flies[i].state == "ME" {
			country.flies[i].energy += randomInRange(MayJuneWeather["ME"].Y, MayJuneWeather["ME"].X)
		} else if country.flies[i].state == "MI" {
			country.flies[i].energy += randomInRange(MayJuneWeather["MI"].Y, MayJuneWeather["MI"].X)
		} else if country.flies[i].state == "MO" {
			country.flies[i].energy += randomInRange(MayJuneWeather["MO"].Y, MayJuneWeather["MO"].X)
		} else if country.flies[i].state == "NC" {
			country.flies[i].energy += randomInRange(MayJuneWeather["NC"].Y, MayJuneWeather["NC"].X)
		} else if country.flies[i].state == "NJ" {
			country.flies[i].energy += randomInRange(MayJuneWeather["NJ"].Y, MayJuneWeather["NJ"].X)
		} else if country.flies[i].state == "NM" {
			country.flies[i].energy += randomInRange(MayJuneWeather["NM"].Y, MayJuneWeather["NM"].X)
		} else if country.flies[i].state == "NY" {
			country.flies[i].energy += randomInRange(MayJuneWeather["NY"].Y, MayJuneWeather["NY"].X)
		} else if country.flies[i].state == "OH" {
			country.flies[i].energy += randomInRange(MayJuneWeather["OH"].Y, MayJuneWeather["OH"].X)
		} else if country.flies[i].state == "OR" {
			country.flies[i].energy += randomInRange(MayJuneWeather["OR"].Y, MayJuneWeather["OR"].X)
		} else if country.flies[i].state == "PA" {
			country.flies[i].energy += randomInRange(MayJuneWeather["PA"].Y, MayJuneWeather["PA"].X)
		} else if country.flies[i].state == "RI" {
			country.flies[i].energy += randomInRange(MayJuneWeather["RI"].Y, MayJuneWeather["RI"].X)
		} else if country.flies[i].state == "SC" {
			country.flies[i].energy += randomInRange(MayJuneWeather["SC"].Y, MayJuneWeather["SC"].X)
		} else if country.flies[i].state == "UT" {
			country.flies[i].energy += randomInRange(MayJuneWeather["UT"].Y, MayJuneWeather["UT"].X)
		} else if country.flies[i].state == "VA" {
			country.flies[i].energy += randomInRange(MayJuneWeather["VA"].Y, MayJuneWeather["VA"].X)
		} else if country.flies[i].state == "VT" {
			country.flies[i].energy += randomInRange(MayJuneWeather["VT"].Y, MayJuneWeather["VT"].X)
		} else if country.flies[i].state == "WV" {
			country.flies[i].energy += randomInRange(MayJuneWeather["WV"].Y, MayJuneWeather["WV"].X)
		}
	}

	//randomly choose 10% of the flies to be alive
	for i := 0; i < population/10; i++ {
		country.flies[rand.Intn(population)].isAlive = true
	}
	return country
}

func randomInRange(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano()) // Initialize the random number generator
	return min + rand.Float64()*(max-min)
}
