package res

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"image"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/net/html/charset"

	x "github.com/b7c/nx/xml"
)

// An index describing the furni visualization and logic types.
type Index struct {
	Type          string
	Visualization string
	Logic         string
}

// A manifest describing the library and assets contained.
type Manifest struct {
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

func (asset *Asset) SourceImage() image.Image {
	for asset.Source != nil {
		asset = asset.Source
	}
	return asset.Image
}

type assetManager struct {
	libs map[string]AssetLibrary
}

func NewManager() LibraryManager {
	return &assetManager{
		libs: map[string]AssetLibrary{},
	}
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

func (mgr *assetManager) Load(loader LibraryLoader) (err error) {
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

	*m = Manifest{
		Name:    xm.Library.Name,
		Version: xm.Library.Version,
		Assets:  Assets{},
	}

	for _, xa := range xm.Library.Assets {
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
		m.Assets[asset.Name] = asset
	}

	return
}

func (a *Assets) Unmarshal(data []byte) (err error) {
	// Used in furni data
	err = fmt.Errorf("not implemented")
	return
}
