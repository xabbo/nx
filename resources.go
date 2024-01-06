package nx

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"image"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/net/html/charset"

	"github.com/b7c/swfx"

	x "github.com/b7c/nx/xml"
)

type AssetManager interface {
	Library(name string) AssetLibrary
	Libraries() []string
	LibraryExists(name string) bool
	Load(AssetLibraryLoader) error
}

type AssetLibraryLoader interface {
	Load() (AssetLibrary, error)
}

type AssetLibrary interface {
	Name() string
	Asset(name string) (Asset, error)
	Assets() []string
	AssetExists(name string) bool
}

type swfFigureLibraryLoader struct {
	swf *swfx.Swf
}

type swfFigurePartLibrary struct {
	name   string
	swf    *swfx.Swf
	assets Assets
}

type Manifest struct {
	Libraries map[string]ManifestLibrary
}

type ManifestLibrary struct {
	Name    string
	Version string
	Assets  Assets
}

type Assets map[string]Asset

type Asset struct {
	Name   string      // The asset's name.
	Source *Asset      // The source asset.
	FlipH  bool        // Whether the asset is flipped horizontally.
	FlipV  bool        // Whether the asset is flipped vertically.
	Offset image.Point // The asset's image offset.
	Image  image.Image // The asset's image.
}

type assetManager struct {
	libs map[string]AssetLibrary
}

func (mgr *assetManager) Library(name string) AssetLibrary {
	return mgr.libs[name]
}

func (mgr *assetManager) Libraries() []string {
	return maps.Keys(mgr.libs)
}

func (mgr *assetManager) LibraryExists(name string) bool {
	_, exists := mgr.libs[name]
	return exists
}

func (mgr *assetManager) Load(loader AssetLibraryLoader) (err error) {
	library, err := loader.Load()
	if err == nil {
		if _, exists := mgr.libs[library.Name()]; exists {
			err = fmt.Errorf("library already loaded: %q", library.Name())
		} else {
			mgr.libs[library.Name()] = library
		}
	}
	return
}

func (loader swfFigureLibraryLoader) Load() (lib AssetLibrary, err error) {
	swf := loader.swf

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

	lib = &swfFigurePartLibrary{
		name:   libraryName,
		swf:    swf,
		assets: library.Assets,
	}
	return
}

func (lib *swfFigurePartLibrary) Name() string {
	return lib.name
}

func (lib *swfFigurePartLibrary) Asset(name string) (asset Asset, err error) {
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

func parsePoint(s string) (pt image.Point, err error) {
	split := strings.Split(s, ",")
	if len(split) != 2 {
		err = fmt.Errorf("invalid point")
		return
	}
	var x, y int
	x, err = strconv.Atoi(split[0])
	if err != nil {
		err = fmt.Errorf("invalid point")
		return
	}
	y, err = strconv.Atoi(split[1])
	if err != nil {
		err = fmt.Errorf("invalid point")
		return
	}
	pt = image.Pt(x, y)
	return
}

func (m *Manifest) Unmarshal(data []byte) (err error) {
	var xm x.Manifest
	reader := bytes.NewReader(data)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&xm)
	if err != nil {
		return
	}

	*m = Manifest{Libraries: make(map[string]ManifestLibrary)}
	for _, xl := range xm.Libraries {
		lib := ManifestLibrary{
			Name:    xl.Name,
			Version: xl.Version,
			Assets:  Assets{},
		}

		for _, xa := range xl.Assets {
			asset := Asset{
				Name: xa.Name,
			}
			for _, param := range xa.Params {
				if param.Key == "offset" {
					offset, err := parsePoint(param.Value)
					if err != nil {
						return fmt.Errorf("invalid offset: %q", param.Value)
					}
					asset.Offset = offset
				}
			}
			lib.Assets[asset.Name] = asset
		}

		m.Libraries[xl.Name] = lib
	}

	return
}

func (a *Assets) Unmarshal(data []byte) (err error) {
	// Used in furni data
	err = fmt.Errorf("not implemented")
	return
}

// Returns a loader for the specified SWF figure part library.
func SwfFigureLibraryLoader(swf *swfx.Swf) AssetLibraryLoader {
	return &swfFigureLibraryLoader{swf}
}
