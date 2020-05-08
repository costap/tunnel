package cmd

import (
	"github.com/spf13/cobra"
)

// hostsCmd represents the hosts command
var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "Commands to manage hosts",
	Long: "",
}

func init() {
	rootCmd.AddCommand(hostsCmd)

	hostsCmd.Flags().StringP("region", "r", "lon1", "Region for the host")
	hostsCmd.Flags().StringP("token", "t", "", "API token for provisioning host")
	hostsCmd.Flags().StringP("token", "t", "", "API token for provisioning host")
	hostsCmd.Flags().StringP("token", "t", "", "API token for provisioning host")
}
