package imager

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/disintegration/imaging"
	"github.com/phrozen/blend"
	"xabbo.io/nx/res"
)

// Sprite defines parameters for an asset to be drawn on a canvas.
type Sprite struct {
	Asset  *res.Asset  // A reference to the asset used by this sprite.
	FlipH  bool        // FlipH defines whether to flip the asset horizontally.
	FlipV  bool        // FlipV defines whether to flip the asset vertically.
	Offset image.Point // Offset defines the offset point from the origin to draw the asset.
	Color  color.Color // Color defines the color to blend with the asset.
	Blend  Blend       // Blend defines the blending mode used to draw the asset.
	Alpha  uint8       // Alpha defines the alpha transparency of the asset.
}

// Blend represents a color blending mode.
type Blend int

const (
	BlendNone Blend = iota // No blend mode.
	BlendAdd               // Additive blending.
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

// Bounds returns the bounds of the sprite's image translated by the sprite's offset.
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
func (s *Sprite) Draw(canvas draw.Image, offset image.Point, drawer draw.Drawer) {
	srcImg := s.Image()
	if srcImg == nil {
		return
	}
	if s.Color != nil && s.Color != color.White {
		srcImg = blend.BlendNewImage(srcImg, image.NewUniform(s.Color), blend.Multiply)
	}
	bounds := srcImg.Bounds()
	offset = offset.Add(s.Offset)
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

// Bounds returns the union of the bounds of all sprites in this frame.
func (f Frame) Bounds() (bounds image.Rectangle) {
	for _, sprite := range f {
		bounds = bounds.Union(sprite.Bounds())
	}
	return
}

// Animation represents an animated asset.
type Animation struct {
	Background color.Color            // Background defines the color to fill the canvas with when rendering.
	Layers     map[int]AnimationLayer // Layers is a map of animation layers by index.
}

// AnimationLayer defines a set of frames and frame sequences.
type AnimationLayer struct {
	Frames      map[int]Frame       // Frames is a map of frames by index.
	FrameRepeat int                 // FrameRepeat defines the duration of each frame for this layer.
	Sequences   []res.FrameSequence // Sequences contains a list of frame sequences.
	Z           int                 // Z defines the Z-order of this layer.
}

// SequenceOrDefault gets the specified frame sequence if it exists, the first sequence if `i` is out of range,
// or a sequence with a single zero frame if the animation has no frame sequences.
func (animationLayer AnimationLayer) SequenceOrDefault(seqIndex int) res.FrameSequence {
	if seqIndex >= 0 && seqIndex < len(animationLayer.Sequences) {
		return animationLayer.Sequences[seqIndex]
	} else if len(animationLayer.Sequences) > 0 {
		return animationLayer.Sequences[0]
	} else {
		return []int{0}
	}
}

// Bounds gets the union of the bounds of all frames in this animation for the specified sequence.
func (animation Animation) Bounds(seqIndex int) (bounds image.Rectangle) {
	for _, layer := range animation.Layers {
		for _, frame := range layer.Frames {
			bounds = bounds.Union(frame.Bounds())
		}
	}
	return
}

// TotalFrames gets the total number of frames in this animation for the specified sequence.
func (animation *Animation) TotalFrames(seqIndex int) int {
	n := 1
	for _, layer := range animation.Layers {
		seq := layer.SequenceOrDefault(seqIndex)
		n = lcm(n, len(seq)*max(1, layer.FrameRepeat))
	}
	return n
}

// LongestSequence gets the longest frame sequence of all layers in this animation for the specified sequence.
func (animation *Animation) LongestSequence(seqIndex int) int {
	n := 1
	for _, layer := range animation.Layers {
		seq := layer.SequenceOrDefault(seqIndex)
		n = max(n, len(seq)*max(1, layer.FrameRepeat))
	}
	return n
}
