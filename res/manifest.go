package res

import (
	"fmt"

	x "xabbo.b7c.io/nx/raw/xml"
)

// A manifest describing the library and assets contained.
type Manifest struct {
	Name    string
	Version string
	Assets  Assets
}

func (m *Manifest) Unmarshal(data []byte) (err error) {
	var xm x.Manifest
	err = decodeXml(data, &xm)
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
		m.Assets[asset.Name] = &asset
	}

	return
}
