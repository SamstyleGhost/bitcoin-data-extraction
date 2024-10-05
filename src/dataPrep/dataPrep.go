package dataprep

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/SamstyleGhost/bitcoin-data-extraction/src/custom"
)

func GetTxs(transactions []custom.TransactionRow, address string) {
	// txCount, startDepth, maxDepth := 0, 0, 3

	var cashFlowTransactions []custom.CashFlowTransaction
	var transactionsMu sync.Mutex
	var wg sync.WaitGroup

	for _, tx := range transactions {
		wg.Add(1)
		go getTransactionThroughID(tx.Transaction, &cashFlowTransactions, &transactionsMu, &wg)
	}

	wg.Wait()

	file, err := os.Create(address + "_tx_lowers.json")
	if err != nil {
		fmt.Println("Error creating, ", err)
		return
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(cashFlowTransactions, "", "    ")
	if err != nil {
		fmt.Println("Error unmarshalling, ", err)
		return
	}

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing, ", err)
		return
	}

	fmt.Println("Done")
}

func getTransactionThroughID(txID string, slice *[]custom.CashFlowTransaction, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Working on calling transaction: %s\n", txID)

	apiURL := "http://www.walletexplorer.com/api/1/tx?txid=" + txID + "&caller=Ghost"
	client := http.Client{
		Timeout: 20 * time.Second, // Set a timeout for the HTTP request
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

	var obj custom.CashFlowTransaction
	err = json.Unmarshal(body, &obj)

	if err != nil {
		fmt.Println(txID, err)
		return
	}

	mu.Lock()
	*slice = append(*slice, obj)
	defer mu.Unlock()

	fmt.Printf("Done %s\n", txID)
}
