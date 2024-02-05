package internal

import (
	"context"
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

	ctx, cancelCtx := context.WithDeadline(context.Background(), time.Now().Add(model.DefaultTimeout))
	defer cancelCtx()

	chanPreRun := make(chan int)
	rootCmd := &cobra.Command{
		Use:               "oasdiff",
		Short:             "OpenAPI specification diff",
		PersistentPreRun:  func(cmd *cobra.Command, args []string) { preRun(chanPreRun, cmd) },
		PersistentPostRun: func(*cobra.Command, []string) { postRun(ctx, chanPreRun) },
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
		getQRCodeCmd(),
		getDeltaCmd(),
	)

	return run(rootCmd)
}

func preRun(c chan int, cmd *cobra.Command) {

	go func(c chan int) {
		defer close(c)
		if os.Getenv(model.EnvNoTelemetry) != "1" {
			_ = client.NewCollector(util.NewStringSet().Add("err-ignore").
				Add("warn-ignore").
				Add("match-path").
				Add("prefix-base").
				Add("prefix-revision").
				Add("strip-prefix-base").
				Add("strip-prefix-revision").
				Add("filter-extension")).SendCommand(cmd)
		}
	}(c)
}

func postRun(ctx context.Context, c chan int) {

	select {
	case <-c:
		break
	case <-ctx.Done():
		break
	}
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

func run(cmd *cobra.Command) int {

	if err := cmd.Execute(); err != nil {
		if ret := getReturnValue(cmd); ret != 0 {
			return ret
		}
		return generalExecutionErr
	}

	return getReturnValue(cmd)
}
