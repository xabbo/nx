package render

import (
	"fmt"
	"strings"

	"github.com/b7c/swfx"
)

type AssetLibrary interface {
	Name() string
	Asset(name string) (*Asset, error)
}

type AssetManager interface {
	Asset(libraryName, assetName string) (*Asset, error)
}

type SwfAssetManager struct {
	libs map[string]AssetLibrary
}

func NewSwfAssetManager() *SwfAssetManager {
	return &SwfAssetManager{
		libs: make(map[string]AssetLibrary),
	}
}

func (mgr *SwfAssetManager) LoadFurniLib(swf *swfx.Swf) (err error) {
	return fmt.Errorf("not implemented")
}

func (mgr *SwfAssetManager) LoadFigurePartLib(swf *swfx.Swf) (err error) {
	lib, err := NewSwfFigurePartLibrary(swf)
	if err == nil {
		mgr.libs[lib.name] = lib
	}
	return
}

func (mgr *SwfAssetManager) Asset(libraryName, assetName string) (asset *Asset, err error) {
	lib, ok := mgr.libs[libraryName]
	if !ok {
		err = fmt.Errorf("library not found: %q", libraryName)
		return
	}
	return lib.Asset(assetName)
}

type SwfFigurePartLibrary struct {
	name   string
	swf    *swfx.Swf
	assets Assets
}

func NewSwfFigurePartLibrary(swf *swfx.Swf) (lib *SwfFigurePartLibrary, err error) {
	var libraryName string
	var manifestTag swfx.CharacterTag
	for symbol, id := range swf.Symbols {
		var ok bool
		if libraryName, ok = strings.CutSuffix(symbol, "_manifest"); ok {
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

	library, ok := manifest.Libraries[libraryName]
	if !ok {
		err = fmt.Errorf("failed to find library")
		return
	}

	lib = &SwfFigurePartLibrary{
		name:   libraryName,
		swf:    swf,
		assets: library.Assets,
	}
	return
}

func (lib *SwfFigurePartLibrary) Name() string {
	return lib.name
}

func (lib *SwfFigurePartLibrary) AssetExists(name string) bool {
	_, exists := lib.assets[name]
	return exists
}

func (lib *SwfFigurePartLibrary) Asset(name string) (asset *Asset, err error) {
	var ok bool

	*asset, ok = lib.assets[name]
	if !ok {
		err = fmt.Errorf("asset not found")
		return
	}

	if asset.Image == nil {
		ch, ok := lib.swf.Symbols[lib.name+"_"+name]
		if !ok {
			err = fmt.Errorf("asset not found")
			return
		}

		tag, ok := lib.swf.Characters[ch]
		if !ok {
			err = fmt.Errorf("asset not found")
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
