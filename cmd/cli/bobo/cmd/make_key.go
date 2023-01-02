package cmd

import (
	"github.com/martijnkorbee/gobaboon/cmd/cli/internal/pkg/util"
	butil "github.com/martijnkorbee/gobaboon/pkg/util"
	"github.com/spf13/cobra"
)

var (
	keylength int // default 32
)

var makeKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "Create an encryption key of n length, default 32",
	Long: `Creates an encrytion key of n length, the default length is 32.
	We use the crypto/rand package to get random numbers.`,
	Run: func(cmd *cobra.Command, args []string) {
		key, err := butil.RandomStringGenerator(keylength)
		if err != nil {
			util.PrintError("failed to create key", err)
			return
		}
		util.PrintResult("your key", key)
	},
}

func init() {
	makeKeyCmd.Flags().IntVarP(&keylength, "length", "n", 32, "key length")
}
