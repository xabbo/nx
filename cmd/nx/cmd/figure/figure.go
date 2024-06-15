package figure

import (
	"github.com/spf13/cobra"

	root "xabbo.b7c.io/nx/cmd/nx/cmd"
)

var Cmd = &cobra.Command{
	Use: "figure",
}

func init() {
	root.Cmd.AddCommand(Cmd)
}
