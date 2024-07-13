package imager

import (
	"image"
	"io"

	"xabbo.b7c.io/nx"
)

type FurniImager interface {
	Compose(furni Furni) Animation
}

type AvatarImager interface {
	Compose(avatar Avatar) (Animation, error)
	Parts(figure nx.Figure) ([]AvatarPart, error)
	RequiredLibs(figure nx.Figure) ([]string, error)
}

type Encoder interface {
	StaticEncoder
	AnimatedEncoder
}

type StaticEncoder interface {
	FrameEncoder
	ImageEncoder
}

type AnimatedEncoder interface {
	AnimationEncoder
	AnimatedImageEncoder
}

type AnimationEncoder interface {
	EncodeAnimation(w io.Writer, anim Animation, seqIndex, frameCount int) error
}

type AnimatedImageEncoder interface {
	EncodeImages(w io.Writer, frames []image.Image) error
}

type FrameEncoder interface {
	EncodeFrame(w io.Writer, anim Animation, seqIndex, frameIndex int) error
}

type ImageEncoder interface {
	EncodeImage(w io.Writer, frame image.Image) error
}
