package figure

import (
	"github.com/spf13/cobra"

	_root "xabbo.b7c.io/nx/cmd/nx/cmd"
)

var Cmd = &cobra.Command{
	Use: "figure",
}

func init() {
	_root.Cmd.AddCommand(Cmd)
}
