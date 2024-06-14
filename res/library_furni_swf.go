package res

import (
	"fmt"
	"image"

	"b7c.io/swfx"
	"golang.org/x/exp/maps"
)

type swfFurniLibraryLoader struct {
	name string
	swf  *swfx.Swf
}

type swfFurniLibrary struct {
	swf            *swfx.Swf
	name           string
	index          *Index
	manifest       *Manifest
	logic          *Logic
	visualizations map[int]Visualization
	assets         map[string]*Asset
}

func NewFurniLibraryLoader(name string, swf *swfx.Swf) LibraryLoader {
	return &swfFurniLibraryLoader{name, swf}
}

func (loader *swfFurniLibraryLoader) Load() (assetLib AssetLibrary, err error) {
	swf := loader.swf

	indexTag := getBinaryTag(swf, loader.name+"_index")
	if indexTag == nil {
		err = fmt.Errorf("failed to find index in library %q", loader.name)
		return
	}
	var index Index
	err = index.UnmarshalBytes(indexTag.Data)
	if err != nil {
		return
	}

	manifestTag := getBinaryTag(swf, loader.name+"_manifest")
	if manifestTag == nil {
		err = fmt.Errorf("failed to find manifest in library %q", loader.name)
		return
	}
	var manifest Manifest
	err = manifest.Unmarshal(manifestTag.Data)
	if err != nil {
		return
	}

	logicTag := getBinaryTag(swf, loader.name+"_"+loader.name+"_logic")
	if logicTag == nil {
		err = fmt.Errorf("failed to find logic in library %q", loader.name)
		return
	}
	var logic Logic
	err = logic.UnmarshalBytes(logicTag.Data)
	if err != nil {
		return
	}

	visTag := getBinaryTag(swf, loader.name+"_"+loader.name+"_visualization")
	if visTag == nil {
		err = fmt.Errorf("failed to find visualization in library %q", loader.name)
		return
	}
	var visData VisualizationData
	err = visData.UnmarshalBytes(visTag.Data)
	if err != nil {
		return
	}

	assetsTag := getBinaryTag(swf, loader.name+"_"+loader.name+"_assets")
	if assetsTag == nil {
		err = fmt.Errorf("failed to find assets in library %q", loader.name)
		return
	}
	var assetsMap Assets
	err = assetsMap.UnmarshalBytes(assetsTag.Data)
	if err != nil {
		return
	}

	lib := &swfFurniLibrary{
		swf:            swf,
		name:           loader.name,
		index:          &index,
		manifest:       &manifest,
		logic:          &logic,
		visualizations: visData.Visualizations,
		assets:         assetsMap,
	}

	for assetName := range assetsMap {
		if assetsMap[assetName].Source != nil {
			continue
		}
		imgTag := getImageTag(swf, loader.name+"_"+assetName)
		if imgTag == nil {
			err = fmt.Errorf("failed to find asset %q in %q",
				assetName, loader.name)
			continue
		}
		var img image.Image
		img, err = imgTag.Decode()
		if err != nil {
			return
		}
		assetsMap[assetName].Image = img
	}

	assetLib = lib
	return
}

func (lib *swfFurniLibrary) Name() string {
	return lib.name
}

func (lib *swfFurniLibrary) Index() *Index {
	return lib.index
}

func (lib *swfFurniLibrary) Manifest() *Manifest {
	return lib.manifest
}

func (lib *swfFurniLibrary) Logic() *Logic {
	return lib.logic
}

func (lib *swfFurniLibrary) Visualizations() map[int]Visualization {
	return lib.visualizations
}

func (lib *swfFurniLibrary) Asset(name string) (asset *Asset, err error) {
	asset, ok := lib.assets[name]
	if !ok {
		err = fmt.Errorf("asset %q not found in library %q", name, lib.name)
	}
	return
}

func (lib *swfFurniLibrary) Assets() []string {
	return maps.Keys(lib.assets)
}

func (lib *swfFurniLibrary) AssetExists(name string) bool {
	_, exists := lib.assets[name]
	return exists
}
