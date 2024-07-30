package internal

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tufin/oasdiff/checker/localizations"
	"golang.org/x/exp/slices"
)

func initViper(cmd *cobra.Command, v *viper.Viper) *ReturnError {
	if err := readConfFile(v); err != nil {
		return getErrConfigFileProblem(err)
	}

	if err := validate(v); err != nil {
		return getErrConfigFileProblem(err)
	}

	if err := bindFlags(cmd, v); err != nil {
		return getErrConfigFileProblem(err)
	}

	return nil
}

func readConfFile(v *viper.Viper) error {
	v.SetConfigName(".oasdiff")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error
		} else {
			return fmt.Errorf("config file error: %s \n", err)
		}
	}

	return nil
}

func validate(v *viper.Viper) error {
	type Config struct {
		Lang       string   `mapstructure:"lang"`
		Attributes []string `mapstructure:"attributes"`
	}

	var config Config

	if err := v.Unmarshal(&config); err != nil {
		return fmt.Errorf("config file error: %s \n", err)
	}

	if config.Lang != "" && !slices.Contains(localizations.GetSupportedLanguages(), config.Lang) {
		return fmt.Errorf("invalid language %q, supported languages are: %v", config.Lang, localizations.GetSupportedLanguages())
	}

	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) error {
	var result error
	cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		name := flag.Name
		if err := v.BindPFlag(name, cmd.PersistentFlags().Lookup(name)); err != nil {
			result = fmt.Errorf("error binding flag %q to viper", name)
			return
		}
	})

	return result
}

// fixViperStringSlice fixes a limitation in viper that doesn't handle custom flags with multiple values
func fixViperStringSlice(viperString []string) []string {
	// viper returns a slice with a single element if the flag was set with a comma-separated list
	if len(viperString) == 1 {
		return strings.Split(viperString[0], ",")
	}

	return viperString
}
