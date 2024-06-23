package furni

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/kettek/apng"
	"github.com/spf13/cobra"
	"github.com/xyproto/palgen"

	"b7c.io/swfx"

	gd "xabbo.b7c.io/nx/gamedata"
	"xabbo.b7c.io/nx/render"
	"xabbo.b7c.io/nx/res"

	_root "xabbo.b7c.io/nx/cmd/nx/cmd"
	_parent "xabbo.b7c.io/nx/cmd/nx/cmd/render"
	"xabbo.b7c.io/nx/cmd/nx/spinner"
)

const alphaThreshold = 0x8000

var Cmd = &cobra.Command{
	Use:  "furni [flags] [identifier]",
	Args: cobra.MaximumNArgs(1),
	RunE: run,
}

var opts struct {
	inputFilePath string
	size          int // visualization size
	dir           int // furni direction
	state         int
	seq           int
	colors        int
	verbose       bool
	format        string
	fullSequence  bool
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

	f.StringVarP(&opts.format, "format", "f", "apng", "Output image format. (apng, png, gif)")

	_parent.Cmd.AddCommand(Cmd)
}

type alphaThresholdImage struct {
	img       image.Image
	threshold uint32
}

func (i alphaThresholdImage) ColorModel() color.Model {
	return i.img.ColorModel()
}

func (i alphaThresholdImage) Bounds() image.Rectangle {
	return i.img.Bounds()
}

func (i alphaThresholdImage) At(x, y int) color.Color {
	c := i.img.At(x, y)
	switch c.(type) {
	default:
		_, _, _, a := c.RGBA()
		if a >= i.threshold {
			return c
		} else {
			return color.Transparent
		}
	}
}

func run(cmd *cobra.Command, args []string) (err error) {

	spinner.Start()
	defer spinner.Stop()

	var furniIdentifier string

	mgr := gd.NewManager(_root.Host)

	if opts.inputFilePath != "" {
		if len(args) > 0 {
			return errors.New("only one of furni identifier or input file may be specified")
		}
		cmd.SilenceUsage = true
		if strings.HasSuffix(opts.inputFilePath, ".swf") {
			spinner.Message("Loading SWF library...")

			var swf *swfx.Swf
			swf, err = loadSwfFile(opts.inputFilePath)
			if err != nil {
				return err
			}

			var lib res.AssetLibrary
			lib, err = res.LoadFurniLibrarySwf(swf)
			if err != nil {
				return
			}

			mgr.AddLibrary(lib)
			furniIdentifier = lib.Name()
		} else {
			return fmt.Errorf("input file format not supported")
		}
	} else {
		if len(args) != 1 {
			return errors.New("no furni identifier or input file specified")
		}
		cmd.SilenceUsage = true

		furniIdentifier = args[0]

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
	}

	renderer := render.NewFurniRenderer(mgr)
	anim, err := renderer.Render(render.Furni{
		Identifier: furniIdentifier,
		Size:       opts.size,
		Direction:  opts.dir,
		State:      opts.state,
	})
	if err != nil {
		return
	}

	if opts.verbose {
		spinner.Printf("Total animation frames: %d\n", anim.TotalFrames())
	}

	start := time.Now()
	spinner.Message(fmt.Sprintf("Drawing frames with %d cores...", runtime.NumCPU()))

	var renderFrameCount int
	if opts.fullSequence {
		renderFrameCount = anim.TotalFrames()
	} else {
		renderFrameCount = anim.LongestFrameSequence(opts.seq)
	}
	imgs := anim.RenderFrames(opts.seq, renderFrameCount) // anim.DrawQuantizedFrames(opts.seq, colorPalette, renderFrameCount)
	if opts.verbose {
		spinner.Printf("Rendered %d frames in %dms\n", len(imgs), time.Since(start).Milliseconds())
	}

	switch opts.format {
	case "apng":
		err = saveAPNG(imgs)
	case "gif":
		err = saveGIF(imgs)
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

func saveAPNG(imgs []image.Image) (err error) {
	spinner.Message("Encoding image...")
	start := time.Now()

	a := apng.APNG{
		Frames: make([]apng.Frame, len(imgs)),
	}

	for i := range imgs {
		a.Frames[i].Image = imgs[i]
		a.Frames[i].DelayNumerator = 1
		a.Frames[i].DelayDenominator = 30
	}

	f, err := os.Create("out.png")
	if err != nil {
		return
	}
	defer f.Close()

	err = apng.Encode(f, a)
	if err != nil {
		return
	}

	if opts.verbose {
		spinner.Printf("Encoded image in %dms.\n", time.Since(start).Milliseconds())
	}
	return
}

func saveGIF(imgs []image.Image) (err error) {
	colors := make([]color.Color, 0)
	for _, img := range imgs {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
				col := img.At(x, y)
				_, _, _, a := col.RGBA()
				if a >= 0x8000 {
					colors = append(colors, img.At(x, y))
				}
			}
		}
	}

	globalPalette, err := palgen.Generate(paletteImg(colors), 255)
	if err != nil {
		return
	}
	globalPalette = append(globalPalette, color.Transparent)

	delays := make([]int, 0, len(imgs))
	disposals := make([]byte, 0, len(imgs))

	spinner.Message(fmt.Sprintf("Drawing quantized frames with %d cores...", runtime.NumCPU()))
	start := time.Now()

	wg := &sync.WaitGroup{}
	wg.Add(len(imgs))

	paletteImgs := make([]*image.Paletted, len(imgs))
	chImgIndex := make(chan int)
	for range runtime.NumCPU() {
		go func() {
			for i := range chImgIndex {
				bounds := imgs[i].Bounds()
				bounds = bounds.Sub(bounds.Min)
				src := alphaThresholdImage{imgs[i], alphaThreshold}
				img := image.NewPaletted(bounds, globalPalette)
				draw.Src.Draw(img, img.Bounds(), image.Transparent, image.Point{})
				draw.Over.Draw(img, bounds, src, imgs[i].Bounds().Min)
				paletteImgs[i] = img
				wg.Done()
			}
		}()
	}

	for i := range imgs {
		chImgIndex <- i
		delays = append(delays, 3)
		disposals = append(disposals, gif.DisposalPrevious)
	}
	wg.Wait()
	close(chImgIndex)

	if opts.verbose {
		spinner.Printf("Rendered quantized frames in %dms\n", time.Since(start).Milliseconds())
	}

	f, err := os.Create("test.gif")
	if err != nil {
		return
	}
	defer f.Close()

	spinner.Message("Encoding image")
	start = time.Now()

	err = gif.EncodeAll(f, &gif.GIF{
		Image:    paletteImgs,
		Delay:    delays,
		Disposal: disposals,
	})
	if err != nil {
		return
	}

	if opts.verbose {
		spinner.Printf("Encoded image in %dms\n", time.Since(start).Milliseconds())
	}

	return
}

type paletteImg color.Palette

func (p paletteImg) ColorModel() color.Model {
	return color.RGBAModel
}

func (p paletteImg) Bounds() image.Rectangle {
	return image.Rect(0, 0, len(p), 1)
}

func (p paletteImg) At(x, y int) color.Color {
	return p[x]
}
