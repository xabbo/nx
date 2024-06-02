package render

import (
	"github.com/spf13/cobra"

	root "github.com/xabbo/nx/cli/nx/cmd"
	"github.com/xabbo/nx/cli/nx/cmd/render/furni"
)

var Cmd = &cobra.Command{
	Use:   "render",
	Short: "Render resources to images",
}

func init() {
	root.Cmd.AddCommand(Cmd)

	Cmd.AddCommand(furni.Cmd)
}
