package imager

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/disintegration/imaging"
	"github.com/phrozen/blend"
	"xabbo.b7c.io/nx/res"
)

type Sprite struct {
	Asset  *res.Asset // A reference to the asset used by this sprite.
	FlipH  bool
	FlipV  bool
	Offset image.Point
	Color  color.Color
	Blend  Blend
	Alpha  uint8
	Order  int
}

type Blend int

const (
	BlendNone Blend = iota
	BlendAdd
	BlendCopy
)

// Image gets the source image for this sprite.
func (s *Sprite) Image() image.Image {
	asset := s.Asset
	for asset.Source != nil && asset != asset.Source {
		asset = asset.Source
	}
	return asset.Image
}

// Bounds returns the relative bounds of the sprite.
func (s *Sprite) Bounds() image.Rectangle {
	img := s.Image()
	if img == nil {
		return image.Rectangle{}
	}
	return img.Bounds().Sub(s.Offset)
}

// Size returns the size of the sprite.
func (s *Sprite) Size() image.Point {
	return s.Bounds().Size()
}

// Draw draws the sprite onto the canvas using the provided drawer.
// If the drawer is nil, one will be selected automatically based on the sprite's blending mode.
func (s *Sprite) Draw(canvas draw.Image, drawer draw.Drawer) {
	srcImg := s.Image()
	if srcImg == nil {
		return
	}
	if s.Color != color.White {
		srcImg = blend.BlendNewImage(srcImg, image.NewUniform(s.Color), blend.Multiply)
	}
	bounds := srcImg.Bounds()
	offset := s.Offset
	if s.FlipH {
		srcImg = imaging.FlipH(srcImg)
	}
	if drawer == nil {
		switch s.Blend {
		case BlendAdd:
			drawer = additiveDrawer{}
		case BlendCopy:
			fallthrough
		default:
			drawer = alphaDrawer(s.Alpha)
		}
	}
	drawer.Draw(canvas, bounds.Sub(offset), srcImg, image.Point{})
}

type Frame []Sprite

func (f Frame) Bounds() (bounds image.Rectangle) {
	for _, sprite := range f {
		bounds = bounds.Union(sprite.Bounds())
	}
	return
}

type Animation struct {
	Background color.Color
	Layers     map[int]AnimationLayer
}

type AnimationLayer struct {
	Frames      map[int]Frame
	FrameRepeat int
	Sequences   []res.FrameSequence
}

// FrameSequenceOrDefault gets the specified frame sequence if it exists, or the first sequence if `i` is out of range.
// If the animation has no frame sequences, a sequence with a single zero frame is returned.
func (animationLayer AnimationLayer) FrameSequenceOrDefault(seqIndex int) res.FrameSequence {
	if seqIndex < len(animationLayer.Sequences) {
		return animationLayer.Sequences[seqIndex]
	} else if len(animationLayer.Sequences) > 0 {
		return animationLayer.Sequences[0]
	} else {
		return []int{0}
	}
}

// Bounds gets the relative bounds in the animation for the specified sequence.
func (animation Animation) Bounds(seqIndex int) (bounds image.Rectangle) {
	for _, layer := range animation.Layers {
		for _, frame := range layer.Frames {
			bounds = bounds.Union(frame.Bounds())
		}
	}
	return
}

func (animation *Animation) TotalFrames(seqIndex int) int {
	n := 1
	for _, layer := range animation.Layers {
		seq := layer.FrameSequenceOrDefault(seqIndex)
		n = lcm(n, len(seq)*max(1, layer.FrameRepeat))
	}
	return n
}

// LongestFrameSequence gets the longest frame sequence (given the specified sequence index) of all layers in the animation.
func (animation *Animation) LongestFrameSequence(seqIndex int) int {
	n := 1
	for _, layer := range animation.Layers {
		seq := layer.FrameSequenceOrDefault(seqIndex)
		n = max(n, len(seq)*max(1, layer.FrameRepeat))
	}
	return n
}
