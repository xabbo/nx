package render

import (
	"image"
	"image/color"
	"image/draw"
	"runtime"
	"slices"
	"sync"

	"github.com/disintegration/imaging"
	"golang.org/x/exp/maps"

	"xabbo.b7c.io/nx/res"
)

func (frame Frame) Draw(canvas draw.Image, drawer draw.Drawer) {
	for _, layer := range frame {
		for _, sprite := range layer.Sprites {
			srcImg := sprite.Image()
			if srcImg == nil {
				continue
			}
			bounds := srcImg.Bounds()
			offset := sprite.Offset
			if sprite.FlipH {
				offset.X = offset.X*-1 + srcImg.Bounds().Dx() - 64
				srcImg = imaging.FlipH(srcImg)
			}
			drawer.Draw(canvas, bounds.Sub(offset), srcImg, image.Point{})
		}
	}
}

func (frame Frame) ToImage() image.Image {
	canvas := image.NewRGBA(frame.Bounds())
	frame.Draw(canvas, draw.Over)
	return canvas
}

// Gets all assets used in this animation for the specified frame sequence.
func (anim Animation) Assets(sequenceIdx int) []*res.Asset {
	m := map[*res.Asset]struct{}{}
	for _, layer := range anim.Layers {
		var seq res.FrameSequence
		if sequenceIdx < len(layer.Sequences) {
			seq = layer.Sequences[sequenceIdx]
		} else if len(layer.Sequences) > 0 {
			seq = layer.Sequences[0]
		} else {
			seq = []int{0}
		}
		for _, frameId := range seq {
			for _, frameLayer := range layer.Frames[frameId] {
				for _, sprite := range frameLayer.Sprites {
					if sprite.Asset != nil {
						m[sprite.Asset] = struct{}{}
					}
				}
			}
		}
	}
	return maps.Keys(m)
}

func (anim Animation) RenderFrames(sequenceIndex, frameCount int) []image.Image {
	frames := []image.Image{}
	bounds := anim.Bounds()

	wg := sync.WaitGroup{}
	wg.Add(frameCount)

	ch := make(chan int)
	defer close(ch)
	for range runtime.NumCPU() {
		go func() {
			for frameIndex := range ch {
				img := image.NewRGBA(bounds)
				anim.drawFrame(img, draw.Over, sequenceIndex, frameIndex)
				frames[frameIndex] = img
				wg.Done()
			}
		}()
	}

	for frameIndex := range frameCount {
		frames = append(frames, anim.renderFrame(bounds, sequenceIndex, frameIndex))
	}
	return frames
}

func (anim Animation) DrawQuantizedFrames(sequenceIndex int, palette color.Palette, count int) []*image.Paletted {
	frames := make([]*image.Paletted, count)
	bounds := anim.Bounds()

	wg := sync.WaitGroup{}
	wg.Add(count)

	ch := make(chan int)
	defer close(ch)
	for range runtime.NumCPU() {
		go func() {
			for frameIndex := range ch {
				img := image.NewPaletted(bounds, palette)
				draw.Src.Draw(img, bounds, image.Transparent, image.Point{})
				anim.drawFrame(img, draw.Over, sequenceIndex, frameIndex)
				frames[frameIndex] = img
				wg.Done()
			}
		}()
	}
	for frameIndex := range count {
		ch <- frameIndex
	}
	wg.Wait()
	return frames
}

func (anim Animation) drawFrame(canvas draw.Image, drawer draw.Drawer, sequenceIndex int, frameIndex int) {
	layerIds := maps.Keys(anim.Layers)
	slices.Sort(layerIds)
	for layerId := range layerIds {
		layer := anim.Layers[layerId]
		var seq res.FrameSequence
		if sequenceIndex < len(layer.Sequences) {
			seq = layer.Sequences[sequenceIndex]
		} else {
			if len(layer.Sequences) > 0 {
				seq = layer.Sequences[0]
			} else {
				seq = []int{0}
			}
		}
		frameId := seq[(frameIndex/max(1, layer.FrameRepeat))%len(seq)]
		layer.Frames[frameId].Draw(canvas, drawer)
	}
}

func (anim Animation) renderFrame(bounds image.Rectangle, sequenceIndex int, frameIndex int) image.Image {
	canvas := image.NewRGBA(bounds)
	anim.drawFrame(canvas, draw.Over, sequenceIndex, frameIndex)
	return canvas
}
