package get

import (
	"github.com/spf13/cobra"

	root "xabbo.b7c.io/nx/cmd/nx/cmd"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets various resources",
}

func init() {
	root.Cmd.AddCommand(getCmd)
}
