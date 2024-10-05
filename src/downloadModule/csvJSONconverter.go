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
	_, err = reader.Read() // This is because the CSV that walletexplorer gives has some text in the first line

	reader.FieldsPerRecord = 7
	if err != nil {
		log.Fatal(err)
	}

	csvRecords, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	csvRecords = csvRecords[1:] // This is done to remove the headers

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

	dataprep.GetTxs(rows, f)
}
