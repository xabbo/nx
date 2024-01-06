package get

import (
	"github.com/spf13/cobra"

	root "cli/cmd"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets various resources",
}

func init() {
	root.Cmd.AddCommand(getCmd)
}
