package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	rebaseCmd = &cobra.Command{
		Use:   "rebase",
		Short: "Perform all required rebases for current branch",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("rebasing")
		},
	}
)

func init() {
	rootCmd.AddCommand(rebaseCmd)
}
