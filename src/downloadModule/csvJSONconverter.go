package downloadmodule

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"github.com/SamstyleGhost/bitcoin-data-extraction/src/custom"
	dataprep "github.com/SamstyleGhost/bitcoin-data-extraction/src/dataPrep"
)

func CSVReader(f string) {
	path := "/home/rohan/Work/bitcoin-data-extraction/data/"

	filename := path + f + "_tx_history.csv"
	rf, err := os.Open(filename)

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

	// TODO: There are also a couple of rows where the values of addresses are like: Huobi.com (23cvx34...), Polniex.com (00xv23...), etc.
	// Will need to see what to do with those
	// Thinking of just ignoring them when converting to JSON

	csvRecords, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// csvRecords = csvRecords[1:]

	var row custom.TransactionRow
	var rows []custom.TransactionRow

	for _, r := range csvRecords {
		row.Date = r[0]
		if r[1] != "" {
			row.ReceivedFrom = r[1]
			row.ReceivedAmount, err = strconv.ParseFloat(r[2], 32)
		}
		if err != nil {
			log.Fatal(err)
		}
		if r[4] != "" {
			row.SentAmount, err = strconv.ParseFloat(r[3], 32)
			if err != nil {
				log.Fatal(err)
			}
			row.SentTo = r[4]
		}
		row.Balance, err = strconv.ParseFloat(r[5], 32)
		if err != nil {
			log.Fatal(err)
		}
		row.Transaction = r[6]

		rows = append(rows, row)
	}

	dataprep.GetTxs(rows)
	// transactionsFile, err := os.Create(path + f + "_tx_history.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// defer transactionsFile.Close()
	// transactionsFile.Write(transactionsJSON)
}
