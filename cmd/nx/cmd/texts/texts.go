package texts

import (
	"fmt"

	"github.com/spf13/cobra"

	gd "xabbo.io/nx/gamedata"

	_root "xabbo.io/nx/cmd/nx/cmd"
	"xabbo.io/nx/cmd/nx/util"
)

var Cmd = &cobra.Command{
	Use:   "texts",
	Short: "List and search external texts",
	RunE:  runTexts,
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

func runTexts(cmd *cobra.Command, args []string) (err error) {
	mgr := gd.NewManager(_root.Host)
	err = util.LoadTexts(mgr)
	if err != nil {
		return
	}

	for k, v := range mgr.Texts() {
		if !filterText(k, v) {
			fmt.Printf("%s=%s\n", k, v)
		}
	}

	return
}

func filterText(key, value string) bool {
	return opts.searchKey.Filter(key) || opts.searchValue.Filter(value)
}
