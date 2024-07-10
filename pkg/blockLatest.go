package pkg

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
)

var BlockLatest bool

func BlockLatestFlag() {
	flag.BoolVar(&BlockLatest, "block-latest", false, "Return information about the latest block")
	flag.BoolVar(&BlockLatest, "bl", false, "Consulta o bloco mais recente")
}

func BlockLatestAction() {
	var jsonBlock jsonBlock

	apiResponseProcessed = requestGet("https://mempool.space/api/blocks/tip/hash")
	apiResponseProcessed = requestGet(fmt.Sprintf("https://mempool.space/api/block/%s", apiResponseProcessed))
	err := json.Unmarshal(apiResponseProcessed, &jsonBlock)
	if err != nil {
		log.Fatal(tableError(&apiResponseProcessed, err))
	}
	
	tableBlock(&jsonBlock)
}
