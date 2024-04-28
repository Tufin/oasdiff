package internal

import (
	"fmt"

	"github.com/spf13/cobra"
)

func getQRCodeCmd() *cobra.Command {

	cmd := cobra.Command{
		Use:               "qr",
		Short:             "Display QR code of oasdiff repo",
		Long:              "Display QR code of the URL of the oasdiff repository",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions, // see https://github.com/spf13/cobra/issues/1969
		RunE: func(cmd *cobra.Command, args []string) error {

			fmt.Fprintf(cmd.OutOrStdout(), "%s\n",
				`			▄▄▄▄▄▄▄  ▄  ▄▄  ▄▄▄▄  ▄▄▄▄▄▄▄  
			█ ▄▄▄ █ ▄▄█▀█▄█▀█▀ █  █ ▄▄▄ █  
			█ ███ █  ▀█▀ ▀█▀ ▄██  █ ███ █  
			█▄▄▄▄▄█ ▄▀▄ ▄▀▄ ▄ █▀█ █▄▄▄▄▄█  
			▄▄▄ ▄▄▄▄█▀█▄▀▀█▄▀█▄▀▀▄▄   ▄    
			▀██▄▀ ▄▀█▄▀▀ ▄▀▀ ▄  █▄█ ▄▀▄▄█  
			 █▄█ █▄██  ▀▄  ▀▄█▀▀▀▄▄▀▀▄ █▄  
			 ▀▄█▀▄▄▀  ▄██▀▄██▄ ▀▀▀█ ▄█ ▄█  
			█▄▀▀▄▀▄▀█▄▀▄▀▀█▄██▀▄█▀█▄ █ █▄  
			▄▀▄  ▄▄▄▀▀▀▀ ▄▀▀▄▀▄  █▀▄▄▀▀▄█  
			▄▀▄▄▄▄▄██ ▄▀▄   ██▀██████▀ ▀   
			▄▄▄▄▄▄▄ █▄ ██▀▄ █████ ▄ █▄▀██  
			█ ▄▄▄ █ █▄▄▄▀▀█ ██▀ █▄▄▄█▀ ▀▄  
			█ ███ █ ▄█ ▀ ▄▀ ▀▀ ▄▀  ██▄▀▀█  
			█▄▄▄▄▄█ █▀▄▀▄  █▄█ ███ ▄█  █▄`)
			return nil
		},
	}

	return &cmd
}
