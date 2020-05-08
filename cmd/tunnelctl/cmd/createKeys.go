package cmd

import (
	"github.com/costap/tunnel/internal/pkg/keys"
	"log"

	"github.com/spf13/cobra"
)

// createKeysCmd represents the create command
var createKeysCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new ssh keys pair",
	Long: "",
	Run: func(cmd *cobra.Command, args []string) {
		p := cmd.Flag("path").Value.String()
		n := cmd.Flag("name").Value.String()

		err := keys.GenerateKeyPair(p, n)
		if err != nil {
			log.Fatalf("error creating key pair %v", err)
		}
	},
}

func init() {
	keysCmd.AddCommand(createKeysCmd)

	createKeysCmd.Flags().StringP("path", "p", ".", "Path where keys will be created")
	createKeysCmd.Flags().StringP("name", "n", "id_rsa", "Path where keys will be created")
}
