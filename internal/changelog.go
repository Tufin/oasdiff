package internal

import "github.com/spf13/cobra"

func getChangelogCmd() *cobra.Command {

	return &cobra.Command{
		Use:   `changelog`,
		Short: "Display changelog",
		Args:  cobra.ExactArgs(1),
		Run:   getChangeLog,
	}
}

func getChangeLog(cmd *cobra.Command, args []string) {

	failEmpty, returnErr := runDiffInternal(cmd, args)
	exit(failEmpty, returnErr, cmd.ErrOrStderr())
}
