package furni

import (
	"errors"
	"fmt"
	"image"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"b7c.io/swfx"

	gd "xabbo.b7c.io/nx/gamedata"
	"xabbo.b7c.io/nx/imager"
	"xabbo.b7c.io/nx/raw/nitro"
	"xabbo.b7c.io/nx/res"

	_root "xabbo.b7c.io/nx/cmd/nx/cmd"
	_parent "xabbo.b7c.io/nx/cmd/nx/cmd/render"
	"xabbo.b7c.io/nx/cmd/nx/spinner"
)

var Cmd = &cobra.Command{
	Use:  "furni [flags] [identifier]",
	Args: cobra.MaximumNArgs(1),
	RunE: run,
}

var opts struct {
	inputFilePath  string
	size           int // visualization size
	dir            int // furni direction
	state          int
	seq            int
	colors         int
	verbose        bool
	format         string
	fullSequence   bool
	alphaThreshold float64
}

func init() {
	f := Cmd.Flags()

	f.StringVarP(&opts.inputFilePath, "input", "i", "", "Path to a furni library in SWF format.")
	f.IntVarP(&opts.dir, "dir", "d", 0, "The direction.")
	f.IntVar(&opts.size, "size", 64, "The visualization size.")
	f.IntVarP(&opts.state, "state", "s", 0, "The animation state.")
	f.IntVar(&opts.seq, "seq", 0, "The animation sequence index.")
	f.IntVar(&opts.colors, "colors", 256, "Number of colors to quantize when encoding to GIF.")
	f.BoolVarP(&opts.verbose, "verbose", "v", false, "Output detailed information.")
	f.BoolVar(&opts.fullSequence, "full-sequence", false, "Render the full animation sequence.")
	f.Float64Var(&opts.alphaThreshold, "alpha-threshold", 0, "Alpha threshold for GIF encoding.")

	f.StringVarP(&opts.format, "format", "f", "png", "Output image format. (apng, png, gif, svg)")

	_parent.Cmd.AddCommand(Cmd)
}

func run(cmd *cobra.Command, args []string) (err error) {

	spinner.Start()
	defer spinner.Stop()

	var libraryName string

	mgr := gd.NewManager(_root.Host)

	if opts.inputFilePath != "" {
		if len(args) > 0 {
			return errors.New("only one of furni identifier or input file may be specified")
		}
		cmd.SilenceUsage = true
		switch {
		case strings.HasSuffix(opts.inputFilePath, ".swf"):
			spinner.Message("Loading SWF library...")

			var swf *swfx.Swf
			swf, err = loadSwfFile(opts.inputFilePath)
			if err != nil {
				return
			}

			var lib res.AssetLibrary
			lib, err = res.LoadFurniLibrarySwf(swf)
			if err != nil {
				return
			}

			mgr.AddLibrary(lib)
			libraryName = lib.Name()
		case strings.HasSuffix(opts.inputFilePath, ".nitro"):
			spinner.Message("Loading Nitro library...")

			var archive nitro.Archive
			archive, err = loadNitroArchive(opts.inputFilePath)
			if err != nil {
				return
			}

			var lib res.AssetLibrary
			lib, err = res.LoadFurniLibraryNitro(archive)
			if err != nil {
				return
			}

			mgr.AddLibrary(lib)
			libraryName = lib.Name()
		default:
			return fmt.Errorf("input file format not supported")
		}
	} else {
		if len(args) != 1 {
			return errors.New("no furni identifier or input file specified")
		}
		cmd.SilenceUsage = true

		libraryName = args[0]

		spinner.Message("Loading game data...")
		err = mgr.Load(gd.GameDataVariables, gd.GameDataFurni)
		if err != nil {
			return
		}

		spinner.Message("Loading furni library...")
		err = mgr.LoadFurni(libraryName)
		if err != nil {
			return
		}

		libraryName = strings.Split(libraryName, "*")[0]
	}

	imgr := imager.NewFurniImager(mgr)
	anim, err := imgr.Compose(imager.Furni{
		Identifier: libraryName,
		Size:       opts.size,
		Direction:  opts.dir,
		State:      opts.state,
	})
	if err != nil {
		return
	}

	if opts.verbose {
		spinner.Printf("Total animation frames: %d\n", anim.TotalFrames(opts.seq))
	}

	if opts.format == "svg" {
		outName := fmt.Sprintf("%s_%d_%d_%d_%d.%d",
			libraryName, opts.size, opts.dir, opts.state, opts.seq, 0)

		var f *os.File
		f, err = os.OpenFile(outName+".svg", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		encoder := imager.NewEncoderSVG()
		err = encoder.EncodeFrame(f, anim, opts.seq, 0)
		return nil
	}

	start := time.Now()
	spinner.Message(fmt.Sprintf("Drawing frames with %d cores...", runtime.NumCPU()))

	var renderFrameCount int
	switch opts.format {
	case "apng", "gif":
		if opts.fullSequence {
			renderFrameCount = anim.TotalFrames(opts.seq)
		} else {
			renderFrameCount = anim.LongestFrameSequence(opts.seq)
		}
	default:
		renderFrameCount = 1
	}

	outName := fmt.Sprintf("%s_%d_%d_%d_%d.%d",
		libraryName, opts.size, opts.dir, opts.state, opts.seq, renderFrameCount)

	imgs := imager.RenderFrames(anim, opts.seq, renderFrameCount)
	if opts.verbose {
		spinner.Printf("Rendered %d frame(s) in %dms\n", len(imgs), time.Since(start).Milliseconds())
	}

	switch opts.format {
	case "png":
		encoder := imager.NewEncoderPNG()
		saveEncoded(outName+".png", encoder, imgs)
	case "apng":
		encoder := imager.NewEncoderAPNG()
		saveEncoded(outName+".apng", encoder, imgs)
	case "gif":
		threshold := uint16(opts.alphaThreshold * 0xffff)
		encoder := imager.NewEncoderGIF(imager.WithAlphaThreshold(threshold))
		saveEncoded(outName+".gif", encoder, imgs)
	}

	return
}

func loadSwfFile(filePath string) (swf *swfx.Swf, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()
	return swfx.ReadSwf(f)
}

func loadNitroArchive(filePath string) (archive nitro.Archive, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()
	r := nitro.NewReader(f)
	return r.ReadArchive()
}

func saveEncoded(name string, encoder imager.ImageEncoder, imgs []image.Image) (err error) {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	return encoder.EncodeImages(f, imgs)
}
