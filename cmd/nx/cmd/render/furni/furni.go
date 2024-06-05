package furni

import (
	"github.com/spf13/cobra"
)

var opts struct {
	swfPath string
	states  bool
}

var Cmd = &cobra.Command{
	Use:  "furni [flags] identifier",
	Args: cobra.MaximumNArgs(1),
	RunE: run,
}

func init() {
	f := Cmd.Flags()

	f.StringVar(&opts.swfPath, "swf", "", "Path to a furni library in SWF format.")
	f.BoolVar(&opts.states, "states", false, "Print number of states.")
}

func run(cmd *cobra.Command, args []string) (err error) {

	return
}
