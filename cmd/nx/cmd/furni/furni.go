package furni

import (
	"fmt"
	"slices"

	"github.com/spf13/cobra"

	"github.com/b7c/nx"

	root "cli/cmd"
	"cli/util"
)

var furniCmd = &cobra.Command{
	Use:  "furni",
	RunE: runFurni,
}

var (
	listSwitch util.MutexValue
)

func init() {
	f := furniCmd.Flags()
	listSwitch.Switch(f, "lines", "List furni lines")
	listSwitch.Switch(f, "categories", "List furni categories")
	listSwitch.Switch(f, "environments", "List furni environments")

	root.Cmd.AddCommand(furniCmd)
}

func runFurni(cmd *cobra.Command, args []string) (err error) {
	if listSwitch.Selected() == "" {
		return fmt.Errorf("no options specified")
	}

	mgr := nx.NewGamedataManager(root.Host)
	err = util.LoadFurni(mgr)
	if err != nil {
		return
	}

	furnis := make([]nx.FurniInfo, len(mgr.Furni))
	for _, furni := range mgr.Furni {
		furnis = append(furnis, furni)
	}

	switch listSwitch.Selected() {
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

func getName(fi nx.FurniInfo) string {
	return fi.Name
}

func getIdentifier(fi nx.FurniInfo) string {
	return fi.Identifier
}

func getLine(fi nx.FurniInfo) string {
	return fi.Line
}

func getCategory(fi nx.FurniInfo) string {
	return fi.Category
}

func getEnvironment(fi nx.FurniInfo) string {
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
