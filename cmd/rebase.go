package cmd

import (
	"fmt"

	"github.com/mhristof/stacks/bash"
	"github.com/mhristof/stacks/git"
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

			branch, err := cmd.Flags().GetString("branch")
			if err != nil {
				panic(err)
			}

			fmt.Println("rebasing (branch:", branch, ") (dry:", dry, ")")
			commands, err := git.Rebase(".", branch)
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
	rebaseCmd.PersistentFlags().StringP("branch", "b", git.Branch("."), "Branch regex to match")
	rootCmd.AddCommand(rebaseCmd) //nolint
}
