package dataprep

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/SamstyleGhost/bitcoin-data-extraction/src/custom"
)

/*
? My end goal is to create a JSON object that would have the cash-flow network (I will have to see if I need to create a JSON file to store the data, because thats what Adam's code does)
* I will need to have a mutex on the object so that there is no read-write issues while writing
? How do I limit how many goroutines are spawned
* Feel like I need to use a buffer which would also have a mutex lock on it. So, whenever a goroutine that has completed its work gives out a JSON object, the forked routine would be joined back and then forked again for a different address
* This is the Producer-Consumer problem: I have to make sure that the Producer (making the API calls and parsing JSON) writes to a buffer that is being read by the Consumer (writing data to the JSON file)
*/

func GetTxs(transactions []custom.TransactionRow) {
	// txCount, startDepth, maxDepth := 0, 0, 3

	var inwardTransactions [][]byte
	var outwardTransactions [][]byte
	var inmu sync.Mutex
	var outmu sync.Mutex
	var wg sync.WaitGroup

	// var dataSlice []custom.CashFlowTransaction
	// dataChan := make(chan custom.CashFlowTransaction, 1000)

	for _, tx := range transactions {
		// dataChan <- custom.CashFlowTransaction{}
		wg.Add(1)
		if tx.SentTo == "(fee)" {
			// This means that the tx is inward transaction
			go getTransactionThroughID(tx.Transaction, &inwardTransactions, &inmu, &wg)
		} else {
			// This means that the tx is outward transaction
			go getTransactionThroughID(tx.Transaction, &outwardTransactions, &outmu, &wg)
		}
	}

	wg.Wait()

	fmt.Println("Inward")
	for _, tx := range inwardTransactions {
		fmt.Println(string(tx))
	}

}

func getTransactionThroughID(txID string, slice *[][]byte, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Working on calling transaction: %s\n", txID)

	apiURL := "http://www.walletexplorer.com/api/1/tx?txid=" + txID + "&caller=Ghost"
	client := http.Client{
		Timeout: 10 * time.Second, // Set a timeout for the HTTP request
	}
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(txID, err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error fetching transaction %s: received status %d\n", txID, resp.StatusCode)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(txID, err)
		return
	}

	// var result interface{}
	var result map[string]interface{} // This helps in getting data that we do not know the structure of beforehand
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(txID, err)
		return
	}
	var obj custom.CashFlowTransaction

	if val, ok := result["block_pos"].(float64); ok {
		obj.BlockPos = int32(val)
	}

	jsonValue, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(txID, err)
		return
	}

	mu.Lock()
	*slice = append(*slice, jsonValue)
	defer mu.Unlock()

	fmt.Printf("Done %s\n", txID)
}
