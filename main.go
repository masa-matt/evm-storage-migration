package main

import (
	"flag"
	"fmt"
	"os"

	"evm-storage-migration/migration"

	"github.com/joho/godotenv"
)

var (
	target string
)

func init() {
	flag.CommandLine.Init("migrate", flag.ExitOnError)

	flag.CommandLine.Usage = func() {
		o := flag.CommandLine.Output()
		fmt.Fprintf(o, "\nUsage: %s\n", flag.CommandLine.Name())
		fmt.Fprintf(o, "\nDescription: EVM Storage Migration Tools.\n\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&target, "target", "", "Specify the Contract Name")

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	if flag.NArg() > 0 {
		args := flag.Args()
		switch args[0] {
		case "migrate":
			if err := flag.CommandLine.Parse(args[1:]); err != nil {
				if err != flag.ErrHelp {
					fmt.Fprintf(os.Stderr, "error: %s\n", err)
				}
				os.Exit(1)
			}
			if target == "" {
				flag.CommandLine.Usage()
				os.Exit(1)
			}

			migration.Migrate(target)
			os.Exit(0)
		}
	}

	flag.CommandLine.Usage()
	os.Exit(1)
}
