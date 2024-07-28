package internal

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func readConfFile(v *viper.Viper) {
	v.SetConfigName(".oasdiff")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error
		} else {
			log.Fatal(fmt.Errorf("config file error: %s \n", err))
		}
	}
}

func bindViperFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		name := flag.Name
		if err := v.BindPFlag(name, cmd.PersistentFlags().Lookup(name)); err != nil {
			log.Fatalf("error binding flag %q to viper", name)
		}
	})
}

// fixViperStringSlice fixes a limitation in viper that doesn't handle custom flags with multiple values
func fixViperStringSlice(viperString []string) []string {
	// viper returns a slice with a single element if the flag was set with a comma-separated list
	if len(viperString) == 1 {
		return strings.Split(viperString[0], ",")
	}

	return viperString
}
