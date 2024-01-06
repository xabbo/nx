package furni

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/b7c/nx"

	root "cli/cmd"
	"cli/util"
)

var infoCmd = &cobra.Command{
	Use:  "info",
	Args: cobra.MaximumNArgs(1),
	RunE: runInfo,
}

var (
	kind       int
	identifier string
)

func init() {
	furniCmd.AddCommand(infoCmd)

	infoCmd.Flags().IntVarP(&kind, "kind", "k", 0, "The furni kind (type ID)")
	infoCmd.Flags().StringVarP(&identifier, "identifier", "i", "", "The furni identifier (class name)")
}

func runInfo(cmd *cobra.Command, args []string) (err error) {
	cmd.SilenceUsage = true

	mgr := nx.NewGamedataManager(root.Host)
	err = util.LoadFurni(mgr)
	if err != nil {
		return
	}

	var fi *nx.FurniInfo

	if len(args) > 0 {
		identifier := args[0]
		if furniInfo, ok := mgr.Furni[identifier]; ok {
			fi = &furniInfo
		}
	} else if kind > 0 {
		for _, f := range mgr.Furni {
			if f.Kind == kind {
				fi = &f
				break
			}
		}
	} else if identifier != "" {
		for _, f := range mgr.Furni {
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
