package main

import (
	"flag"
	"fmt"
	"os"

	"evm-storage-migration/migration"

	"github.com/joho/godotenv"
)

var (
	all     = FlagSet{description: "EVM Storage Migration Tools"}
	migrate = FlagSet{description: "Migrate the storage between EVM compatible chains"}
	verify  = FlagSet{description: "Verify the migration results"}

	target string
)

type FlagSet struct {
	fls         *flag.FlagSet
	description string
}

func init() {
	all.fls = flag.CommandLine
	migrate.fls = flag.NewFlagSet("migrate", flag.ExitOnError)
	verify.fls = flag.NewFlagSet("verify", flag.ExitOnError)

	flag.Usage = usage(&all, true)
	migrate.fls.Usage = usage(&migrate, false)
	verify.fls.Usage = usage(&verify, false)

	migrate.fls.StringVar(&target, "target", "", "Specify the Contract Name")
	verify.fls.StringVar(&target, "target", "", "Specify the Contract Name")

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
			parse(migrate.fls, args)
			migration.Migrate(target)
			os.Exit(0)
		case "verify":
			parse(verify.fls, args)
			migration.Verify(target)
			os.Exit(0)
		}
	}

	flag.Usage()
	os.Exit(0)
}

func usage(set *FlagSet, all bool) func() {
	return func() {
		o := set.fls.Output()
		fmt.Fprintf(o, "\nUsage: %s\n", set.fls.Name())
		fmt.Fprintf(o, "\nDescription: %s\n\nOptions:\n", set.description)
		if all {
			fmt.Fprintf(o, "\nCommands:\n")
			fmt.Fprintf(o, "  %s\t%s\n", migrate.fls.Name(), migrate.description)
			fmt.Fprintf(o, "  %s\t%s\n", verify.fls.Name(), verify.description)
			return
		}
		set.fls.PrintDefaults()
	}
}

func parse(fls *flag.FlagSet, args []string) {
	if err := fls.Parse(args[1:]); err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
		}
		os.Exit(1)
	}
	if target == "" {
		fls.Usage()
		os.Exit(1)
	}
}
