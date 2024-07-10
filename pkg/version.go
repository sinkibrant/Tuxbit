package pkg

import (
	"flag"
	"os"

	"github.com/aquasecurity/table"
)

var Version bool

func VersionFlag() {
	flag.BoolVar(&Version, "version", false, "Show the version number and exit")

}

func VersionAction() {
	tableVersion := table.New(os.Stdout)

	tableVersion.SetLineStyle(94)
	tableVersion.SetPadding(2)

	tableVersion.AddRow("\033[33m"+"VERSION"+"\033[0m", "1.0.0")

	tableVersion.Render()
}
