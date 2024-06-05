package figure

import (
	"github.com/spf13/cobra"

	root "xabbo.b7c.io/nx/cmd/nx/cmd"
)

var figureCmd = &cobra.Command{
	Use: "figure",
}

func init() {
	root.Cmd.AddCommand(figureCmd)
}
