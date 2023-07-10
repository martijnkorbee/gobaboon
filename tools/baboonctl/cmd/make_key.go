package cmd

import (
	"github.com/martijnkorbee/gobaboon/pkg/util"
	ctl "github.com/martijnkorbee/gobaboon/tools/baboonctl/internal/util"
	"github.com/spf13/cobra"
)

var (
	keylength int // default 32
)

var makeKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "Create an encryption key of n length, default 32",
	Long:  "Creates an encrytion key of n length, the default length is 32. Using the crypto/rand package.",
	Run: func(cmd *cobra.Command, args []string) {
		key, err := util.RandomStringGenerator(keylength)
		if err != nil {
			ctl.PrintError("failed to create key", err)
			return
		}
		ctl.PrintResult("your key", key)
	},
}

func init() {
	makeKeyCmd.Flags().IntVarP(&keylength, "length", "n", 32, "key length")
}
