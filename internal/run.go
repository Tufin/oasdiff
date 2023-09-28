package internal

import (
	"io"
	"os"
	"strconv"
	"time"

	"github.com/oasdiff/go-common/util"
	"github.com/oasdiff/telemetry/client"
	"github.com/oasdiff/telemetry/model"
	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/build"
)

func Run(args []string, stdout io.Writer, stderr io.Writer) int {
	rootCmd := &cobra.Command{
		Use:   "oasdiff",
		Short: "OpenAPI specification diff",
	}

	rootCmd.SetArgs(args[1:])
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.Version = build.Version

	rootCmd.AddCommand(
		getDiffCmd(),
		getSummaryCmd(),
		getBreakingChangesCmd(),
		getChangelogCmd(),
		getFlattenCmd(),
		getChecksCmd(),
	)

	return strategy()(rootCmd)
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

func strategy() func(*cobra.Command) int {

	if os.Getenv(model.EnvNoTelemetry) == "1" {
		return run
	}

	return func(cmd *cobra.Command) int {
		c := make(chan int)
		go func() {
			defer close(c)
			_ = client.NewCollector(util.NewStringSet().Add("err-ignore").
				Add("warn-ignore").
				Add("match-path").
				Add("prefix-base").
				Add("prefix-revision").
				Add("strip-prefix-base").
				Add("strip-prefix-revision").
				Add("filter-extension")).SendCommand(cmd)
		}()

		ret := run(cmd)

		select {
		case <-c:
		case <-time.After(model.DefaultTimeout):
		}

		return ret
	}
}

func run(cmd *cobra.Command) int {

	if err := cmd.Execute(); err != nil {
		if ret := getReturnValue(cmd); ret != 0 {
			return ret
		}
		return 100
	}

	return getReturnValue(cmd)
}
