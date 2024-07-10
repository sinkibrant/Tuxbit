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

type jsonDifficulty struct {
	ProgressPercent       float64 `json:"progressPercent"`
	DifficultyChange      float64 `json:"difficultyChange"`
	EstimatedRetargetDate int64   `json:"estimatedRetargetDate"`
	RemainingBlocks       int     `json:"remainingBlocks"`
	PreviousRetarget      float64 `json:"previousRetarget"`
	PreviousTime          int64     `json:"previousTime"`
	NextRetargetHeight    int     `json:"nextRetargetHeight"`
}

var Difficulty bool

func DifficultyFlag() {
	flag.BoolVar(&Difficulty, "difficulty", false, "Return information about the network difficulty adjustment")
	flag.BoolVar(&Difficulty, "d", false, "Return information about the network difficulty adjustment")
}

func DifficultyAction() {
	var jsonDifficulty jsonDifficulty

	apiResponseProcessed = requestGet("https://mempool.space/api/v1/difficulty-adjustment")
	err := json.Unmarshal(apiResponseProcessed, &jsonDifficulty)
	if err != nil {
		log.Fatal(tableError(&apiResponseProcessed, err))
	}

	tableDifficulty(&jsonDifficulty)
}

func tableDifficulty(jsonDifficulty *jsonDifficulty) {

	tableDifficulty := table.New(os.Stdout)
	
	tableDifficulty.SetLineStyle(94)
	tableDifficulty.SetPadding(2)

	tableDifficulty.AddRow("\033[33m"+"Difficulty Progress"+"\033[0m", fmt.Sprintf("%.3f%%", jsonDifficulty.ProgressPercent))
	tableDifficulty.AddRow("\033[33m"+"Variation"+"\033[0m", fmt.Sprintf("%.4f%%", jsonDifficulty.DifficultyChange))
	tableDifficulty.AddRow("\033[33m"+"Adjustment Block"+"\033[0m", strconv.Itoa(jsonDifficulty.NextRetargetHeight))
	tableDifficulty.AddRow("\033[33m"+"Blocks Remaining"+"\033[0m", strconv.Itoa((jsonDifficulty.RemainingBlocks)))
	tableDifficulty.AddRow("\033[33m"+"Estimated Adjustment Date"+"\033[0m", fmt.Sprint(time.UnixMilli(jsonDifficulty.EstimatedRetargetDate).Local()))
	tableDifficulty.AddRow("\033[33m"+"Previous Adjustment"+"\033[0m", fmt.Sprintf("%.4f%%", jsonDifficulty.PreviousRetarget))
	tableDifficulty.AddRow("\033[33m"+"Last Adjustment"+"\033[0m", fmt.Sprint(time.Unix(jsonDifficulty.PreviousTime, 0).Local()))

	tableDifficulty.Render()
}
