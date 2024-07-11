package imager

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"slices"

	"golang.org/x/exp/maps"
)

type svgEncoder struct{}

func NewEncoderSVG() FrameEncoder {
	return svgEncoder{}
}

func (e svgEncoder) EncodeFrame(w io.Writer, anim Animation, seqIndex, frameIndex int) (err error) {
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
	<svg xmlns="http://www.w3.org/2000/svg" xmlns:svg="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" >
	<g inkscape:Label="Layer 1" inkscape:groupmode="layer" id="layer1">`))

	layerIds := maps.Keys(anim.Layers)
	slices.Sort(layerIds)
	for layerId := range layerIds {
		layer := anim.Layers[layerId]
		seq := layer.FrameSequenceOrDefault(seqIndex)
		frameId := seq[(frameIndex/max(1, layer.FrameRepeat))%len(seq)]
		frame := layer.Frames[frameId]
		for _, sprite := range frame {
			bounds := sprite.Bounds()
			img := image.NewRGBA(bounds)
			sprite.Draw(img, draw.Over)
			svgImgLayer{
				x:    bounds.Min.X,
				y:    bounds.Min.Y,
				w:    bounds.Dx(),
				h:    bounds.Dy(),
				name: sprite.Asset.Name,
				img:  img,
			}.Write(w)
		}
	}
	w.Write([]byte("</g></svg>"))

	return
}

type svgImgLayer struct {
	x, y, w, h int
	name       string
	img        image.Image
}

func (l svgImgLayer) Write(w io.Writer) {
	template := `<image x="%d" y="%d" width="%d" height="%d" ` +
		`inkscape:label="%s" preserveAspectRatio="none" style="image-rendering:optimizeSpeed" ` +
		`xlink:href="data:image/png;base64,%s"></image>`
	buf := &bytes.Buffer{}
	png.Encode(buf, l.img)
	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	w.Write([]byte(fmt.Sprintf(template, l.x, l.y, l.w, l.h, l.name, b64)))
}
