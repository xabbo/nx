package figure

import (
	"github.com/spf13/cobra"

	root "github.com/xabbo/nx/cli/nx/cmd"
)

var figureCmd = &cobra.Command{
	Use: "figure",
}

func init() {
	root.Cmd.AddCommand(figureCmd)
}
