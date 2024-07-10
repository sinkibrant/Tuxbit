package pkg

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aquasecurity/table"
)

type jsonAddress struct {
	Address    string `json:"address"`
	ChainStats struct {
		FundedTxoCount int `json:"funded_txo_count"`
		FundedTxoSum   int `json:"funded_txo_sum"`
		SpentTxoCount  int `json:"spent_txo_count"`
		SpentTxoSum    int `json:"spent_txo_sum"`
		TxCount        int `json:"tx_count"`
	} `json:"chain_stats"`
	Isvalid bool `json:"isvalid"`
}

var Address string
var transactions bool

func AddressFlags() {
	flag.StringVar(&Address, "address", "", "Return information about a Bitcoin address")
	flag.StringVar(&Address, "a", "", "Return information about a Bitcoin address")

	flag.BoolVar(&transactions, "transactions", false, "Return the latest 50 transactions associated with the address")
	flag.BoolVar(&transactions, "t", false, "Return the latest 50 transactions associated with the address")
}

func AddressAction(address *string) {
	var jsonAddress jsonAddress

	apiResponseProcessed = requestGet(fmt.Sprintf("https://mempool.space/api/address/%s", *address))
	err := json.Unmarshal(apiResponseProcessed, &jsonAddress)
	if err != nil {
		log.Fatal(tableError(&apiResponseProcessed, err))
	}

	apiResponseProcessed = requestGet(fmt.Sprintf("https://mempool.space/api/v1/validate-address/%s", *address))
	err = json.Unmarshal(apiResponseProcessed, &jsonAddress)
	if err != nil {
		log.Fatal(tableError(&apiResponseProcessed, err))
	}

	if transactions {
		var jsonTx []jsonTx

		apiResponseProcessed := requestGet(fmt.Sprintf("https://mempool.space/api/address/%s/txs", *address))
		err := json.Unmarshal(apiResponseProcessed, &jsonTx)
		if err != nil {
			log.Fatal(tableError(&apiResponseProcessed, err))
		}

		tableAddressTransactions(&jsonAddress, &jsonTx)
		return
	}

	tableAddress(&jsonAddress)
}

func tableAddressTransactions(jsonAddress *jsonAddress, jsonTx *[]jsonTx) {
	tableAddress(jsonAddress)

	for j := range *jsonTx {
		tableTx(&(*jsonTx)[j])
	}

}

func tableAddress(jsonAddress *jsonAddress) {
	tableAddress := table.New(os.Stdout)

	tableAddress.SetHeaderStyle(1)
	tableAddress.SetLineStyle(94)
	tableAddress.SetPadding(2)

	tableAddress.SetHeaders("\033[33m"+"Address"+"\033[0m", jsonAddress.Address)

	tableAddress.AddRow("\033[33m"+"Confirmed Balance"+"\033[0m", fmt.Sprintf("%.8f BTC", float64(jsonAddress.ChainStats.FundedTxoSum-jsonAddress.ChainStats.SpentTxoSum)/100000000))
	tableAddress.AddRow("\033[33m"+"Total Received"+"\033[0m", fmt.Sprintf("%.8f BTC", float64(jsonAddress.ChainStats.FundedTxoSum)/100000000))
	tableAddress.AddRow("\033[33m"+"Total Transactions"+"\033[0m", strconv.Itoa(jsonAddress.ChainStats.TxCount))

	tableAddress.Render()
}
