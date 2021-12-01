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

var Version = "DEV"

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&genCmd{}, "")
	subcommands.Register(&initCmd{}, "")

	flag.Parse()

	allCmds := map[string]bool{
		"commands": true, // builtin
		"help":     true, // builtin
		"flags":    true, // builtin
		"gen":      true,
		"init":     true,
	}

	log.Printf("v%s", Version)

	// Default to running the "gen" command.
	if args := flag.Args(); len(args) == 0 || !allCmds[args[0]] {
		defCmd := &genCmd{}
		os.Exit(int(defCmd.Execute(context.Background(), flag.CommandLine)))
	}
	os.Exit(int(subcommands.Execute(context.Background())))
}

type genCmd struct {
	debug bool
}

func (*genCmd) Name() string { return "gen" }

func (*genCmd) Synopsis() string {
	return "generate handler functions and Open API specs from Go code"
}

func (*genCmd) Usage() string {
	return `

gen [packages]:

generate documentation and bindings for template in specified package, defaults to "."
	example: "zipline ./..."
	example: "zipline gen ./..."
	example: "zipline gen -d ./..."

`
}

func (p *genCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.debug, "d", false, "run with trace logging enabled")
}

func (p *genCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	zipline := internal.NewZipline()

	debug.Debug = p.debug

	err := zipline.Start(packageParam(f))

	if err != nil {
		log.Println(fmt.Errorf("%s", err.Error()))
	}

	log.Println("-----")
	return subcommands.ExitSuccess
}

func packageParam(f *flag.FlagSet) []string {
	pkgs := f.Args()
	if len(pkgs) == 0 {
		pkgs = []string{"."}
	}
	return pkgs
}
