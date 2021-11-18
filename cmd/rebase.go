package cmd

import (
	"fmt"

	"github.com/mhristof/go-stacks/bash"
	"github.com/mhristof/go-stacks/git"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rebaseCmd = &cobra.Command{
		Use:   "rebase",
		Short: "Perform all required rebases for current branch",
		Run: func(cmd *cobra.Command, args []string) {
			dry, err := cmd.Flags().GetBool("dryrun")
			if err != nil {
				panic(err)
			}

			fmt.Println("rebasing (dry:", dry, ")")
			commands, err := git.Rebase(".")
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("cannot calculate commands")
			}

			for _, command := range commands {
				fmt.Println(fmt.Sprintf("command: %+v", command))
				if dry {
					continue
				}

				bash.Run(command)
			}

		},
	}
)

func init() {
	rootCmd.AddCommand(rebaseCmd)
}
