package figure

import (
	"github.com/spf13/cobra"

	root "cli/cmd"
)

var figureCmd = &cobra.Command{
	Use: "figure",
}

func init() {
	root.Cmd.AddCommand(figureCmd)
}
