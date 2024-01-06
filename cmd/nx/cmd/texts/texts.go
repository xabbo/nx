package texts

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/b7c/nx"

	root "cli/cmd"
	"cli/util"
)

var textsCmd = &cobra.Command{
	Use:   "texts",
	Short: "List and search external texts",
	RunE:  runTexts,
}

var (
	searchKey   util.Wildcard
	searchValue util.Wildcard
)

func init() {
	root.Cmd.AddCommand(textsCmd)

	textsCmd.Flags().VarP(&searchKey, "key", "k", "Key search text")
	textsCmd.Flags().VarP(&searchValue, "value", "v", "Value search text")
}

func runTexts(cmd *cobra.Command, args []string) (err error) {
	mgr := nx.NewGamedataManager(root.Host)
	err = util.LoadTexts(mgr)
	if err != nil {
		return
	}

	for k, v := range mgr.Texts {
		if !filterText(k, v) {
			fmt.Printf("%s=%s\n", k, v)
		}
	}

	return
}

func filterText(key, value string) bool {
	return searchKey.Filter(key) || searchValue.Filter(value)
}
