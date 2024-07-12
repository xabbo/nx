package res

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"strings"

	"golang.org/x/exp/maps"

	"xabbo.b7c.io/nx/raw/nitro"
)

type nitroFurniLibrary struct {
	name           string
	index          *Index
	manifest       *Manifest
	logic          *Logic
	visualizations Visualizations
	assets         map[string]*Asset
}

func LoadFurniLibraryNitro(archive nitro.Archive) (furniLibrary FurniLibrary, err error) {
	nitroLib := &nitroFurniLibrary{
		assets: map[string]*Asset{},
	}

	// find metadata

	var metadataFile nitro.File
	for name := range archive.Files {
		if strings.HasSuffix(name, ".json") {
			metadataFile = archive.Files[name]
		}
	}

	if metadataFile.Data == nil {
		err = fmt.Errorf("failed to find metadata in Nitro archive")
		return
	}

	var nitroFurni nitro.Furni
	err = json.Unmarshal(metadataFile.Data, &nitroFurni)
	if err != nil {
		return
	}

	nitroLib.name = nitroFurni.Name
	nitroLib.index = &Index{
		Type:          nitroFurni.Name,
		Logic:         nitroFurni.LogicType,
		Visualization: nitroFurni.VisualizationType,
	}

	nitroLib.manifest = &Manifest{
		Name:   nitroFurni.Name,
		Assets: Assets{},
	}

	// populate manifest
	for assetName := range nitroFurni.Assets {
		nitroLib.manifest.Assets[assetName] = &Asset{
			Name: assetName,
		}
	}

	nitroLib.logic = new(Logic).fromNitro(&nitroFurni.Logic)
	nitroLib.visualizations = Visualizations{}.fromNitro(nitroFurni.Visualizations)

	sourceMap := map[string]string{}
	for name, asset := range nitroFurni.Assets {
		nitroLib.assets[name] = new(Asset).fromNitro(name, asset)
		if asset.Source != "" {
			sourceMap[name] = asset.Source
		}
	}
	for dstName, srcName := range sourceMap {
		nitroLib.assets[dstName].Source = nitroLib.assets[srcName]
	}

	// extract images from spritesheet
	bytesSpritesheet := archive.Files[nitroFurni.Spritesheet.Meta.Image].Data
	imgSprites, err := png.Decode(bytes.NewReader(bytesSpritesheet))
	if err != nil {
		return
	}

	for name, asset := range nitroLib.assets {
		name = nitroLib.name + "_" + name
		spriteInfo, ok := nitroFurni.Spritesheet.Frames[name]
		if !ok {
			spriteInfo, ok = nitroFurni.Spritesheet.Frames[name+".png"]
			if !ok {
				continue
			}
		}
		frame := spriteInfo.Frame
		size := image.Rect(0, 0, frame.W, frame.H)
		spriteImg := image.NewRGBA(image.Rect(0, 0, frame.W, frame.H))
		draw.Src.Draw(spriteImg, size, imgSprites, image.Point{frame.X, frame.Y})
		asset.Image = spriteImg
	}

	furniLibrary = nitroLib
	return
}

func (lib *nitroFurniLibrary) Name() string {
	return lib.name
}

func (lib *nitroFurniLibrary) Index() *Index {
	return lib.index
}

func (lib *nitroFurniLibrary) Manifest() *Manifest {
	return lib.manifest
}

func (lib *nitroFurniLibrary) Logic() *Logic {
	return lib.logic
}

func (lib *nitroFurniLibrary) Visualizations() map[int]*Visualization {
	return lib.visualizations
}

func (lib *nitroFurniLibrary) Asset(name string) (asset *Asset, err error) {
	asset, ok := lib.assets[name]
	if !ok {
		err = fmt.Errorf("asset %q not found in library %q", name, lib.name)
	}
	return
}

func (lib *nitroFurniLibrary) Assets() []string {
	return maps.Keys(lib.assets)
}

func (lib *nitroFurniLibrary) AssetExists(name string) bool {
	_, exists := lib.assets[name]
	return exists
}
