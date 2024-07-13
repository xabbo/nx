package imager

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"io"
	"runtime"
	"sync"

	"github.com/xyproto/palgen"
)

type gifEncoder struct {
	opts GifEncoderOptions
}

type GifEncoderOptions struct {
	AlphaThreshold uint16
	Colors         int
}

func NewEncoderGIF(options ...EncoderOption) Encoder {
	opts := GifEncoderOptions{
		AlphaThreshold: 0x8000,
		Colors:         256,
	}
	for _, configure := range options {
		configure(&opts)
	}
	return gifEncoder{
		opts: opts,
	}
}

func (e gifEncoder) EncodeImage(w io.Writer, frame image.Image) error {
	return e.EncodeImages(w, []image.Image{frame})
}

func (g gifEncoder) EncodeImages(w io.Writer, frames []image.Image) (err error) {
	colors := make([]color.Color, 0)
	for _, img := range frames {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
				col := img.At(x, y)
				_, _, _, a := col.RGBA()
				if a >= uint32(g.opts.AlphaThreshold) {
					colors = append(colors, img.At(x, y))
				}
			}
		}
	}

	globalPalette, err := palgen.Generate(paletteImg(colors), g.opts.Colors-1)
	if err != nil {
		return
	}
	globalPalette = append(globalPalette, color.Transparent)

	delays := make([]int, 0, len(frames))
	disposals := make([]byte, 0, len(frames))

	wg := &sync.WaitGroup{}
	wg.Add(len(frames))

	paletteImgs := make([]*image.Paletted, len(frames))
	chImgIndex := make(chan int)
	for range runtime.NumCPU() {
		go func() {
			for i := range chImgIndex {
				bounds := frames[i].Bounds()
				bounds = bounds.Sub(bounds.Min)
				src := alphaThresholdImage{frames[i], uint32(g.opts.AlphaThreshold)}
				img := image.NewPaletted(bounds, globalPalette)
				draw.Src.Draw(img, img.Bounds(), image.Transparent, image.Point{})
				draw.Over.Draw(img, bounds, src, frames[i].Bounds().Min)
				paletteImgs[i] = img
				wg.Done()
			}
		}()
	}

	for i := range frames {
		chImgIndex <- i
		delays = append(delays, 4)
		disposals = append(disposals, gif.DisposalBackground)
	}
	wg.Wait()
	close(chImgIndex)

	err = gif.EncodeAll(w, &gif.GIF{
		Image:    paletteImgs,
		Delay:    delays,
		Disposal: disposals,
	})
	return
}

func (g gifEncoder) EncodeAnimation(w io.Writer, anim Animation, seqIndex int, frameCount int) error {
	imgs := RenderFrames(anim, seqIndex, frameCount)
	return g.EncodeImages(w, imgs)
}

func (g gifEncoder) EncodeFrame(w io.Writer, anim Animation, sequenceIndex int, frameIndex int) (err error) {
	img := RenderFrame(anim, sequenceIndex, frameIndex)
	return g.EncodeImages(w, []image.Image{img})
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
