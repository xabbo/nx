package info

import (
	"fmt"

	"github.com/spf13/cobra"

	gd "xabbo.b7c.io/nx/gamedata"

	_root "xabbo.b7c.io/nx/cmd/nx/cmd"
	"xabbo.b7c.io/nx/cmd/nx/cmd/furni"
	"xabbo.b7c.io/nx/cmd/nx/util"
)

var Cmd = &cobra.Command{
	Use:  "info",
	Args: cobra.MaximumNArgs(1),
	RunE: runInfo,
}

var (
	kind       int
	identifier string
)

func init() {
	f := Cmd.Flags()
	f.IntVarP(&kind, "kind", "k", 0, "The furni kind (type ID)")
	f.StringVarP(&identifier, "identifier", "i", "", "The furni identifier (class name)")

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
			fi = &furniInfo
		}
	} else if kind > 0 {
		for _, f := range mgr.Furni() {
			if f.Kind == kind {
				fi = &f
				break
			}
		}
	} else if identifier != "" {
		for _, f := range mgr.Furni() {
			if f.Identifier == identifier {
				fi = &f
				break
			}
		}
	}

	if fi != nil {
		util.RenderFurniInfo(fi)
	} else {
		return fmt.Errorf("furni not found")
	}

	return
}
