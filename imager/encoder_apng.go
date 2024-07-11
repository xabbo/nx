package imager

import (
	"image"
	"io"

	"github.com/kettek/apng"
)

type apngEncoder struct{}

func NewEncoderAPNG() Encoder {
	return apngEncoder{}
}

func (e apngEncoder) EncodeImages(w io.Writer, imgs []image.Image) error {
	a := apng.APNG{}
	for _, img := range imgs {
		a.Frames = append(a.Frames, apng.Frame{
			Image:            img,
			DelayNumerator:   1,
			DelayDenominator: 24,
		})
	}
	return apng.Encode(w, a)
}

func (e apngEncoder) EncodeAnimation(w io.Writer, anim Animation, seqIndex, frameCount int) error {
	imgs := RenderFrames(anim, seqIndex, frameCount)
	return e.EncodeImages(w, imgs)
}

func (e apngEncoder) EncodeFrame(w io.Writer, anim Animation, seqIndex, frameIndex int) error {
	img := RenderFrame(anim, seqIndex, frameIndex)
	return e.EncodeImages(w, []image.Image{img})
}
