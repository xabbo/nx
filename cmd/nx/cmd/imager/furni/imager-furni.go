package furni

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"

	"b7c.io/swfx"

	gd "xabbo.b7c.io/nx/gamedata"
	"xabbo.b7c.io/nx/imager"
	"xabbo.b7c.io/nx/raw/nitro"
	"xabbo.b7c.io/nx/res"

	_root "xabbo.b7c.io/nx/cmd/nx/cmd"
	_parent "xabbo.b7c.io/nx/cmd/nx/cmd/imager"
	"xabbo.b7c.io/nx/cmd/nx/spinner"
)

var Cmd = &cobra.Command{
	Use:  "furni [flags] [identifier]",
	Args: cobra.MaximumNArgs(1),
	RunE: run,
}

var (
	lib  res.FurniLibrary
	opts struct {
		inputFilePath  string
		size           int // visualization size
		dir            int // furni direction
		state          int
		seq            int
		color          int
		colors         int
		verbose        bool
		format         string
		fullSequence   bool
		alphaThreshold float64
		allDirections  bool
		allStates      bool
		allColors      bool
	}
)

func init() {
	f := Cmd.Flags()

	f.StringVarP(&opts.inputFilePath, "input", "i", "", "Path to a furni library in SWF format.")
	f.IntVarP(&opts.dir, "dir", "d", 0, "The direction.")
	f.IntVar(&opts.size, "size", 64, "The visualization size.")
	f.IntVarP(&opts.state, "state", "s", 0, "The animation state.")
	f.IntVar(&opts.seq, "seq", 0, "The animation sequence index.")
	f.IntVarP(&opts.color, "color", "c", 0, "The color index to use.")
	f.IntVar(&opts.colors, "num-colors", 256, "Number of colors to quantize when encoding to GIF.")
	f.BoolVarP(&opts.verbose, "verbose", "v", false, "Output detailed information.")
	f.BoolVar(&opts.fullSequence, "full-sequence", false, "Render the full animation sequence.")
	f.Float64Var(&opts.alphaThreshold, "alpha-threshold", 0, "Alpha threshold for GIF encoding.")
	f.BoolVarP(&opts.allDirections, "dirs", "D", false, "Output all directions.")
	f.BoolVarP(&opts.allStates, "states", "S", false, "Output all states.")
	f.BoolVarP(&opts.allColors, "colors", "C", false, "Output all colors.")

	f.StringVarP(&opts.format, "format", "f", "png", "Output image format. (apng, png, gif, svg)")

	_parent.Cmd.AddCommand(Cmd)
}

