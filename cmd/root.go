package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var version = "devel"

var rootCmd = &cobra.Command{
	Use:     "go-stacks",
	Short:   "Command like tool to handle git PR stacks ",
	Long:    `TODO: changeme`,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		Verbose(cmd)
	},
}

// Verbose Increase verbosity
func Verbose(cmd *cobra.Command) {
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		log.Panic(err)
	}

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}
func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Increase verbosity")
	rootCmd.PersistentFlags().BoolP("dryrun", "n", false, "Dry run")
}

// Execute The main function for the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
