package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/bilal-bhatti/zipline/internal/debug"

	"github.com/bilal-bhatti/zipline/internal"
	"github.com/google/subcommands"
)

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

generate documentation and bindings for ZiplineTemplate in specified package, defaults to "."
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
