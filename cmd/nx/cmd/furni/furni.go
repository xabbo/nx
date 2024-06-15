package furni

import (
	"fmt"
	"slices"

	"github.com/spf13/cobra"

	gd "xabbo.b7c.io/nx/gamedata"

	root "xabbo.b7c.io/nx/cmd/nx/cmd"
	"xabbo.b7c.io/nx/cmd/nx/util"
)

var furniCmd = &cobra.Command{
	Use:  "furni",
	RunE: runFurni,
}

var opts struct {
	listSwitch util.MutexValue
}

func init() {
	f := furniCmd.Flags()
	opts.listSwitch.Switch(f, "lines", "List furni lines")
	opts.listSwitch.Switch(f, "categories", "List furni categories")
	opts.listSwitch.Switch(f, "environments", "List furni environments")

	root.Cmd.AddCommand(furniCmd)
}

func runFurni(cmd *cobra.Command, args []string) (err error) {
	if opts.listSwitch.Selected() == "" {
		return fmt.Errorf("no options specified")
	}

	mgr := gd.NewManager(root.Host)
	err = util.LoadFurni(mgr)
	if err != nil {
		return
	}

	furnis := make([]gd.FurniInfo, len(mgr.Furni()))
	for _, furni := range mgr.Furni() {
		furnis = append(furnis, furni)
	}

	switch opts.listSwitch.Selected() {
	case "lines":
		listDistinctBy(furnis, getLine)
	case "categories":
		listDistinctBy(furnis, getCategory)
	case "environments":
		listDistinctBy(furnis, getEnvironment)
	}

	return
}

func listDistinctBy[T any](items []T, get func(T) string) {
	values := distinctBy(items, get)
	slices.Sort(values)
	for _, value := range values {
		fmt.Println(value)
	}
}

func distinctBy[T any](items []T, get func(T) string) []string {
	known := make(map[string]struct{})
	distinct := make([]string, 0)
	for _, it := range items {
		value := get(it)
		if _, exist := known[value]; !exist {
			known[value] = struct{}{}
			distinct = append(distinct, value)
		}
	}
	return distinct
}

func getName(fi gd.FurniInfo) string {
	return fi.Name
}

func getIdentifier(fi gd.FurniInfo) string {
	return fi.Identifier
}

func getLine(fi gd.FurniInfo) string {
	return fi.Line
}

func getCategory(fi gd.FurniInfo) string {
	return fi.Category
}

func getEnvironment(fi gd.FurniInfo) string {
	return fi.Environment
}

// Returns the number of boolean flags that are true.
func nFlags(bools ...bool) int {
	n := 0
	for _, b := range bools {
		if b {
			n++
		}
	}
	return n
}
