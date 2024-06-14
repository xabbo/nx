package furni

import (
	"image/png"
	"os"

	"github.com/spf13/cobra"

	root "xabbo.b7c.io/nx/cmd/nx/cmd"
	"xabbo.b7c.io/nx/cmd/nx/spinner"
	gd "xabbo.b7c.io/nx/gamedata"
	"xabbo.b7c.io/nx/render"
)

var opts struct {
	swfPath    string
	states     bool
	identifier string
	dir        int
}

var Cmd = &cobra.Command{
	Use:  "furni [flags] identifier",
	Args: cobra.MaximumNArgs(1),
	RunE: run,
}

func init() {
	f := Cmd.Flags()

	f.StringVar(&opts.swfPath, "swf", "", "Path to a furni library in SWF format.")
	f.BoolVar(&opts.states, "states", false, "Print number of states.")
	f.StringVarP(&opts.identifier, "identifier", "i", "", "The furni identifier to load.")
	f.IntVarP(&opts.dir, "dir", "d", 0, "The direction to render.")
}

func run(cmd *cobra.Command, args []string) (err error) {

	mgr := gd.NewManager(root.Host)

	err = spinner.DoErr("Loading game data...", func() (err error) {
		return mgr.Load(gd.GameDataVariables, gd.GameDataFurni)
	})

	err = spinner.DoErr("Loading furni library...", func() (err error) {
		return mgr.LoadFurni(opts.identifier)
	})
	if err != nil {
		return
	}

	renderer := render.NewFurniRenderer(mgr)
	anim, err := renderer.Render(render.Furni{
		Identifier: opts.identifier,
		Size:       64,
		Direction:  opts.dir,
	})
	if err != nil {
		return
	}

	cmd.Printf("%+v\n", anim.Frames[0])
	img := anim.Frames[0].ToImage()
	f, err := os.OpenFile("test.png", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return
	}
	defer f.Close()

	err = png.Encode(f, img)
	return
}
