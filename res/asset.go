package res

import (
	"image"

	"xabbo.b7c.io/nx/raw/nitro"
	x "xabbo.b7c.io/nx/raw/xml"
)

type Assets map[string]*Asset

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

func (a *Assets) UnmarshalBytes(b []byte) (err error) {
	var xAssets x.Assets
	err = decodeXml(b, &xAssets)
	if err != nil {
		return err
	}
	sourceMap := make(map[string]string)
	*a = make(map[string]*Asset, len(xAssets.Assets))
	for _, xAsset := range xAssets.Assets {
		var asset Asset
		asset.fromXml(xAsset)
		(*a)[asset.Name] = &asset
		if xAsset.Source != "" {
			sourceMap[asset.Name] = xAsset.Source
		}
	}
	for k, v := range sourceMap {
		(*a)[k].Source = (*a)[v]
	}
	return
}

func (a *Asset) fromXml(xAsset x.Asset) {
	a.Name = xAsset.Name
	a.FlipH = xAsset.FlipH
	a.FlipV = xAsset.FlipV
	a.Offset = image.Point{xAsset.X, xAsset.Y}
}

func (a *Asset) fromNitro(name string, src nitro.Asset) *Asset {
	*a = Asset{
		Name:   name,
		FlipH:  src.FlipH,
		FlipV:  src.FlipV,
		Offset: image.Point{X: src.X, Y: src.Y},
	}
	return a
}
