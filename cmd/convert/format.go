package main

import "github.com/spf13/cobra"

func formatCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "format",
		Short: "get the video format for conversion",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
}
