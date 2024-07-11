package imager

import (
	"image"
	"image/color"
	"image/draw"
)

type additiveDrawer struct{}

func (additiveDrawer) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	sBounds := src.Bounds()
	sw, sh := src.Bounds().Dx(), src.Bounds().Dy()
	for y := 0; y < r.Dy() && (sp.Y+y) < sh; y++ {
		for x := 0; x < r.Dx() && (sp.X+x) < sw; x++ {
			dc := dst.At(r.Min.X+x, r.Min.Y+y)
			dr, dg, db, da := dc.RGBA()
			if da == 0 {
				continue
			}
			sc := src.At(sBounds.Min.X+sp.X+x, sBounds.Min.Y+sp.Y+y)
			sr, sg, sb, _ := sc.RGBA()
			dc = color.RGBA{
				R: uint8(float64(min(0xffff, dr+sr)) / 65535.0 * 255.0),
				G: uint8(float64(min(0xffff, dg+sg)) / 65535.0 * 255.0),
				B: uint8(float64(min(0xffff, db+sb)) / 65535.0 * 255.0),
				A: uint8(float64(da) / 65535.0 * 255.0),
			}
			dst.Set(r.Min.X+x, r.Min.Y+y, dc)
		}
	}
}
