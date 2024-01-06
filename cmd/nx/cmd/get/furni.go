package get

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/b7c/nx"

	root "cli/cmd"
	"cli/spinner"
	"cli/util"
)

var ErrNotFound = errors.New("not found")

var getFurniCmd = &cobra.Command{
	Use:  "furni <identifier>",
	RunE: runGetFurni,
}

func init() {
	getCmd.AddCommand(getFurniCmd)
}

func runGetFurni(cmd *cobra.Command, args []string) (err error) {
	if len(args) == 0 {
		return fmt.Errorf("no furni identifier specified")
	}
	cmd.SilenceUsage = true

	mgr := nx.NewGamedataManager(root.Host)
	err = util.LoadFurni(mgr)
	if err != nil {
		return
	}

	defer spinner.Stop()

	for _, identifier := range args {
		if furni, ok := mgr.Furni[identifier]; ok {
			err := downloadFurni(&furni)
			if err != nil {
				return fmt.Errorf("failed to get %d/%s: %s", furni.Revision, furni.Identifier, err)
			}
		} else {
			cmd.PrintErrf("%s: identifier not found\n", identifier)
		}
	}

	return nil
}

func downloadFurni(fi *nx.FurniInfo) (err error) {
	defer spinner.Stop()

	identifier := fi.Identifier
	idx := strings.Index(identifier, "*")
	if idx > 0 {
		identifier = identifier[:idx]
	}
	filePath := identifier + ".swf"

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0755)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			spinner.Stop()
			fmt.Printf("%s: file exists\n", filePath)
			return nil
		}
		return
	}

	success := false
	defer func() {
		f.Close()
		if !success {
			os.Remove(filePath)
		}
	}()

	spinner.Message(fmt.Sprintf(
		"Downloading %d/%s.swf...",
		fi.Revision, identifier,
	))
	spinner.Start()

	res, err := http.Get(fmt.Sprintf("https://images.habbo.com/dcr/hof_furni/%d/%s.swf",
		fi.Revision, identifier))
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return ErrNotFound
		}
		return fmt.Errorf("server responded %s", res.Status)
	}

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return
	}

	success = true
	spinner.Stop()
	fmt.Println(filePath)
	return
}
