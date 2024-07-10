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

type jsonBlock struct {
	ID                string  `json:"id"`
	Height            int     `json:"height"`
	Version           int     `json:"version"`
	Timestamp         int64   `json:"timestamp"`
	TxCount           int     `json:"tx_count"`
	Size              int     `json:"size"`
	Weight            int     `json:"weight"`
	MerkleRoot        string  `json:"merkle_root"`
	Nextblockhash     string  `json:"next_best"`
	Previousblockhash string  `json:"previousblockhash"`
	Mediantime        int     `json:"mediantime"`
	Nonce             int64   `json:"nonce"`
	Bits              int     `json:"bits"`
	Difficulty        float64 `json:"difficulty"`
}

var Block string

func BlockFlag() {
	flag.StringVar(&Block, "block", "", "Return information about a block using its hash or height")
	flag.StringVar(&Block, "b", "", "Return information about a block using its hash or height")
}

func BlockAction(block *string) {
	var jsonBlock jsonBlock 

	blockHeight, err := strconv.Atoi(*block)
	if err != nil {

		apiResponseProcessed = requestGet(fmt.Sprintf("https://mempool.space/api/block/%s", *block))
		err := json.Unmarshal(apiResponseProcessed, &jsonBlock)
		if err != nil {
			log.Fatal(tableError(&apiResponseProcessed, err))
		}

		apiResponseProcessed = requestGet(fmt.Sprintf("https://mempool.space/api/block/%s/status", *block))
		err = json.Unmarshal(apiResponseProcessed, &jsonBlock)
		if err != nil {
			log.Fatal(tableError(&apiResponseProcessed, err))
		}

		tableBlock(&jsonBlock)
		return
	}

	apiResponseProcessed = requestGet(fmt.Sprintf("https://mempool.space/api/block-height/%d", blockHeight))
	apiResponseProcessed = requestGet(fmt.Sprintf("https://mempool.space/api/block/%s", apiResponseProcessed))
	err = json.Unmarshal(apiResponseProcessed, &jsonBlock)
	if err != nil {
		log.Fatal(tableError(&apiResponseProcessed, err))
	}

	apiResponseProcessed = requestGet(fmt.Sprintf("https://mempool.space/api/block/%s/status", jsonBlock.ID))
	err = json.Unmarshal(apiResponseProcessed, &jsonBlock)
	if err != nil {
		log.Fatal(tableError(&apiResponseProcessed, err))
	}

	tableBlock(&jsonBlock)
}

func tableBlock(jsonBlock *jsonBlock) {
	table := table.New(os.Stdout)
	table.SetLineStyle(94)
	table.SetPadding(2)

	table.AddRow("\033[33m"+"Hash"+"\033[0m", jsonBlock.ID)
	table.AddRow("\033[33m"+"Height"+"\033[0m", strconv.Itoa(jsonBlock.Height))
	table.AddRow("\033[33m"+"Version"+"\033[0m", fmt.Sprintf("0x%x", jsonBlock.Version))
	table.AddRow("\033[33m"+"Transaction Count"+"\033[0m", strconv.Itoa(jsonBlock.TxCount))
	table.AddRow("\033[33m"+"Timestamp"+"\033[0m", fmt.Sprint(time.Unix(jsonBlock.Timestamp, 0).Local()))

	if jsonBlock.Size >= 1000000 {
		size := float32(jsonBlock.Size) / 1000000
		table.AddRow("\033[33m"+"Size"+"\033[0m", fmt.Sprintf("%.2fMB", size))
	} else if jsonBlock.Size >= 1000 && jsonBlock.Size < 1000000 {
		size := float32(jsonBlock.Size) / 1000
		table.AddRow("\033[33m"+"Size"+"\033[0m", fmt.Sprintf("%.2fKB", size))
	} else {
		table.AddRow("\033[33m"+"Size"+"\033[0m", fmt.Sprintf("%.dB", jsonBlock.Size))
	}

	if jsonBlock.Weight >= 1000000 {
		weight := float32(jsonBlock.Weight) / 1000000
		table.AddRow("\033[33m"+"Weight"+"\033[0m", fmt.Sprintf("%.2fMWU", weight))
	} else if jsonBlock.Weight >= 1000 && jsonBlock.Weight < 1000000 {
		weight := float32(jsonBlock.Weight) / 1000
		table.AddRow("\033[33m"+"Weight"+"\033[0m", fmt.Sprintf("%.2fKWU", weight))
	} else {
		table.AddRow("\033[33m"+"Weight"+"\033[0m", fmt.Sprintf("%dWU", jsonBlock.Weight))
	}

	if jsonBlock.Nextblockhash != "" {
		table.AddRow("\033[33m"+"Next Bloco Hash"+"\033[0m", jsonBlock.Nextblockhash)
	}

	table.AddRow("\033[33m"+"Previous Block Hash"+"\033[0m", jsonBlock.Previousblockhash)

	table.AddRow("\033[33m"+"Merkle Root"+"\033[0m", jsonBlock.MerkleRoot)
	table.AddRow("\033[33m"+"Nonce"+"\033[0m", fmt.Sprintf("0x%x", jsonBlock.Nonce))
	table.AddRow("\033[33m"+"Bits"+"\033[0m", fmt.Sprintf("0x%x", jsonBlock.Bits))
	table.AddRow("\033[33m"+"Difficulty"+"\033[0m", fmt.Sprintf("%.1f", jsonBlock.Difficulty))

	defer table.Render()
}
