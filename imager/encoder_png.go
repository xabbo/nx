package imager

import (
	"image"
	"image/png"
	"io"
)

type pngEncoder struct{}

func NewEncoderPNG() Encoder {
	return pngEncoder{}
}

func (e pngEncoder) EncodeImages(w io.Writer, imgs []image.Image) error {
	return png.Encode(w, imgs[0])
}

func (e pngEncoder) EncodeAnimation(w io.Writer, anim Animation, seqIndex, frameCount int) error {
	img := RenderFrame(anim, seqIndex, 0)
	return png.Encode(w, img)
}

func (e pngEncoder) EncodeFrame(w io.Writer, anim Animation, seqIndex, frameIndex int) error {
	img := RenderFrame(anim, seqIndex, frameIndex)
	return png.Encode(w, img)
}
