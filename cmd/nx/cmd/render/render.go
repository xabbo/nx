package render

import (
	"github.com/spf13/cobra"

	root "xabbo.b7c.io/nx/cmd/nx/cmd"
	"xabbo.b7c.io/nx/cmd/nx/cmd/render/furni"
)

var Cmd = &cobra.Command{
	Use:   "render",
	Short: "Render resources to images",
}

func init() {
	root.Cmd.AddCommand(Cmd)

	Cmd.AddCommand(furni.Cmd)
}
