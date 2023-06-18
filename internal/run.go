package internal

import (
	"io"
	"strconv"

	"github.com/spf13/cobra"
)

func Run(args []string, stdout io.Writer, stderr io.Writer) int {
	rootCmd := &cobra.Command{
		Use:   `oasdiff`,
		Short: "Compare and detect breaking changes in OpenAPI specs",
	}

	rootCmd.SetArgs(args[1:])
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)

	rootCmd.AddCommand(
		getDiffCmd(),
		getSummaryCmd(),
		getBreakingChangesCmd(),
		getChangelogCmd(),
		getLintCmd(),
	)

	rootCmd.Flags().StringP("version", "v", "", "show version and quit")

	if err := rootCmd.Execute(); err != nil {
		if ret := getReturnValue(rootCmd); ret != 0 {
			return ret
		}
		return 100
	}

	return getReturnValue(rootCmd)
}

func setReturnValue(cmd *cobra.Command, code int) {
	if cmd.Root().Annotations == nil {
		cmd.Root().Annotations = map[string]string{}
	}

	cmd.Root().Annotations["return"] = strconv.Itoa(code)
}

func getReturnValue(cmd *cobra.Command) int {
	if cmd.Root().Annotations == nil {
		return 0
	}

	codeStr := cmd.Root().Annotations["return"]
	if codeStr == "" {
		return 0
	}

	code, err := strconv.Atoi(codeStr)
	if err != nil {
		// this shouldn't happen
		return 0
	}

	return code
}
