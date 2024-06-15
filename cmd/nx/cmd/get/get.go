package get

import (
	"github.com/spf13/cobra"

	_root "xabbo.b7c.io/nx/cmd/nx/cmd"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Gets various resources",
}

func init() {
	_root.Cmd.AddCommand(Cmd)
}
