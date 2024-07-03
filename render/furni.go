package render

import (
	"errors"
	"fmt"
	"os"

	"xabbo.b7c.io/nx/res"
)

type FurniRenderer interface {
	// Render assembles the furni into an Animation.
	Render(furni Furni) Animation
}

type furniRenderer struct {
	mgr res.LibraryManager
}

type Furni struct {
	Identifier string
	Size       int
	Direction  int
	State      int
	Seq        int // Animation sequence to use.
}

func NewFurniRenderer(mgr res.LibraryManager) *furniRenderer {
	return &furniRenderer{mgr}
}

// Renders the furni to an Animation.
// A furni with static visualization will return an animation with a single frame and no animation layers.
func (r *furniRenderer) Render(furni Furni) (anim Animation, err error) {
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

	switch index.Visualization {
	case "furniture_static":
		var fr Frame
		fr, err = r.renderStatic(lib, furni)
		if err != nil {
			return
		}
		anim.Layers = map[int]AnimationLayer{}
		anim.Layers[0] = AnimationLayer{
			Frames:    map[int]Frame{0: fr},
			Sequences: []res.FrameSequence{[]int{0}},
		}
	case "furniture_animated":
		anim, err = r.renderAnimated(lib, furni)
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("visualization type not implemented: %s", index.Visualization)
		return
	}

	return
}

// Renders a static furniture.
func (r *furniRenderer) renderStatic(lib res.FurniLibrary, furni Furni) (frame Frame, err error) {
	// get the visualization for the specified size
	vis, ok := lib.Visualizations()[furni.Size]
	if !ok {
		err = fmt.Errorf("no visualization for size %d [%s]", furni.Size, furni.Identifier)
		return
	}

	for i := range vis.LayerCount {
		spec := res.FurniAssetSpec{
			Name:      furni.Identifier,
			Size:      furni.Size,
			Layer:     i,
			Direction: furni.Direction,
			Frame:     0,
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

		layer := Layer{
			Id: i,
			Sprites: []Sprite{
				{
					Asset:  asset,
					FlipH:  asset.FlipH,
					FlipV:  asset.FlipV,
					Offset: asset.Offset,
				},
			},
		}
		frame = append(frame, layer)
	}

	return
}

// Renders an animated furniture.
func (r *furniRenderer) renderAnimated(lib res.FurniLibrary, furni Furni) (anim Animation, err error) {
	// get the visualization for the specified size
	vis, ok := lib.Visualizations()[furni.Size]
	if !ok {
		err = fmt.Errorf("no visualization for size %d [%s]", furni.Size, furni.Identifier)
		return
	}

	vAnim, ok := vis.Animations[furni.State]
	if !ok {
		err = fmt.Errorf("no animation for state %d [%s]", furni.State, furni.Identifier)
		return
	}

	anim.Layers = map[int]AnimationLayer{}

	for _, layer := range vAnim.Layers {
		requiredFrames := map[int]struct{}{}
		for _, seq := range layer.FrameSequences {
			for _, id := range seq {
				requiredFrames[id] = struct{}{}
			}
		}
		if len(requiredFrames) == 0 {
			requiredFrames[0] = struct{}{}
		}

		frames := map[int]Frame{}
		for frameId := range requiredFrames {
			spec := res.FurniAssetSpec{
				Name:      furni.Identifier,
				Size:      furni.Size,
				Layer:     layer.Id,
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

			ink := ""
			if visLayer, ok := vis.Layers[layer.Id]; ok {
				ink = visLayer.Ink
			}

			frames[frameId] = Frame{
				Layer{
					Id:    frameId,
					Blend: ink,
					Sprites: []Sprite{
						{
							Asset:  asset,
							FlipH:  asset.FlipH,
							FlipV:  asset.FlipV,
							Offset: asset.Offset,
						},
					},
				},
			}
		}

		anim.Layers[layer.Id] = AnimationLayer{
			Frames:      frames,
			FrameRepeat: layer.FrameRepeat,
			Sequences:   layer.FrameSequences,
		}
	}

	return
}
