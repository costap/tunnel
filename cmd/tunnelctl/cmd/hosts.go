package cmd

import (
	"github.com/spf13/cobra"
)

// hostsCmd represents the hosts command
var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "Commands to manage hosts",
	Long:  "",
}

func init() {
	rootCmd.AddCommand(hostsCmd)
}
