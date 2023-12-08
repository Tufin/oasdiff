package internal

import (
	"github.com/spf13/cobra"
)

func getQRCodeCmd() *cobra.Command {

	cmd := cobra.Command{
		Use:               "qr",
		Short:             "Display QR code",
		Long:              "Display QR code",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions, // see https://github.com/spf13/cobra/issues/1969
		RunE: func(cmd *cobra.Command, args []string) error {

			println(`                                  
			▄▄▄▄▄▄▄  ▄  ▄▄  ▄▄▄▄  ▄▄▄▄▄▄▄  
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
			█▄▄▄▄▄█ █▀▄▀▄  █▄█ ███ ▄█  █▄  
										   `)
			return nil
		},
	}

	return &cmd
}
