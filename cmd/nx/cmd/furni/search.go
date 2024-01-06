package furni

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/b7c/nx"

	root "cli/cmd"
	"cli/util"
)

var (
	searchName       util.Wildcard
	searchIdentifier util.Wildcard
	searchCategory   util.Wildcard
	searchLine       util.Wildcard
	searchParam      util.Wildcard
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for furni.",
	RunE:  runSearch,
}

func init() {
	searchCmd.Flags().VarP(&searchIdentifier, "identifier", "i", "The furni identifier")
	searchCmd.Flags().VarP(&searchCategory, "category", "c", "The furni category")
	searchCmd.Flags().VarP(&searchLine, "line", "l", "The furni line")
	searchCmd.Flags().VarP(&searchParam, "param", "p", "The furni parameters")

	furniCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) (err error) {
	err = searchName.Set(strings.Join(args, " "))
	if err != nil {
		return
	}

	cmd.SilenceUsage = true

	mgr := nx.NewGamedataManager(root.Host)
	err = util.LoadFurni(mgr)
	if err != nil {
		return
	}

	for _, f := range mgr.Furni {
		if !filterFurni(f) {
			fmt.Printf("%s [%s]\n", f.Name, f.Identifier)
		}
	}

	return
}

func filterFurni(f nx.FurniInfo) bool {
	return searchName.Filter(f.Name) ||
		searchIdentifier.Filter(f.Identifier) ||
		searchCategory.Filter(f.Category) ||
		searchLine.Filter(f.Line) ||
		searchParam.Filter(f.CustomParams)
}
