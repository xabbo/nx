package imager

import (
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/color"
	"os"

	"xabbo.b7c.io/nx/res"
)

const shadowAlpha = 46

type Furni struct {
	Identifier string
	Size       int
	Direction  int
	State      int
	Sequence   int // Animation sequence to use.
	Color      int
	Shadow     bool // Whether to render the shadow.
}

type furniImager struct {
	mgr res.LibraryManager
}

func NewFurniImager(mgr res.LibraryManager) *furniImager {
	return &furniImager{mgr}
}

// Compose composes the furni to an Animation.
func (r *furniImager) Compose(furni Furni) (anim Animation, err error) {
	assetLib := r.mgr.Library(furni.Identifier)
	if assetLib == nil {
		err = errors.New("no library found")
		return
	}

	var lib res.FurniLibrary
	var ok bool
	if lib, ok = assetLib.(res.FurniLibrary); !ok {
		err = errors.New("not a furni library")
		return
	}

	index := lib.Index()
	if index == nil {
		err = errors.New("no index found")
		return
	}

	visuals := lib.Visualizations()
	vis, ok := visuals[furni.Size]
	if !ok {
		err = errors.New("invalid size")
		return
	}

	if _, ok := vis.Directions[furni.Direction]; !ok {
		err = fmt.Errorf("no visualization for direction %d [%s]", furni.Direction, furni.Identifier)
		return
	}

	vAnim, ok := vis.Animations[furni.State]
	if !ok {
		if furni.State == 0 {
			vAnim = &res.Animation{}
		} else {
			err = fmt.Errorf("no animation for state %d [%s]", furni.State, furni.Identifier)
			return
		}
	}

	anim.Layers = map[int]AnimationLayer{}

	for i := range vis.LayerCount + 1 {
		layerId := i - 1
		if layerId < 0 && !furni.Shadow {
			continue
		}

		layer := vAnim.Layers[layerId]
		frameRepeat := 0
		frameSequences := []res.FrameSequence{[]int{0}}

		requiredFrames := map[int]struct{}{}
		if layer != nil {
			frameRepeat = layer.FrameRepeat
			frameSequences = layer.FrameSequences
			for _, seq := range layer.FrameSequences {
				for _, id := range seq {
					requiredFrames[id] = struct{}{}
				}
			}
		}
		if len(requiredFrames) == 0 {
			requiredFrames[0] = struct{}{}
		}

		ink := ""
		alpha := uint8(255)
		z := 0
		if visLayer, ok := vis.Layers[layerId]; ok {
			ink = visLayer.Ink
			if visLayer.Alpha > 0 {
				alpha = uint8(visLayer.Alpha)
			}
			z = visLayer.Z
		}

		var blend Blend
		switch ink {
		case "ADD":
			blend = BlendAdd
		case "COPY":
			blend = BlendCopy
		default:
			blend = BlendNone
		}

		if layerId < 0 {
			blend = BlendCopy
			alpha = shadowAlpha
		}

		frames := map[int]Frame{}
		for frameId := range requiredFrames {
			spec := res.FurniAssetSpec{
				Name:      furni.Identifier,
				Size:      furni.Size,
				Layer:     layerId,
				Direction: furni.Direction,
				Frame:     frameId,
			}
			assetName := spec.String()
			if !lib.AssetExists(assetName) {
				continue
			}
			var asset *res.Asset
			asset, err = lib.Asset(spec.String())
			if err != nil {
				fmt.Fprintf(os.Stderr, "frame not found: %q", spec.String())
				return
			}

			col := color.Color(color.White)
			if colors, ok := vis.Colors[furni.Color]; ok {
				if colorLayer, ok := colors.Layers[layerId]; ok {
					if bytes, err := hex.DecodeString(colorLayer.Color); err == nil {
						col = color.RGBA{
							R: bytes[0],
							G: bytes[1],
							B: bytes[2],
							A: 255,
						}
					}
				}
			}

			offset := asset.Offset
			if asset.FlipH {
				offset = flipOffsetFurni(offset, asset.SourceImage().Bounds())
			}

			frames[frameId] = Frame{Sprite{
				Asset:  asset,
				FlipH:  asset.FlipH,
				FlipV:  asset.FlipV,
				Offset: offset,
				Blend:  blend,
				Color:  col,
				Alpha:  alpha,
			}}
		}

		anim.Layers[layerId] = AnimationLayer{
			Frames:      frames,
			FrameRepeat: frameRepeat,
			Sequences:   frameSequences,
			Z:           z,
		}
	}

	return
}

func flipOffsetFurni(offset image.Point, bounds image.Rectangle) image.Point {
	offset.X = -offset.X + bounds.Dx()
	return offset
}
