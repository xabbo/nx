package render

import (
	"errors"
	"fmt"
	"slices"

	"golang.org/x/exp/maps"
	"xabbo.b7c.io/nx/gamedata"
	"xabbo.b7c.io/nx/res"
)

/*

API goals

Furni Rendering Pipeline
	Furni{Identifier, Direction, State}

gdm := gd.NewGamedataManager()
rn := render.NewFurniRenderer()

layerGrp := rn.Render(Furni{
	Identifier: "duck",
	Direction: 2,
	State: 0,
})

img := layerGrp.ToPng()
img.Save("duck.png")

*/

type FurniRenderer interface {
	Render(furni Furni) Animation
}

type furniRenderer struct {
	mgr *gamedata.GamedataManager
}

type Furni struct {
	Identifier string
	Size       int
	Direction  int
	State      int
}

func NewFurniRenderer(mgr *gamedata.GamedataManager) *furniRenderer {
	return &furniRenderer{mgr}
}

func (r *furniRenderer) Render(furni Furni) (anim Animation, err error) {
	assetLib := r.mgr.Assets.Library(furni.Identifier)
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

	// r.renderFrame(lib, )
	fr, err := r.renderFrame(lib, furni, 0)
	if err != nil {
		return
	}

	anim.Frames = map[int]Frame{0: fr}
	return
}

// Renders a single furni frame.
func (r *furniRenderer) renderFrame(lib res.FurniLibrary, furni Furni, frameNo int) (frame Frame, err error) {

	// get the visualization for the specified size
	vis, ok := lib.Visualizations()[furni.Size]
	if !ok {
		err = fmt.Errorf("no visualization for size %d [%s]", furni.Size, furni.Identifier)
		return
	}

	var layers []int
	if len(vis.Layers) > 0 {
		layers = maps.Keys(vis.Layers)
	} else {
		layers = []int{0}
	}

	for _, layerId := range layers {
		spec := res.FurniAssetSpec{
			Name:      furni.Identifier,
			Size:      furni.Size,
			Layer:     layerId,
			Direction: furni.Direction,
			Frame:     frameNo,
		}
		assetName := spec.String()
		if !lib.AssetExists(assetName) {
			continue
		}
		var asset *res.Asset
		asset, err = lib.Asset(spec.String())
		if err != nil {
			return
		}

		layer := Layer{
			Id: layerId,
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

	slices.SortFunc(frame, func(a, b Layer) int {
		return a.Id - b.Id
	})
	return
}

