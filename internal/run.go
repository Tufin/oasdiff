package internal

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func Run(args []string, stdout io.Writer, stderr io.Writer) int {
	rootCmd := &cobra.Command{
		Use:   `oasdiff`,
		Short: "Compare and detect breaking changes in OpenAPI specs",
	}

	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)

	rootCmd.AddCommand(
		getDiffCmd(),
		getBreakingChangesCmd(),
		getChangelogCmd(),
		getLintCmd(),
	)

	rootCmd.Flags().StringP("version", "v", "", "show version and quit")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(stderr, "%v\n", err)
		return 101
	}

	return 0
}

func exit(failEmpty bool, returnErr *ReturnError, stderr io.Writer) {

	if returnErr != nil {
		if returnErr.Err != nil {
			fmt.Fprintf(stderr, "%v\n", returnErr.Err)
		}
		os.Exit(returnErr.Code)
	}

	if failEmpty {
		os.Exit(1)
	}

	os.Exit(0)
}
