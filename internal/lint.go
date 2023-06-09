package internal

import "github.com/spf13/cobra"

func getLintCmd() *cobra.Command {

	return &cobra.Command{
		Use:   `lint`,
		Short: "report problems in an OpenAPI spec",
		Args:  cobra.ExactArgs(1),
		Run:   getLint,
	}
}

func getLint(cmd *cobra.Command, args []string) {

	failEmpty, returnErr := runDiffInternal(cmd, args)
	exit(failEmpty, returnErr, cmd.ErrOrStderr())
}
