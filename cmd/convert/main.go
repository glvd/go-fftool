package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fftool",
	Short: "fftool is a very fast video conversion tool",
	Long:  `A Fast Video Conversion Tool`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func main() {
	rootCmd.AddCommand(formatCmd())
	if err := rootCmd.Execute(); err != nil {
		return
	}
}
