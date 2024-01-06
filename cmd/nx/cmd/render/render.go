package render

import (
	"github.com/spf13/cobra"

	root "cli/cmd"
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render resources to images",
}

func init() {
	root.Cmd.AddCommand(renderCmd)
}
