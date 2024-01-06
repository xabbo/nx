package render

import (
	"encoding/xml"
	"fmt"
	"image"
	"strconv"
	"strings"

	x "github.com/b7c/nx/xml"
)

type Manifest struct {
	Libraries map[string]Library
}

type Library struct {
	Name    string
	Version string
	Assets  Assets
}

type Assets map[string]Asset

type Asset struct {
	Name   string
	Source *Asset
	FlipH  bool
	FlipV  bool
	Offset image.Point
	Image  image.Image
}

func parsePoint(s string) (pt image.Point, err error) {
	split := strings.Split(s, ",")
	if len(split) != 2 {
		err = fmt.Errorf("invalid point")
	}
	var x, y int
	x, err = strconv.Atoi(split[0])
	if err != nil {
		err = fmt.Errorf("invalid point")
	}
	y, err = strconv.Atoi(split[1])
	if err != nil {
		err = fmt.Errorf("invalid point")
	}
	pt = image.Pt(x, y)
	return
}

func (m *Manifest) Unmarshal(data []byte) (err error) {
	var xm x.Manifest
	err = xml.Unmarshal(data, &xm)
	if err != nil {
		return
	}

	*m = Manifest{Libraries: make(map[string]Library)}
	for _, xl := range xm.Libraries {
		lib := Library{
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
	return
}
