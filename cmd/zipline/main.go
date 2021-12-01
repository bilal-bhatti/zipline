package main

import (
	"context"
	"flag"
	"log"
	"os"


	"github.com/google/subcommands"
)

var Version = "DEV"

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	// subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&genCmd{}, "")
	subcommands.Register(&initCmd{}, "")

	flag.Parse()

	allCmds := map[string]bool{
		"commands": true, // builtin
		"help":     true, // builtin
		// "flags":    true, // builtin
		"gen":  true,
		"init": true,
	}

	log.Printf("v%s", Version)

	// Default to running the "gen" command.
	if args := flag.Args(); len(args) == 0 || !allCmds[args[0]] {
		defCmd := &genCmd{}
		os.Exit(int(defCmd.Execute(context.Background(), flag.CommandLine)))
	}
	os.Exit(int(subcommands.Execute(context.Background())))
}

func packageParam(f *flag.FlagSet) []string {
	pkgs := f.Args()
	if len(pkgs) == 0 {
		pkgs = []string{"."}
	}
	return pkgs
}
