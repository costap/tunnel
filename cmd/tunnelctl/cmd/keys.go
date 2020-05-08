package cmd

import (
	"github.com/spf13/cobra"
)

// keysCmd represents the keys command
var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Commands to manage keys",
	Long: "",
}

func init() {
	rootCmd.AddCommand(keysCmd)
}
