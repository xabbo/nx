package get

import (
	"github.com/spf13/cobra"

	root "github.com/xabbo/nx/cli/nx/cmd"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets various resources",
}

func init() {
	root.Cmd.AddCommand(getCmd)
}
