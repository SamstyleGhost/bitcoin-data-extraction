package dataprep

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/SamstyleGhost/bitcoin-data-extraction/src/custom"
)

/*
Since I am writing this in GO for the purpose of using concurrency, I will need to use goroutines and also a channel and a mutex. Here is my thoughts:
? My end goal is to create a JSON object that would have the cash-flow network (I will have to see if I need to create a JSON file to store the data, because thats what Adam's code does)
* I will need to have a mutex on the object so that there is no read-write issues while writing
? How do I limit how many goroutines are spawned
* Feel like I need to use a buffer which would also have a mutex lock on it. So, whenever a goroutine that has completed its work gives out a JSON object, the forked routine would be joined back and then forked again for a different address
* This is the Producer-Consumer problem: I have to make sure that the Producer (making the API calls and parsing JSON) writes to a buffer that is being read by the Consumer (writing data to the JSON file)
*/

func GetTxs(transactions []custom.TransactionRow) {
	// txCount, startDepth, maxDepth := 0, 0, 3

	// var inwardTransactions []custom.TransactionRow
	// var outwardTransactions []custom.TransactionRow

	// for _, tx := range transactions {
	// 	if tx.SentAmount == 0 {
	// 		// This means that the tx is inward transaction
	// 	} else if tx.ReceivedAmount == 0 {
	// 		// This means that the tx is outward transaction
	// 	}
	// }

	var wg sync.WaitGroup

	for _, tx := range transactions {
		wg.Add(1)
		go getTransactionThroughID(tx.Transaction, &wg)
	}

	wg.Wait()

}

func getTransactionThroughID(txID string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Working on calling transaction: %s\n", txID)
	apiURL := "http://www.walletexplorer.com/api/1/tx?txid=" + txID + "&caller=Ghost"
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}

	switch v := result.(type) {
	case map[string]interface{}:
		for key, value := range v {
			fmt.Println(key, value)
		}
	case []interface{}:
		for _, value := range v {
			fmt.Println(value)
		}
	default:
		fmt.Println(v)
	}
	fmt.Println("--------")
}
