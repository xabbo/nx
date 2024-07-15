package imager

import (
	"image"
	"image/color"
	"image/draw"
	"runtime"
	"slices"
	"sync"

	"golang.org/x/exp/maps"

	"xabbo.b7c.io/nx/res"
)

func (frame Frame) Draw(canvas draw.Image, drawer draw.Drawer) {
	for _, sprite := range frame {
		sprite.Draw(canvas, drawer)
	}
}

func (frame Frame) ToImage() image.Image {
	canvas := image.NewRGBA(frame.Bounds())
	frame.Draw(canvas, draw.Over)
	return canvas
}

// Gets all assets used in this animation for the specified frame sequence.
func (anim Animation) RequiredAssets(seqIndex int) []*res.Asset {
	m := map[*res.Asset]struct{}{}
	for _, layer := range anim.Layers {
		var seq res.FrameSequence
		if seqIndex < len(layer.Sequences) {
			seq = layer.Sequences[seqIndex]
		} else if len(layer.Sequences) > 0 {
			seq = layer.Sequences[0]
		} else {
			seq = []int{0}
		}
		for _, frameId := range seq {
			for _, sprite := range layer.Frames[frameId] {
				if sprite.Asset != nil {
					m[sprite.Asset] = struct{}{}
				}
			}
		}
	}
	return maps.Keys(m)
}

func RenderFramesBounds(bounds image.Rectangle, anim Animation, seqIndex, frameCount int) []image.Image {
	frames := make([]image.Image, frameCount)

	wg := sync.WaitGroup{}
	wg.Add(frameCount)

	ch := make(chan int)
	defer close(ch)
	for range runtime.NumCPU() {
		go func() {
			for frameIndex := range ch {
				img := image.NewRGBA(bounds)
				if anim.Background != nil && anim.Background != color.Transparent {
					draw.Over.Draw(img, img.Bounds(), image.NewUniform(anim.Background), image.Point{})
				}
				DrawFrame(anim, img, nil, seqIndex, frameIndex)
				frames[frameIndex] = img
				wg.Done()
			}
		}()
	}

	for frameIndex := range frames {
		ch <- frameIndex
	}
	wg.Wait()
	return frames
}

func RenderFrames(anim Animation, seqIndex, frameCount int) []image.Image {
	bounds := anim.Bounds(seqIndex)
	return RenderFramesBounds(bounds, anim, seqIndex, frameCount)
}

func RenderQuantizedFrames(anim Animation, seqIndex int, palette color.Palette, count int) []*image.Paletted {
	frames := make([]*image.Paletted, count)
	bounds := anim.Bounds(seqIndex)

	wg := sync.WaitGroup{}
	wg.Add(count)

	ch := make(chan int)
	defer close(ch)
	for range runtime.NumCPU() {
		go func() {
			for frameIndex := range ch {
				img := image.NewPaletted(bounds, palette)
				draw.Src.Draw(img, bounds, image.Transparent, image.Point{})
				DrawFrame(anim, img, nil, seqIndex, frameIndex)
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

func DrawFrame(anim Animation, canvas draw.Image, drawer draw.Drawer, sequenceIndex int, frameIndex int) {
	layerIds := maps.Keys(anim.Layers)
	slices.SortFunc(layerIds, func(a, b int) int {
		// Shadow layer index should be -1.
		// Ensure that the shadow is always on the bottom layer.
		if a < 0 {
			return -1
		} else if b < 0 {
			return 1
		}
		la := anim.Layers[a]
		lb := anim.Layers[b]
		diff := la.Z - lb.Z
		// If both layer Z indexes are equal, order by layer ID.
		if diff == 0 {
			diff = a - b
		}
		return diff
	})
	for _, layerId := range layerIds {
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

func RenderFrame(anim Animation, seqIndex int, frameIndex int) image.Image {
	canvas := image.NewRGBA(anim.Bounds(seqIndex))
	if anim.Background != nil && anim.Background != color.Transparent {
		draw.Over.Draw(canvas, canvas.Bounds(), image.NewUniform(anim.Background), image.Point{})
	}
	DrawFrame(anim, canvas, nil, seqIndex, frameIndex)
	return canvas
}
