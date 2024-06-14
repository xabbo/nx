package res

import (
	"fmt"
	"strings"

	"b7c.io/swfx"
	"golang.org/x/exp/maps"
)

type swfFigureLibraryLoader struct {
	swf *swfx.Swf
}

type swfFigurePartLibrary struct {
	name   string
	swf    *swfx.Swf
	assets Assets
}

// Returns a loader for the specified SWF figure part library.
func NewSwfFigureLibraryLoader(swf *swfx.Swf) LibraryLoader {
	return &swfFigureLibraryLoader{swf}
}

func (loader swfFigureLibraryLoader) Load() (lib AssetLibrary, err error) {
	swf := loader.swf

	var manifestTag swfx.CharacterTag
	for symbol, id := range swf.Symbols {
		if strings.HasSuffix(symbol, "_manifest") {
			manifestTag = swf.Characters[id]
			break
		}
	}

	if manifestTag == nil {
		err = fmt.Errorf("manifest not found")
		return
	}

	manifestData, ok := manifestTag.(*swfx.DefineBinaryData)
	if !ok {
		err = fmt.Errorf("invalid manifest tag type: %T", manifestTag)
		return
	}

	var manifest Manifest
	err = manifest.Unmarshal(manifestData.Data)
	if err != nil {
		return
	}

	lib = &swfFigurePartLibrary{
		name:   manifest.Name,
		swf:    swf,
		assets: manifest.Assets,
	}
	return
}

func (lib *swfFigurePartLibrary) Name() string {
	return lib.name
}

func (lib *swfFigurePartLibrary) Asset(name string) (asset *Asset, err error) {
	var ok bool

	asset, ok = lib.assets[name]
	if !ok {
		err = fmt.Errorf("asset %s/%q not found", lib.name, name)
		return
	}

	if asset.Image == nil {
		ch, ok := lib.swf.Symbols[lib.name+"_"+name]
		if !ok {
			err = fmt.Errorf("symbol %q not found", lib.name+"_"+name)
			return
		}

		tag, ok := lib.swf.Characters[ch]
		if !ok {
			err = fmt.Errorf("character %d not found", ch)
			return
		}

		imageTag, ok := tag.(swfx.ImageTag)
		if !ok {
			err = fmt.Errorf("asset is not an image")
			return
		}

		asset.Image, err = imageTag.Decode()
		if err != nil {
			return
		}
	}

	return
}

func (lib *swfFigurePartLibrary) Assets() []string {
	return maps.Keys(lib.assets)
}

func (lib *swfFigurePartLibrary) AssetExists(name string) bool {
	_, exists := lib.assets[name]
	return exists
}
