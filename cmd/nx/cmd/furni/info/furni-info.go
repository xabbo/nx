package info

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"xabbo.io/nx"
	gd "xabbo.io/nx/gamedata"

	_root "xabbo.io/nx/cmd/nx/cmd"
	"xabbo.io/nx/cmd/nx/cmd/furni"
	"xabbo.io/nx/cmd/nx/util"
)

var Cmd = &cobra.Command{
	Use:  "info",
	Args: cobra.MaximumNArgs(1),
	RunE: runInfo,
}

var opts struct {
	itemType   string
	kind       int
	identifier string
	json       bool
}

func init() {
	f := Cmd.Flags()

	f.StringVarP(&opts.itemType, "type", "t", "", "The furni type (floor/wall)")
	f.IntVarP(&opts.kind, "kind", "k", 0, "The furni kind (type ID)")
	f.StringVarP(&opts.identifier, "identifier", "i", "", "The furni identifier (class name)")
	f.BoolVar(&opts.json, "json", false, "Output JSON")

	furni.Cmd.AddCommand(Cmd)
}

func runInfo(cmd *cobra.Command, args []string) (err error) {
	cmd.SilenceUsage = true

	mgr := gd.NewManager(_root.Host)
	err = util.LoadFurni(mgr)
	if err != nil {
		return
	}

	var fi *gd.FurniInfo

	if len(args) > 0 {
		identifier := args[0]
		if furniInfo, ok := mgr.Furni()[identifier]; ok {
			fi = furniInfo
		}
	} else if opts.kind > 0 {
		var itemType nx.ItemType
		switch opts.itemType {
		case "":
			return fmt.Errorf("item type is required when kind is specified")
		case "floor", "f", "s":
			itemType = nx.ItemFloor
		case "wall", "w", "i":
			itemType = nx.ItemWall
		default:
			return fmt.Errorf("invalid item type: %q", opts.itemType)
		}
		for _, f := range mgr.Furni() {
			if f.Type == itemType && f.Kind == opts.kind {
				fi = f
				break
			}
		}
	} else if opts.identifier != "" {
		for _, f := range mgr.Furni() {
			if f.Identifier == opts.identifier {
				fi = f
				break
			}
		}
	}

	if fi != nil {
		if opts.json {
			utf8, err := json.Marshal(fi)
			if err != nil {
				return err
			}
			fmt.Print(string(utf8))
		} else {
			util.RenderFurniInfo(fi)
		}
	} else {
		return fmt.Errorf("furni not found")
	}

	return
}
