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
	for asset.Source != nil {
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
	Sequence []int
	Frames   map[int]Frame
}

func (animation Animation) Bounds() (bounds image.Rectangle) {
	for _, frame := range animation.Frames {
		bounds = bounds.Union(frame.Bounds())
	}
	return
}
