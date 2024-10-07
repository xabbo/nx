package res

import (
	x "xabbo.io/nx/raw/xml"
)

// An Index describes the furni visualization and logic types.
type Index struct {
	Type          string
	Visualization string
	Logic         string
}

func (index *Index) UnmarshalBytes(b []byte) (err error) {
	var xIndex x.Index
	err = decodeXml(b, &xIndex)
	if err != nil {
		return
	}
	*index = Index{
		Type:          xIndex.Type,
		Visualization: xIndex.Visualization,
		Logic:         xIndex.Logic,
	}
	return
}
