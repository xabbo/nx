package res

import (
	"fmt"
	"image"
	"strings"

	"b7c.io/swfx"
	"golang.org/x/exp/maps"
)

type swfFurniLibrary struct {
	swf            *swfx.Swf
	name           string
	index          *Index
	manifest       *Manifest
	logic          *Logic
	visualizations map[int]*Visualization
	assets         map[string]*Asset
}

func LoadFurniLibrarySwf(swf *swfx.Swf) (assetLib AssetLibrary, err error) {
	// find manifest tag to extract library name
	var manifestTag *swfx.DefineBinaryData
	for symbolName, chId := range swf.Symbols {
		if strings.HasSuffix(symbolName, "_manifest") {
			manifestTag = swf.Characters[chId].(*swfx.DefineBinaryData)
			break
		}
	}
	if manifestTag == nil {
		err = fmt.Errorf("failed to find manifest in library")
		return
	}
	var manifest Manifest
	err = manifest.Unmarshal(manifestTag.Data)
	if err != nil {
		return
	}
	libName := manifest.Name

	indexTag := getBinaryTag(swf, libName+"_index")
	if indexTag == nil {
		err = fmt.Errorf("failed to find index in library %q", libName)
		return
	}
	var index Index
	err = index.UnmarshalBytes(indexTag.Data)
	if err != nil {
		return
	}

	logicTag := getBinaryTag(swf, libName+"_"+libName+"_logic")
	if logicTag == nil {
		err = fmt.Errorf("failed to find logic in library %q", libName)
		return
	}
	var logic Logic
	err = logic.UnmarshalBytes(logicTag.Data)
	if err != nil {
		return
	}

	visTag := getBinaryTag(swf, libName+"_"+libName+"_visualization")
	if visTag == nil {
		err = fmt.Errorf("failed to find visualization in library %q", libName)
		return
	}
	var visData VisualizationData
	err = visData.UnmarshalBytes(visTag.Data)
	if err != nil {
		return
	}

	assetsTag := getBinaryTag(swf, libName+"_"+libName+"_assets")
	if assetsTag == nil {
		err = fmt.Errorf("failed to find assets in library %q", libName)
		return
	}
	var assetsMap Assets
	err = assetsMap.UnmarshalBytes(assetsTag.Data)
	if err != nil {
		return
	}

	lib := &swfFurniLibrary{
		swf:            swf,
		name:           libName,
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
		imgTag := getImageTag(swf, libName+"_"+assetName)
		if imgTag == nil {
			err = fmt.Errorf("failed to find asset %q in %q",
				assetName, libName)
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

func (lib *swfFurniLibrary) Visualizations() map[int]*Visualization {
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
