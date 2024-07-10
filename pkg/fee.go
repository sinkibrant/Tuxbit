package pkg

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/aquasecurity/table"
)

type jsonFee struct {
	FastestFee  int `json:"fastestFee"`
	HalfHourFee int `json:"halfHourFee"`
	HourFee     int `json:"hourFee"`
	EconomyFee  int `json:"economyFee"`
	MinimumFee  int `json:"minimumFee"`
}

var Fee bool

func FeeFlag() {
	flag.BoolVar(&Fee, "fee", false, "Return the currently suggested transaction fees")
	flag.BoolVar(&Fee, "f", false, "Return the currently suggested transaction fees")
}

func FeeAction() {
	var jsonFee jsonFee

	apiResponseProcessed = requestGet("https://mempool.space/api/v1/fees/recommended")
	err := json.Unmarshal(apiResponseProcessed, &jsonFee)
	if err != nil {
		log.Fatal(tableError(&apiResponseProcessed, err))
	}

	tableFee(&jsonFee)
}

func tableFee(jsonFee *jsonFee) {
	tableFee := table.New(os.Stdout)

	tableFee.SetLineStyle(94)
	tableFee.SetPadding(2)

	tableFee.AddRow("\033[33m"+"High Priority"+"\033[0m", strconv.Itoa(jsonFee.FastestFee))
	tableFee.AddRow("\033[33m"+"Medium Priority"+"\033[0m", strconv.Itoa(jsonFee.HalfHourFee))
	tableFee.AddRow("\033[33m"+"Low Priority"+"\033[0m", strconv.Itoa(jsonFee.HourFee))
	tableFee.AddRow("\033[33m"+"No Priority"+"\033[0m", strconv.Itoa(jsonFee.EconomyFee))

	tableFee.Render()
}
