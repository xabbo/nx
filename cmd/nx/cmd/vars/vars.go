package vars

import (
	"fmt"

	"github.com/spf13/cobra"

	gd "github.com/xabbo/nx/gamedata"

	root "cli/cmd"
	"cli/util"
)

var varsCommand = &cobra.Command{
	Use:   "vars",
	Short: "List and search external variables",
	RunE:  runVars,
}

var (
	searchKey   util.Wildcard
	searchValue util.Wildcard
)

func init() {
	root.Cmd.AddCommand(varsCommand)

	varsCommand.Flags().VarP(&searchKey, "key", "k", "Key search text")
	varsCommand.Flags().VarP(&searchValue, "value", "v", "Value search text")
}

func runVars(cmd *cobra.Command, args []string) (err error) {
	mgr := gd.NewGamedataManager(root.Host)
	err = util.LoadGamedata(mgr, "Loading external variables...", gd.GamedataVariables)
	if err != nil {
		return
	}

	for k, v := range mgr.Variables {
		if !filterVar(k, v) {
			fmt.Printf("%s=%s\n", k, v)
		}
	}

	return
}

func filterVar(k, v string) bool {
	return searchKey.Filter(k) || searchValue.Filter(v)
}
