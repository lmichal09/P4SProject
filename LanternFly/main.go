package main

import (
	"fmt"
)

func main() {
	fmt.Println("Lantern Flies simulation!")

	// step 1: reading input from a single file.

	filename := "Data/lydetext.txt"
	allData := ReadSamplesFromDirectory(filename)

	//step 2: reading input from a directory

	for sampleName, data := range allData {
		csvFilename := sampleName + ".csv"
		err := WriteToFile(csvFilename, data)
		if err != nil {
			fmt.Printf("Error writing to file %s: %v\n", csvFilename, err)
		} else {
			fmt.Printf("Data written to %s\n", csvFilename)
		}
	}

	fmt.Println("Success! Now we are ready to do something cool with our data.")
}
