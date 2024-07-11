package vars

import (
	"fmt"

	"github.com/spf13/cobra"

	gd "xabbo.b7c.io/nx/gamedata"

	_root "xabbo.b7c.io/nx/cmd/nx/cmd"
	"xabbo.b7c.io/nx/cmd/nx/util"
)

var Cmd = &cobra.Command{
	Use:   "vars",
	Short: "List and search external variables",
	RunE:  runVars,
}

var opts struct {
	searchKey   util.Wildcard
	searchValue util.Wildcard
}

func init() {
	f := Cmd.Flags()
	f.VarP(&opts.searchKey, "key", "k", "Key search text")
	f.VarP(&opts.searchValue, "value", "v", "Value search text")

	_root.Cmd.AddCommand(Cmd)
}

func runVars(cmd *cobra.Command, args []string) (err error) {
	mgr := gd.NewManager(_root.Host)
	err = util.LoadGameData(mgr, "Loading external variables...", gd.GameDataVariables)
	if err != nil {
		return
	}

	for k, v := range mgr.Variables() {
		if !filterVar(k, v) {
			fmt.Printf("%s=%s\n", k, v)
		}
	}

	return
}

func filterVar(k, v string) bool {
	return opts.searchKey.Filter(k) || opts.searchValue.Filter(v)
}
