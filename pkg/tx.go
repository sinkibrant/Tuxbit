package pkg

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aquasecurity/table"
)

type jsonTx struct {
	Txid     string `json:"txid"`
	Version  int    `json:"version"`
	Locktime int    `json:"locktime"`
	Vin      []struct {
		Txid    string `json:"txid"`
		Vout    int    `json:"vout"`
		Prevout struct {
			Scriptpubkey        string `json:"scriptpubkey"`
			ScriptpubkeyAsm     string `json:"scriptpubkey_asm"`
			ScriptpubkeyType    string `json:"scriptpubkey_type"`
			ScriptpubkeyAddress string `json:"scriptpubkey_address"`
			Value               int    `json:"value"`
		} `json:"prevout"`
		Scriptsig    string   `json:"scriptsig"`
		ScriptsigAsm string   `json:"scriptsig_asm"`
		Witness      []string `json:"witness"`
		IsCoinbase   bool     `json:"is_coinbase"`
		Sequence     int      `json:"sequence"`
	} `json:"vin"`
	Vout []struct {
		Scriptpubkey        string `json:"scriptpubkey"`
		ScriptpubkeyAsm     string `json:"scriptpubkey_asm"`
		ScriptpubkeyType    string `json:"scriptpubkey_type"`
		ScriptpubkeyAddress string `json:"scriptpubkey_address"`
		Value               int    `json:"value"`
	} `json:"vout"`
	Size   int `json:"size"`
	Weight int `json:"weight"`
	Sigops int `json:"sigops"`
	Fee    int `json:"fee"`
	Status struct {
		Confirmed   bool   `json:"confirmed"`
		BlockHeight int    `json:"block_height"`
		BlockHash   string `json:"block_hash"`
		BlockTime   int64  `json:"block_time"`
	} `json:"status"`
}

var Tx string
var txVerbose bool

func TxFlag() {
	flag.StringVar(&Tx, "tx", "", "Return simplified information about a transaction")

	flag.BoolVar(&txVerbose, "verbose", false, "Return detailed information about the associated transaction")
	flag.BoolVar(&txVerbose, "v", false, "Return detailed information about the associated transaction")
}

func TxAction(tx *string) {
	var jsonTx jsonTx

	apiResponseProcessed = requestGet(fmt.Sprintf("https://mempool.space/api/tx/%s", *tx))
	err := json.Unmarshal(apiResponseProcessed, &jsonTx)
	if err != nil {
		log.Fatal(tableError(&apiResponseProcessed, err))
	}

	if txVerbose {
		tableTxVerbose(&jsonTx)
		return
	}

	tableTx(&jsonTx)
}

func tableTx(jsonTx *jsonTx) {
	var txValue float64
	for i := range jsonTx.Vout {
		txValue += float64(jsonTx.Vout[i].Value)
	}
	txValue /= 100000000

	var statusTX string
	if jsonTx.Status.Confirmed {
		statusTX = "Confirmed"
	} else {
		statusTX = "Not confirmed"
	}

	tableTransaction := table.New(os.Stdout)

	tableTransaction.SetHeaderStyle(1)
	tableTransaction.SetLineStyle(94)
	tableTransaction.SetPadding(2)

	tableTransaction.SetHeaders("\033[33m"+"Transaction"+"\033[0m", "\033[33m"+"Detalhes"+"\033[0m")

	tableTransaction.AddRow("\033[33m"+"TXID"+"\033[0m", jsonTx.Txid)
	tableTransaction.AddRow("\033[33m"+"TX Value"+"\033[0m", fmt.Sprintf("%.8f", txValue))
	tableTransaction.AddRow("\033[33m"+"Check"+"\033[0m", statusTX)
	if statusTX != "Not confirmed" {
	tableTransaction.AddRow("\033[33m"+"Block Height"+"\033[0m", strconv.Itoa(jsonTx.Status.BlockHeight))
	tableTransaction.AddRow("\033[33m"+"Block Hash"+"\033[0m", jsonTx.Status.BlockHash)
	tableTransaction.AddRow("\033[33m"+"Block Time"+"\033[0m", fmt.Sprint(time.Unix(jsonTx.Status.BlockTime, 0)))
	}

	tableTransaction.Render()

}

