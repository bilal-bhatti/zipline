package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bilal-bhatti/zipline/internal/debug"

	"github.com/bilal-bhatti/zipline/internal"
	"github.com/google/subcommands"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&defaultCmd{}, "")

	flag.Parse()

	allCmds := map[string]bool{
		"commands": true, // builtin
		"help":     true, // builtin
		"flags":    true, // builtin
		"gen":      true,
	}

	// Default to running the "gen" command.
	if args := flag.Args(); len(args) == 0 || !allCmds[args[0]] {
		defCmd := &defaultCmd{}
		os.Exit(int(defCmd.Execute(context.Background(), flag.CommandLine)))
	}
	os.Exit(int(subcommands.Execute(context.Background())))
}

type defaultCmd struct {
	debug bool
}

func (*defaultCmd) Name() string { return "gen" }

func (*defaultCmd) Synopsis() string {
	return "Generate handler functions and Open API specs from Go code."
}

func (*defaultCmd) Usage() string {
	return `gen [packages]:
	Generates bindings_gen.go for given packages.
	If no packages provided, defaults to ".".
	example: "zipline gen ./..."
	example: "zipline gen -d ./..."
  `
}

func (p *defaultCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.debug, "d", false, "run with trace logging enabled")
}

func (p *defaultCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	zipline := internal.NewZipline()

	debug.Debug = p.debug

	err := zipline.Start(packages(f))

	if err != nil {
		log.Println(fmt.Errorf("%s", err.Error()))
	}

	fmt.Println()
	return subcommands.ExitSuccess
}

func packages(f *flag.FlagSet) []string {
	pkgs := f.Args()
	if len(pkgs) == 0 {
		pkgs = []string{"."}
	}
	return pkgs
}