func run(cmd *cobra.Command, args []string) (err error) {

	spinner.Start()
	defer spinner.Stop()

	mgr := gd.NewManager(_root.Host)

	if opts.inputFilePath != "" {
		if len(args) > 0 {
			return errors.New("only one of furni identifier or input file may be specified")
		}

		cmd.SilenceUsage = true
		lib, err = loadLibraryFile(opts.inputFilePath)
		if err != nil {
			return
		}

		mgr.AddLibrary(lib)
	} else {
		if len(args) != 1 {
			return errors.New("no furni identifier or input file specified")
		}
		cmd.SilenceUsage = true

		furniIdentifier := args[0]
		split := strings.SplitN(furniIdentifier, "*", 2)
		libName := split[0]
		if len(split) >= 2 {
			strColorIndex := split[1]
			if colorIndex, err := strconv.Atoi(strColorIndex); err == nil {
				if !cmd.Flags().Lookup("color").Changed {
					opts.color = colorIndex
				}
			}
		}

		spinner.Message("Loading game data...")
		err = mgr.Load(gd.GameDataVariables, gd.GameDataFurni)
		if err != nil {
			return
		}

		spinner.Message("Loading furni library...")
		err = mgr.LoadFurni(furniIdentifier)
		if err != nil {
			return
		}

		var ok bool
		if lib, ok = mgr.Library(libName).(res.FurniLibrary); !ok {
			err = fmt.Errorf("failed to load furni library")
			return
		}
	}

	vis, ok := lib.Visualizations()[opts.size]
	if !ok {
		err = fmt.Errorf("no visualization for size: %d", opts.size)
		return
	}

	directions := []int{}
	if opts.allDirections {
		directions = maps.Keys(vis.Directions)
	} else {
		dir := opts.dir
		if !cmd.Flags().Lookup("dir").Changed {
			if vis, ok := lib.Visualizations()[opts.size]; ok {
				for i := range 4 {
					d := (2 + i*2) % 8
					if _, ok := vis.Directions[d]; ok {
						dir = d
						break
					}
				}
			}
		}
		directions = append(directions, dir)
	}

	states := []int{}
	if opts.allStates {
		states = maps.Keys(vis.Animations)
		slices.Sort(states)
		if len(states) == 0 {
			states = append(states, 0)
		}
	} else {
		states = append(states, opts.state)
	}

	colors := []int{}
	if opts.allColors {
		colors = maps.Keys(vis.Colors)
		slices.Sort(colors)
		if len(colors) == 0 {
			colors = append(colors, 0)
		}
	} else {
		colors = append(colors, opts.color)
	}

	imgr := imager.NewFurniImager(mgr)

	for _, dir := range directions {
		for _, state := range states {
			for _, color := range colors {
				furni := imager.Furni{
					Identifier: lib.Name(),
					Size:       opts.size,
					Direction:  dir,
					State:      state,
					Color:      color,
				}

				var anim imager.Animation
				anim, err = imgr.Compose(furni)
				if err != nil {
					return
				}
				if len(anim.Layers) == 0 {
					continue
				}

				var name string
				name, err = saveAnimation(furni, anim, opts.seq, 0)
				if err != nil {
					return
				}
				spinner.Printf("%s\n", name)
			}
		}
	}

	return
}

func loadLibraryFile(name string) (lib res.FurniLibrary, err error) {
	switch {
	case strings.HasSuffix(strings.ToLower(name), ".swf"):
		var swf *swfx.Swf
		swf, err = loadSwf(name)
		if err != nil {
			return
		}
		lib, err = res.LoadFurniLibrarySwf(swf)
	case strings.HasSuffix(strings.ToLower(name), ".nitro"):
		var archive nitro.Archive
		archive, err = loadNitroArchive(name)
		if err != nil {
			return
		}
		lib, err = res.LoadFurniLibraryNitro(archive)
	default:
		err = fmt.Errorf("input file format not supported")
	}
	return
}

func loadSwf(filePath string) (swf *swfx.Swf, err error) {
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

func saveAnimation(furni imager.Furni, anim imager.Animation, seqIndex, frameIndex int) (name string, err error) {
	frameCount := 1
	if opts.format == "apng" || opts.format == "gif" {
		if opts.fullSequence {
			frameCount = anim.TotalFrames(opts.seq)
		} else {
			frameCount = anim.LongestFrameSequence(opts.seq)
		}
	}

	outName := fmt.Sprintf("%s_%d_%d_%d_%d_%d.%d",
		furni.Identifier, furni.Size, furni.Direction, furni.State, furni.Color, seqIndex, frameCount)

	var encoder any
	switch opts.format {
	case "png":
		encoder = imager.NewEncoderPNG()
	case "apng":
		encoder = imager.NewEncoderAPNG()
	case "gif":
		threshold := uint16(opts.alphaThreshold * 0xffff)
		encoder = imager.NewEncoderGIF(
			imager.WithAlphaThreshold(threshold),
			imager.WithColors(opts.colors),
		)
	case "svg":
		encoder = imager.NewEncoderSVG()
	}

	outName += "." + opts.format
	return outName, saveEncoder(outName, encoder, anim, seqIndex, frameIndex, frameCount)
}

func saveEncoder(output string, encoder any, anim imager.Animation, seqIndex, frameIndex, frameCount int) (err error) {
	f, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	switch encoder := encoder.(type) {
	case imager.AnimationEncoder:
		encoder.EncodeAnimation(f, anim, seqIndex, frameCount)
	case imager.FrameEncoder:
		encoder.EncodeFrame(f, anim, seqIndex, frameIndex)
	default:
		err = fmt.Errorf("unknown encoder type: %T", encoder)
	}

	return
}
