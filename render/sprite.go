package render

import (
	"image"
	"image/color"

	"github.com/b7c/nx"
)

type Sprite struct {
	Asset  nx.Asset
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

type Frame struct {
	Sprites []Sprite
}

type Animation struct {
	Frames []Frame
}
