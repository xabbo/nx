package render

import (
	"image"
	"image/color"

	"xabbo.b7c.io/nx/res"
)

type Sprite struct {
	Asset  res.Asset
	Name   string
	FlipH  bool
	FlipV  bool
	Offset image.Point
	Color  color.Color
	Blend  int
	Order  int
}

func (s *Sprite) Size() image.Point {
	return s.Asset.Image.Bounds().Size()
}

type Layer struct {
	Names   string
	Sprites []Sprite
}

type Frame struct {
	Layers []Layer
}

type Animation struct {
	Frames []Frame
}
