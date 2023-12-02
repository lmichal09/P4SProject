package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
)



/* The life cycle of the lanternfly is as follows:
Eggs can be found on any outdoor surface from October through June.

Egg: October- June ** 
Nymph1: May-July **
Nypmh2 : July- September
Adult: July- December
Adult lays eggs (fecundity): September - November

Data of monthly temperature: October, November, December, May, June, July, August, September
Time: 2015-2021
25 States: "PA" "NJ" "VA" "DE" "MD" "NY" "UT" "MA" "MI" "NC" "WV" "CT" "VT" "OH" "IN" "KY" "DC" "SC" "NM" "AZ" "RI" "OR" "MO" "KS" "ME"
*/


/* Thermal constant for each stage of the lanternfly life cycle	
Table 2 Values of K1, K2, K3 and K4 (degree-days)
K1 		K2 		K3 		K4
39.5 	250 	108.7 	180

*/

func InitializeHabitat(width float64, numberOfFlies, numberOfPredators int) Country {
	country := Country{} // Create the initial country.
	country.width = width
	country.flies = make([]Fly, numberOfFlies)
	country.predators = make([]Predator, numberOfPredators)
	country.states = make([]Quadrant, 25)
	country.states = LoadStates("states.csv")

	// Initialize the flies.
	for i := range country.flies {
		country.flies[i].position = LoadLatternFlyPosition("latternfly.csv")
		
		// Velocity and acceleration are random since no data is available
		country.flies[i].velocity.x = rand.Float64() * 2
		country.flies[i].velocity.y = rand.Float64() * 5
		country.flies[i].acceleration.x = rand.Float64() * rand.Float64() * 2
		country.flies[i].acceleration.y = rand.Float64() * rand.Float64() * 5
		
		//latternfly's stage random from 1-4
		country.flies[i].stage = rand.Intn(4)

		// Initialize the energy of flies
		// stage: egg (1), nymph1 (2), nymph2 (3), adult1 (4), adult 2 (5)
		if country.flies[i].stage == 1 { /
			country.flies[i].energy = rand.Float64() * 39.5
		} else if country.flies[i].stage == 2 { // instar1
			country.flies[i].energy = rand.Float64() * 250
		} else if country.flies[i].stage == 3 {
			country.flies[i].energy = rand.Float64() * 108.7
		} else if country.flies[i].stage == 4 {
			country.flies[i].energy = rand.Float64() * 180
		}
		// when initialize, consider all flies are alive
		country.flies[i].isAlive = true

		// locationID is random from 0-24
		country.flies[i].locationID = rand.Intn(25) // Get data from file, can be changed later
	}

	// Initialize the predators.
	for i := range country.predators {
		country.predators[i].position.x = rand.Float64() * width 
		country.predators[i].position.y = rand.Float64() * width 
		country.predators[i].velocity.x = rand.Float64() * 2
		country.predators[i].velocity.y = rand.Float64() * 5
		country.predators[i].acceleration.x = rand.Float64() * rand.Float64() * 2
		country.predators[i].acceleration.y = rand.Float64() * rand.Float64() * 5
		country.predators[i].PercentEaten = rand.Float64() * 100
	}

	return country
}

func LoadLatternFlyPosition(filePath string) OrderedPair {
	// Read state widths from file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	// File is space-separated

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			log.Fatalf("invalid line format: %v", line)
		}

		// Parse the x and y coordinates
		var x, y float64
		_, err := fmt.Sscanf(parts[0], "%f", &x)
		if err != nil {
			log.Fatalf("invalid x for latternfly: %v", err)
		}
		_, err = fmt.Sscanf(parts[1], "%f", &y)
		if err != nil {
			log.Fatalf("invalid y for latternfly: %v", err)
		}
	}
	return OrderedPair{x, y}
}

// LoadStates reads the state widths from a file and returns a map of state abbreviations to State objects.
// The file is assumed to be comma-separated, with one state per line.
func LoadStates(filePath string) map[string]*State {
	// Initialize the states
	stateAbbreviations := []string{"PA", "NJ", "VA", "DE", "MD", "NY", "UT", "MA", "MI", "NC", "WV", "CT", "VT", "OH", "IN", "KY", "DC", "SC", "NM", "AZ", "RI", "OR", "MO", "KS", "ME"}

	// Create a map to hold state data
	states := make(map[string]*State)
	for _, abbr := range stateAbbreviations {
		states[abbr] = &State{Abbreviation: abbr}
	}

	// Read state widths from file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",") // Assuming the file is comma-separated
		if len(parts) != 2 {
			log.Fatalf("invalid line format: %v", line)
		}

		abbr := parts[0]
		var width float64
		_, err := fmt.Sscanf(parts[1], "%f", &width)
		if err != nil {
			log.Fatalf("invalid width for state %s: %v", abbr, err)
		}

		if state, ok := states[abbr]; ok {
			state.Width = width
		} else {
			log.Printf("Warning: State '%s' not found in the initial list", abbr)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return states
}





func InitializePreyPredatorModel() {

}
