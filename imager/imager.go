package imager

import (
	"image"
	"io"

	"xabbo.io/nx"
)

// FurniImager represents an imager that can compose furni into animations.
type FurniImager interface {
	Compose(furni Furni) Animation
}

// AvatarImager represents an imager that can compose avatars into animations.
type AvatarImager interface {
	Compose(avatar Avatar) (Animation, error)
	Parts(figure nx.Figure) ([]AvatarPart, error)
	RequiredLibs(figure nx.Figure) ([]string, error)
}

// Encoder represents an encoder that can encode animations and frames.
type Encoder interface {
	StaticEncoder
	AnimatedEncoder
}

// StaticEncoder represents an encoder that can encode an image or a single frame from an animation.
type StaticEncoder interface {
	FrameEncoder
	ImageEncoder
}

// AnimatedEncoder represents an encoder that can encode animations and image sequences.
type AnimatedEncoder interface {
	AnimationEncoder
	AnimatedImageEncoder
}

// AnimationSequenceEncoder represents an encoder that can encode a sequence of animations.
type AnimationSequenceEncoder interface {
	EncodeAnimations(w io.Writer, anims []Animation, seqIndex, frameCount int) error
}

// AnimationEncoder represents an encoder that can encode an animation.
type AnimationEncoder interface {
	EncodeAnimation(w io.Writer, anim Animation, seqIndex, frameCount int) error
}

// AnimatedImageEncoder represents an encoder that can encode a sequence of images.
type AnimatedImageEncoder interface {
	EncodeImages(w io.Writer, frames []image.Image) error
}

// FrameEncoder represents an encoder that can encode a single frame within an animation.
type FrameEncoder interface {
	EncodeFrame(w io.Writer, anim Animation, seqIndex, frameIndex int) error
}

// ImageEncoder represents an encoder that can encode an image.
type ImageEncoder interface {
	EncodeImage(w io.Writer, frame image.Image) error
}
