package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bilal-bhatti/zipline/internal"
	"github.com/bilal-bhatti/zipline/internal/debug"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zipline [flags] [<package-glob>, defaults to .]",
	Short: "generate handler functions and Open API specs from Go code",
	Long: `generate documentation, Open API specs, and handler functions from ZiplineTemplate 

examples: 
zipline ./...
zipline --debug ./..."
`,
	// Args: cobra.ExactArgs(1),
	// Args: cobra.MinimumNArgs(1),
	Args: cobra.MaximumNArgs(1),
	// PreRun: func(cmd *cobra.Command, args []string) {
	// 	log.Printf("v%s\n", Version)
	// },

	Run: run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().BoolVar(&debug.Debug, "debug", false, "run with debug")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("curl", "c", false, "print cURL command")
	// rootCmd.Flags().StringP("env", "e", "", "environment settings file")
}

func run(cmd *cobra.Command, args []string) {
	zipline := internal.NewZipline()

	err := zipline.Start(packageList(args))

	if err != nil {
		log.Println(fmt.Errorf("%s", err.Error()))
	}

	log.Println("-----")
}

func packageList(args []string) []string {
	if len(args) == 0 {
		return []string{"."}
	}
	return args
}
