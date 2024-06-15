package search

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	gd "xabbo.b7c.io/nx/gamedata"

	_root "xabbo.b7c.io/nx/cmd/nx/cmd"
	_parent "xabbo.b7c.io/nx/cmd/nx/cmd/furni"
	"xabbo.b7c.io/nx/cmd/nx/util"
)

var opts struct {
	searchName       util.Wildcard
	searchIdentifier util.Wildcard
	searchCategory   util.Wildcard
	searchLine       util.Wildcard
	searchParam      util.Wildcard
}

var Cmd = &cobra.Command{
	Use:   "search",
	Short: "Search for furni.",
	RunE:  runSearch,
}

func init() {
	f := Cmd.Flags()
	f.VarP(&opts.searchIdentifier, "identifier", "i", "The furni identifier")
	f.VarP(&opts.searchCategory, "category", "c", "The furni category")
	f.VarP(&opts.searchLine, "line", "l", "The furni line")
	f.VarP(&opts.searchParam, "param", "p", "The furni parameters")

	_parent.Cmd.AddCommand(Cmd)
}

func runSearch(cmd *cobra.Command, args []string) (err error) {
	err = opts.searchName.Set(strings.Join(args, " "))
	if err != nil {
		return
	}

	cmd.SilenceUsage = true

	mgr := gd.NewManager(_root.Host)
	err = util.LoadFurni(mgr)
	if err != nil {
		return
	}

	for _, f := range mgr.Furni() {
		if !filterFurni(f) {
			fmt.Printf("%s [%s]\n", f.Name, f.Identifier)
		}
	}

	return
}

func filterFurni(f *gd.FurniInfo) bool {
	return opts.searchName.Filter(f.Name) ||
		opts.searchIdentifier.Filter(f.Identifier) ||
		opts.searchCategory.Filter(f.Category) ||
		opts.searchLine.Filter(f.Line) ||
		opts.searchParam.Filter(f.CustomParams)
}
