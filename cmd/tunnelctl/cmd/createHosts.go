package cmd

import (
	"github.com/costap/tunnel/internal/pkg/hosting/do"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"strings"
)

// createHostsCmd represents the create command
var createHostsCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new tunnel host",
	Long: "",
	Run: func(cmd *cobra.Command, args []string) {
		t := viper.GetString("token")
		r := viper.GetString("region")
		i := viper.GetString("image")
		sshP := viper.GetString("sshPath")
		sshN := viper.GetString("sshName")

		prov := do.NewDOProvider(t)
		host := do.NewHostConfig(r,viper.GetString("name"),i,sshP,sshN)

		for _, proxy := range proxies {
			var source, target string
			parts := strings.Split(proxy, ":")
			switch len(parts) {
			case 1:
				source = parts[0]
				target = source
			case 2:
				source = parts[0]
				target = parts[1]
			default:
				log.Fatalf("Wrong proxy ports format: %v should be <source>:<target>", proxy)
			}
			si, err := strconv.Atoi(source)
			if err != nil  {
				log.Fatalf("Invalid proxy ports source: %v should be int", source)
			}
			ti, err := strconv.Atoi(target)
			if err != nil  {
				log.Fatalf("Invalid proxy ports target: %v should be int", target)
			}
			host.AddProxy(si,ti)
		}


		ip, err := prov.CreateInstance(host)
		if err != nil {
			cmd.PrintErrf("error creating host: %v", err)
		}
		cmd.Printf("Host created with ip %v", ip)
	},
}

var proxies []string

func init() {
	hostsCmd.AddCommand(createHostsCmd)

	createHostsCmd.Flags().StringP("sshPath", "p", ".", "Path where ssh keys are store")
	createHostsCmd.Flags().String("sshName", "id_rsa", "Name of the ssh to use or create")
	createHostsCmd.Flags().StringP("image", "i", "ubuntu-18-04-x64", "Image to use for the host")
	createHostsCmd.Flags().StringP("name", "n", "tunnelctl", "Host name to use for the host")
	createHostsCmd.Flags().StringP("token", "t", "", "API token for provisioning host")
	createHostsCmd.Flags().StringP("region", "r", "lon1", "API token for provisioning host")
	createHostsCmd.Flags().StringSliceVar( &proxies,"proxy", []string{}, "Proxy ports in the format <source>:<target>")
	viper.BindPFlag("sshPath", createHostsCmd.LocalFlags().Lookup("sshPath"))
	viper.BindPFlag("sshName", createHostsCmd.LocalFlags().Lookup("sshName"))
	viper.BindPFlag("image", createHostsCmd.LocalFlags().Lookup("image"))
	viper.BindPFlag("token", createHostsCmd.LocalFlags().Lookup("token"))
	viper.BindPFlag("region", createHostsCmd.LocalFlags().Lookup("region"))
	viper.BindPFlag("name", createHostsCmd.LocalFlags().Lookup("name"))
}
