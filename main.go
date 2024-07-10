package main

import (
	"flag"
	"fmt"
	"tuxbit/pkg"
)

func init() {
	pkg.BlockFlag()
	pkg.BlockLatestFlag()
	pkg.TxFlag()
	pkg.AddressFlags()
	pkg.DifficultyFlag()
	pkg.FeeFlag()
	pkg.VersionFlag()

}

func main() {

	flag.Usage = func() {
		w := flag.CommandLine.Output()

		fmt.Fprintf(w, "usage: tuxbit [-a ADDRESS [-t]] [-b HASH/HEIGHT] [-bl] [-d] [-f] [-h] [-tx TXID [-v]] [--version]\n\n")

		fmt.Fprintf(w, "Command-line interface to explore the Bitcoin ecosystem.\n")
		fmt.Fprintf(w, "--------------------------------------------------------------------------\n")
		fmt.Fprintf(w, "https://gitlab.com/sinkibrant/tuxbit\n\n")

		fmt.Fprintf(w, "optional arguments:\n\n")
		fmt.Fprintf(w, " -a, --address HASH_ADDRESS      Return information about a Bitcoin address\n")
		fmt.Fprintf(w, "    -t, --transactions             Return the latest 50 transactions associated with the address\n\n")

		fmt.Fprintf(w, " -b,  --block HASH/HEIGHT        Return information about a block using its hash or height\n\n")

		fmt.Fprintf(w, " -bl, --block-latest             Return information about the latest block\n\n")

		fmt.Fprintf(w, " -d, --difficulty                Return information about the network difficulty adjustment\n\n")

		fmt.Fprintf(w, " -f, --fee                       Return the currently suggested transaction fees\n\n")

		fmt.Fprintf(w, " -h,  --help                     Show this help message and exit\n\n")

		fmt.Fprintf(w, " -tx HASH_TRANSACTION            Return simplified information about a transaction\n")
		fmt.Fprintf(w, "    -v, --verbose                  Return detailed information about the associated transaction\n\n")

		fmt.Fprintf(w, " --version                       Show the version number and exit\n")

	}

	flag.Parse()

	if pkg.Block != "" {
		pkg.BlockAction(&pkg.Block)
	}
	if pkg.BlockLatest {
		pkg.BlockLatestAction()
	}
	if pkg.Tx != "" {
		pkg.TxAction(&pkg.Tx)
	}
	if pkg.Address != "" {
		pkg.AddressAction(&pkg.Address)
	}
	if pkg.Difficulty {
		pkg.DifficultyAction()
	}
	if pkg.Fee {
		pkg.FeeAction()
	}
	if pkg.Version {
		pkg.VersionAction()
	}
}
