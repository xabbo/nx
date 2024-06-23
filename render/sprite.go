package render

import (
	"image"
	"image/color"

	"xabbo.b7c.io/nx/res"
)

type Sprite struct {
	Asset  *res.Asset // A reference to the asset used by this sprite.
	Name   string
	FlipH  bool
	FlipV  bool
	Offset image.Point
	Color  color.Color
	Blend  int
	Order  int
}

func (s *Sprite) Image() image.Image {
	asset := s.Asset
	for asset.Source != nil && asset != asset.Source {
		asset = asset.Source
	}
	return asset.Image
}

// Returns the relative bounds of the sprite.
func (s *Sprite) Bounds() image.Rectangle {
	img := s.Image()
	if img == nil {
		return image.Rectangle{}
	}
	offset := s.Offset
	if s.FlipH {
		offset.X = -offset.X + img.Bounds().Dx() - 64
	}
	return img.Bounds().Sub(offset)
}

func (s *Sprite) Size() image.Point {
	return s.Bounds().Size()
}

type Layer struct {
	Id       int
	Name     string
	Blend    string
	Children []Layer
	Sprites  []Sprite
}

func (layer *Layer) Bounds() (bounds image.Rectangle) {
	for _, child := range layer.Children {
		bounds = bounds.Union(child.Bounds())
	}
	for _, sprite := range layer.Sprites {
		bounds = bounds.Union(sprite.Bounds())
	}
	return
}

type Frame []Layer

func (f Frame) Bounds() (bounds image.Rectangle) {
	for _, layer := range f {
		bounds = bounds.Union(layer.Bounds())
	}
	return
}

type Animation struct {
	Layers map[int]AnimationLayer
}

type AnimationLayer struct {
	Frames      map[int]Frame
	FrameRepeat int
	Sequences   []res.FrameSequence
}

func (animationLayer AnimationLayer) FrameSequenceOrDefault(i int) res.FrameSequence {
	if i < len(animationLayer.Sequences) {
		return animationLayer.Sequences[i]
	} else if len(animationLayer.Sequences) > 0 {
		return animationLayer.Sequences[0]
	} else {
		return []int{0}
	}
}

func (animation Animation) Bounds() (bounds image.Rectangle) {
	for _, layer := range animation.Layers {
		for _, frame := range layer.Frames {
			bounds = bounds.Union(frame.Bounds())
		}
	}
	return
}

func (animation *Animation) TotalFrames() int {
	n := 1
	for _, layer := range animation.Layers {
		sequenceLen := 1
		if len(layer.Sequences) > 0 {
			sequenceLen = len(layer.Sequences[0])
		}
		n = lcm(n, sequenceLen*max(1, layer.FrameRepeat))
	}
	return n
}

func (animation *Animation) LongestFrameSequence(seqIndex int) int {
	n := 1
	for _, layer := range animation.Layers {
		seq := layer.FrameSequenceOrDefault(seqIndex)
		n = max(n, len(seq)*max(1, layer.FrameRepeat))
	}
	return n
}
