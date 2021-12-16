package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/google/subcommands"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

type initCmd struct{}

func (*initCmd) Name() string { return "init" }

func (*initCmd) Synopsis() string {
	return "generate starter ZiplineTemplate in the specified package"
}

func (*initCmd) Usage() string {
	return `

init [package]

generate starter ZiplineTemplate in the specified package, defaults to "."
	example: "zipline init ./pkg/web"
	example: "zipline init"

`
}

func (c *initCmd) SetFlags(f *flag.FlagSet) {
}

func (c *initCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	dests := packageParam(f)
	if len(dests) > 1 {
		log.Println("specify one destination package")
		log.Println(c.Usage())
		return subcommands.ExitFailure
	}

	dest := dests[0]

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("unexpected error", err)
		return subcommands.ExitFailure
	}

	outPath := path.Join(cwd, dest, "bindings.go")

	cfg := &packages.Config{
		Mode: packages.NeedName,
		Dir:  cwd,
	}

	pkgs, err := packages.Load(cfg, dest)
	if err != nil {
		log.Fatalln("failed to load package", err)
		return subcommands.ExitFailure
	}

	if len(pkgs) > 1 {
		log.Println("multiple packages found, specify one destination package")
		log.Println(c.Usage())
		return subcommands.ExitFailure
	}

	name := pkgs[0].Name
	if name == "" {
		log.Println("invalid package specified")
		log.Println(c.Usage())
		return subcommands.ExitFailure
	}

	outFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return subcommands.ExitFailure
	}
	defer outFile.Close()

	t := template.Must(template.New("zTemplate").Parse(zTemplate))

	err = t.Execute(outFile, input{PackageName: name})
	if err != nil {
		log.Fatal("failed to write ZiplineTemplate", err)
		return subcommands.ExitFailure
	}

	log.Println("wrote ZiplineTemplate to", outPath)

	// read back file and format
	content, err := ioutil.ReadFile(outPath)
	if err != nil {
		log.Fatalln("failed to read generated file for formatting", err)
		return subcommands.ExitFailure
	}

	opt := imports.Options{
		Comments:   true,
		FormatOnly: false,
	}
	content, err = imports.Process(outPath, content, &opt)
	if err != nil {
		log.Fatalln("failed to format code", err)
		return subcommands.ExitFailure
	}

	err = ioutil.WriteFile(outPath, content, os.ModePerm)
	if err != nil {
		log.Fatalln("failed to write formatted file", err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

type input struct {
	PackageName string
}
