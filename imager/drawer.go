package imager

import (
	"image"
	"image/color"
	"image/draw"
)

type additiveDrawer struct{}

func (additiveDrawer) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	sx, sy := src.Bounds().Min.X+sp.X, src.Bounds().Min.Y+sp.Y
	sw, sh := src.Bounds().Dx(), src.Bounds().Dy()
	for y := 0; y < r.Dy() && (sp.Y+y) < sh; y++ {
		for x := 0; x < r.Dx() && (sp.X+x) < sw; x++ {
			sr, sg, sb, _ := src.At(sx+x, sy+y).RGBA()
			dr, dg, db, da := dst.At(r.Min.X+x, r.Min.Y+y).RGBA()

			if da > 0 {
				dr = (dr * 0xffff) / da
				dg = (dg * 0xffff) / da
				db = (db * 0xffff) / da
			}

			c := color.NRGBA64{
				R: uint16(min(0xffff, dr+sr)),
				G: uint16(min(0xffff, dg+sg)),
				B: uint16(min(0xffff, db+sb)),
				A: uint16(min(0xffff, da)),
			}
			dst.Set(r.Min.X+x, r.Min.Y+y, c)
		}
	}
}

type alphaDrawer uint8

func (drawer alphaDrawer) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	draw.DrawMask(dst, r, src, sp, image.NewUniform(color.Alpha{uint8(drawer)}), image.Point{}, draw.Over)
}
