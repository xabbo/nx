package imager

type EncoderOption func(any)

func WithAlphaThreshold(threshold uint16) EncoderOption {
	return func(a any) {
		if opts, ok := a.(*GifEncoderOptions); ok {
			opts.AlphaThreshold = threshold
		}
	}
}

func WithColors(n int) EncoderOption {
	return func(a any) {
		if opts, ok := a.(*GifEncoderOptions); ok {
			opts.Colors = n
		}
	}
}
