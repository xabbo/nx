package imager

import (
	"image"
	"image/png"
	"io"
)

type pngEncoder struct{}

func NewEncoderPNG() StaticEncoder {
	return pngEncoder{}
}

func (e pngEncoder) EncodeImage(w io.Writer, frame image.Image) error {
	return png.Encode(w, frame)
}

func (e pngEncoder) EncodeFrame(w io.Writer, anim Animation, seqIndex, frameIndex int) error {
	img := RenderFrame(anim, seqIndex, frameIndex)
	return png.Encode(w, img)
}
