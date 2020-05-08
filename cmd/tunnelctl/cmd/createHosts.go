package cmd

import (
	"github.com/costap/tunnel/internal/pkg/hosting/do"
	"github.com/spf13/cobra"
)

// createHostsCmd represents the create command
var createHostsCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new tunnel host",
	Long: "",
	Run: func(cmd *cobra.Command, args []string) {
		t := cmd.Flag("token").Value.String()
		r := cmd.Flag("region").Value.String()
		i := cmd.Flag("image").Value.String()
		p := cmd.Flag("path").Value.String()
		n := cmd.Flag("name").Value.String()

		prov := do.NewDOProvider(t, r, i, p, n)

		ip, err := prov.CreateInstance()
		if err != nil {
			cmd.PrintErrf("error creating host: %v", err)
		}
		cmd.Printf("Host created with ip %v", ip)
	},
}

func init() {
	hostsCmd.AddCommand(createHostsCmd)

	createHostsCmd.Flags().StringP("sshPath", "p", ".", "Path where ssh keys are store")
	createHostsCmd.Flags().String("sshName", "id_rsa", "Name of the ssh to use or create")
	createHostsCmd.Flags().StringP("image", "i", "ubuntu-18-04-x64", "Image to use for the host")
	createHostsCmd.Flags().StringP("name", "n", "tunnelctl", "Host name to use for the host")
}