func tableTxVerbose(jsonTx *jsonTx) {

	var txValue float64
	for i := range jsonTx.Vout {
		txValue += float64(jsonTx.Vout[i].Value)
	}
	txValue /= 100000000

	var fee float32
	if jsonTx.Fee >= 1000000 {
		fee = float32(jsonTx.Fee) / 1000000
	} else if jsonTx.Fee >= 1000 && jsonTx.Fee < 1000000 {
		fee = float32(jsonTx.Fee) / 1000
	} else {
		fee = float32(jsonTx.Fee)
	}

	var statusTX string
	if jsonTx.Status.Confirmed {
		statusTX = "Confirmed"
	} else {
		statusTX = "Not confirmed"
	}

	tableTransaction := table.New(os.Stdout)

	tableTransaction.SetHeaderStyle(1)
	tableTransaction.SetLineStyle(94)
	tableTransaction.SetPadding(2)

	tableTransaction.SetHeaders("\033[33m"+"Transaction"+"\033[0m", "\033[33m"+"Detalhes"+"\033[0m")

	tableTransaction.AddRow("\033[33m"+"TXID"+"\033[0m", jsonTx.Txid)
	tableTransaction.AddRow("\033[33m"+"TX Value"+"\033[0m", fmt.Sprintf("%.8f", txValue))
	tableTransaction.AddRow("\033[33m"+"Version"+"\033[0m", strconv.Itoa(jsonTx.Version))
	tableTransaction.AddRow("\033[33m"+"LockTime"+"\033[0m", strconv.Itoa(jsonTx.Locktime))

	if jsonTx.Size >= 1000000 {
		size := float32(jsonTx.Size) / 1000000
		tableTransaction.AddRow("\033[33m"+"Size"+"\033[0m", fmt.Sprintf("%.2fMB", size))
	} else if jsonTx.Size >= 1000 && jsonTx.Size < 1000000 {
		size := float32(jsonTx.Size) / 1000
		tableTransaction.AddRow("\033[33m"+"Size"+"\033[0m", fmt.Sprintf("%.2fKB", size))
	} else {
		tableTransaction.AddRow("\033[33m"+"Size"+"\033[0m", fmt.Sprintf("%.dB", jsonTx.Size))
	}

	if jsonTx.Weight >= 1000000 {
		weight := float32(jsonTx.Weight) / 1000000
		tableTransaction.AddRow("\033[33m"+"Weight"+"\033[0m", fmt.Sprintf("%.2fMWU", weight))
	} else if jsonTx.Weight >= 1000 && jsonTx.Weight < 1000000 {
		weight := float32(jsonTx.Weight) / 1000
		tableTransaction.AddRow("\033[33m"+"Weight"+"\033[0m", fmt.Sprintf("%.2fKWU", weight))
	} else {
		tableTransaction.AddRow("\033[33m"+"Weight"+"\033[0m", fmt.Sprintf("%dWU", jsonTx.Weight))
	}

	tableTransaction.AddRow("\033[33m"+"Signature Operation"+"\033[0m", strconv.Itoa(jsonTx.Sigops))
	tableTransaction.AddRow("\033[33m"+"Fee"+"\033[0m", fmt.Sprintf("%.3f", fee))
	tableTransaction.AddRow("\033[33m"+"Check"+"\033[0m", statusTX)
	if statusTX != "Not confirmed" {
	tableTransaction.AddRow("\033[33m"+"Block Height"+"\033[0m", strconv.Itoa(jsonTx.Status.BlockHeight))
	tableTransaction.AddRow("\033[33m"+"Block Hash"+"\033[0m", jsonTx.Status.BlockHash)
	tableTransaction.AddRow("\033[33m"+"Block Time"+"\033[0m", fmt.Sprint(time.Unix(jsonTx.Status.BlockTime, 0)))
	}

	tableTransaction.Render()

	//------------------------------

	tableInputs := table.New(os.Stdout)

	tableInputs.SetHeaderStyle(1)
	tableInputs.SetLineStyle(94)
	tableInputs.SetPadding(2)

	tableInputs.SetHeaders("\033[33m"+"Inputs"+"\033[0m", "\033[33m"+"Detalhes"+"\033[0m")

	for i := range jsonTx.Vin {
		tableInputs.AddRow("\033[33m"+"TXID"+"\033[0m", jsonTx.Vin[i].Txid)
		tableInputs.AddRow("\033[33m"+"Vout"+"\033[0m", strconv.Itoa(jsonTx.Vin[i].Vout))
		tableInputs.AddRow("\033[33m"+"Address"+"\033[0m", jsonTx.Vin[i].Prevout.ScriptpubkeyAddress)
		tableInputs.AddRow("\033[33m"+"Value"+"\033[0m", fmt.Sprintf("%.8f", float64(jsonTx.Vin[i].Prevout.Value)/100000000))
		tableInputs.AddRow("\033[33m"+"Sequence"+"\033[0m", strconv.Itoa(jsonTx.Vin[i].Sequence))
		if i+1 != len(jsonTx.Vin) {
			tableInputs.AddRow()
		}
	}

	tableInputs.Render()

	//------------------------------

	tableOutputs := table.New(os.Stdout)

	tableOutputs.SetHeaderStyle(1)
	tableOutputs.SetLineStyle(94)
	tableOutputs.SetPadding(2)

	tableOutputs.SetHeaders("\033[33m"+"Outputs"+"\033[0m", "\033[33m"+"Detalhes"+"\033[0m")

	for i := range jsonTx.Vout {
		tableOutputs.AddRow("\033[33m"+"Address"+"\033[0m", jsonTx.Vout[i].ScriptpubkeyAddress)
		tableOutputs.AddRow("\033[33m"+"Value"+"\033[0m", fmt.Sprintf("%.8f", float64(jsonTx.Vout[i].Value)/100000000))
		if i != len(jsonTx.Vout) - 1 {
			tableOutputs.AddRow()
		}
	}

	tableOutputs.Render()

}
