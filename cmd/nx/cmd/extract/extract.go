package extract

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	root "xabbo.io/nx/cmd/nx/cmd"
	"xabbo.io/nx/raw/nitro"
)

var Cmd = &cobra.Command{
	Use:   "extract [files...]",
	Short: "Extracts files. Currently supports Nitro archives.",
	RunE:  run,
}

func init() {
	root.Cmd.AddCommand(Cmd)
}

func run(cmd *cobra.Command, args []string) (err error) {
	for _, file := range args {
		err := extractFile(file)
		if err != nil {
			cmd.PrintErrf("%s: %s\n", file, err)
		}
	}
	return
}

func extractFile(name string) (err error) {
	switch {
	case strings.HasSuffix(name, ".nitro"):
		err = extractNitro(name)
	default:
		err = fmt.Errorf("unknown file format")
	}
	return
}

func extractNitro(name string) (err error) {
	dir := strings.TrimSuffix(name, filepath.Ext(name))
	if dir == name {
		dir += "_extracted"
	}

	err = os.Mkdir(dir, 0755)
	if err != nil {
		return
	}

	archive, err := readNitroArchive(name)

	for name, file := range archive.Files {
		fmt.Printf("%s: ", name)
		err = os.WriteFile(filepath.Join(dir, name), file.Data, 0644)
		if err != nil {
			return
		}
		fmt.Println("ok")
	}

	return
}

func readNitroArchive(name string) (archive nitro.Archive, err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()

	return nitro.NewReader(f).ReadArchive()
}
