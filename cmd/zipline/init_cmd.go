package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/spf13/cobra"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

// initCmds represents the swag command
var initCmds = &cobra.Command{
	Use:   "init [flags] [<package>, defaults to .]",
	Short: "generate starter ZiplineTemplate in the specified package",
	Long: `generate starter ZiplineTemplate in the specified package

examples:

	zipline init
	zipline init ./pkg/web
`,
	Run: runInit,
}

func init() {
	rootCmd.AddCommand(initCmds)
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmds.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmds.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runInit(cmd *cobra.Command, args []string) {
	dests := packageList(args)
	if len(dests) > 1 {
		log.Fatalln("specify one destination package")
	}

	dest := dests[0]

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("unexpected error", err)
	}

	outPath := path.Join(cwd, dest, "bindings.go")

	cfg := &packages.Config{
		Mode: packages.NeedName,
		Dir:  cwd,
	}

	pkgs, err := packages.Load(cfg, dest)
	if err != nil {
		log.Fatalln("failed to load package", err)
	}

	if len(pkgs) > 1 {
		log.Fatalln("multiple packages found, specify one destination package")
	}

	name := pkgs[0].Name
	if name == "" {
		log.Fatalln("invalid package specified")
	}

	outFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	defer outFile.Close()

	t := template.Must(template.New("zTemplate").Parse(zTemplate))

	err = t.Execute(outFile, input{PackageName: name})
	if err != nil {
		log.Fatalln("failed to write ZiplineTemplate", err)

	}

	log.Println("wrote ZiplineTemplate to", outPath)

	// read back file and format
	content, err := ioutil.ReadFile(outPath)
	if err != nil {
		log.Fatalln("failed to read generated file for formatting", err)
	}

	opt := imports.Options{
		Comments:   true,
		FormatOnly: false,
	}
	content, err = imports.Process(outPath, content, &opt)
	if err != nil {
		log.Fatalln("failed to format code", err)
	}

	err = ioutil.WriteFile(outPath, content, os.ModePerm)
	if err != nil {
		log.Fatalln("failed to write formatted file", err)
	}
}

type input struct {
	PackageName string
}
