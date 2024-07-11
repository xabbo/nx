package imager

import (
	"image"
	"io"
)

type FurniImager interface {
	Compose(furni Furni) Animation
}

type AvatarImager interface {
	Compose(avatar Avatar) Animation
}

type Renderer interface {
	RenderAnimation(w io.Writer, anim Animation)
	RenderFrame(w io.Writer, anim Animation, frameIndex int, frameSequence int)
}

type ImageRenderer interface {
	RenderAnimation(anim Animation) []image.Image
	RenderFrame(anim Animation, frameIndex int, frameSequence int) image.Image
}

type Encoder interface {
	ImageEncoder
	AnimationEncoder
}

type ImageEncoder interface {
	EncodeImages(w io.Writer, frames []image.Image) error
}

type AnimationEncoder interface {
	EncodeAnimation(w io.Writer, anim Animation, seqIdx, frameCount int) error
	FrameEncoder
}

type FrameEncoder interface {
	EncodeFrame(w io.Writer, anim Animation, seqIdx, frameIdx int) error
}
