package search

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"xabbo.io/nx"
	gd "xabbo.io/nx/gamedata"

	_root "xabbo.io/nx/cmd/nx/cmd"
	_parent "xabbo.io/nx/cmd/nx/cmd/furni"
	"xabbo.io/nx/cmd/nx/util"
)

var opts struct {
	searchName       util.Wildcard
	searchIdentifier util.Wildcard
	searchCategory   util.Wildcard
	searchLine       util.Wildcard
	searchParam      util.Wildcard
	searchTypeStr    string
	searchType       nx.ItemType
	specialTypes     []int
	json             bool
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
	f.StringVarP(&opts.searchTypeStr, "type", "t", "", "The furni type (floor/wall)")
	f.IntSliceVarP(&opts.specialTypes, "special-types", "s", []int{}, "The furni special types")
	f.BoolVar(&opts.json, "json", false, "Output furni info in JSON format")

	_parent.Cmd.AddCommand(Cmd)
}

func runSearch(cmd *cobra.Command, args []string) (err error) {
	err = opts.searchName.Set(strings.Join(args, " "))
	if err != nil {
		return
	}

	switch opts.searchTypeStr {
	case "":
		opts.searchType = 'x'
	case "s", "f", "floor":
		opts.searchType = nx.ItemFloor
	case "i", "w", "wall":
		opts.searchType = nx.ItemWall
	}

	cmd.SilenceUsage = true

	mgr := gd.NewManager(_root.Host)
	err = util.LoadFurni(mgr)
	if err != nil {
		return
	}

	matches := []*gd.FurniInfo{}
	for _, f := range mgr.Furni() {
		if !filterFurni(f) {
			matches = append(matches, f)
		}
	}

	if opts.json {
		json.NewEncoder(os.Stdout).Encode(matches)
	} else {
		for _, f := range matches {
			fmt.Printf("%s [%s]\n", f.Name, f.Identifier)
		}
	}

	return
}

func filterFurni(f *gd.FurniInfo) bool {
	return (opts.searchType != 'x' && f.Type != opts.searchType) ||
		opts.searchName.Filter(f.Name) ||
		opts.searchIdentifier.Filter(f.Identifier) ||
		opts.searchCategory.Filter(f.Category) ||
		opts.searchLine.Filter(f.Line) ||
		opts.searchParam.Filter(f.CustomParams) ||
		(len(opts.specialTypes) > 0 && !slices.Contains(opts.specialTypes, int(f.SpecialType)))
}
