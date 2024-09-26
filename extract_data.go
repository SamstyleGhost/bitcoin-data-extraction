package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func CSVReader() {
	rf, err := os.Open("./output3.csv")

	if err != nil {
		panic(err)
	}

	defer rf.Close()

	reader := csv.NewReader(rf)
	_, err = reader.Read()

	// ? Here I set the reader's fields value to 7 because that is what I will get
	// ? But, I will have to check if in my trials, there comes a file with multiple lines of the stupid first line
	// ? If there is, I can just run reader.Read in a loop, if the fields value is 1, I will set it to 0 and read again
	reader.FieldsPerRecord = 7

	if err != nil {
		log.Fatal(err)
	}

	_, err = reader.Read() // These are for the header lines. I will see if I need those later
	if err != nil {
		log.Fatal(err)
	}

	// TODO: There are also a couple of rows where the values of addresses are like: Huobi.com (23cvx34...), Polniex.com (00xv23...), etc.
	// Will need to see what to do with those
	// Thinking of just ignoring them when converting to JSON
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(records)
}

func main() {
	CSVReader()
}
