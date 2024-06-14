package render

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	"github.com/disintegration/imaging"
)

func (frame Frame) ToImage() image.Image {
	canvas := image.NewRGBA(frame.Bounds())

	for _, layer := range frame {
		for _, sprite := range layer.Sprites {
			srcImg := sprite.Image()
			if srcImg == nil {
				fmt.Fprintf(os.Stderr, "no src img\n")
				continue
			}
			bounds := srcImg.Bounds() // sprite.Bounds()
			offset := sprite.Offset
			if sprite.FlipH {
				// bounds = bounds.Add(image.Point{0, 0})
				offset.X = offset.X*-1 + srcImg.Bounds().Dx() - 64
				srcImg = imaging.FlipH(srcImg)
			}
			bounds = bounds.Sub(offset)
			draw.Over.Draw(canvas, bounds, srcImg, image.Point{})
		}
	}

	return canvas
}
