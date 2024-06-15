package texts

import (
	"fmt"

	"github.com/spf13/cobra"

	gd "xabbo.b7c.io/nx/gamedata"

	root "xabbo.b7c.io/nx/cmd/nx/cmd"
	"xabbo.b7c.io/nx/cmd/nx/util"
)

var textsCmd = &cobra.Command{
	Use:   "texts",
	Short: "List and search external texts",
	RunE:  runTexts,
}

var opts struct {
	searchKey   util.Wildcard
	searchValue util.Wildcard
}

func init() {
	root.Cmd.AddCommand(textsCmd)

	textsCmd.Flags().VarP(&opts.searchKey, "key", "k", "Key search text")
	textsCmd.Flags().VarP(&opts.searchValue, "value", "v", "Value search text")
}

func runTexts(cmd *cobra.Command, args []string) (err error) {
	mgr := gd.NewManager(root.Host)
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
