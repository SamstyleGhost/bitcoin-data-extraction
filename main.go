package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	downloadmodule "github.com/SamstyleGhost/bitcoin-data-extraction/src/downloadModule"
)

func CheckCSVGet() {
	resp, err := http.Get("https://www.walletexplorer.com/wallet/d394a6a98aabeeae?format=csv")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	temp, err := os.CreateTemp("./tempfiles", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(temp.Name())

	fmt.Printf("%s", body)
}

func main() {
	downloadmodule.CSVReader("13AM4VW2dhxYgXeQepoHkHSQuy6NgaEb94")
}
