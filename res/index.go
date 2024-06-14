package res

import (
	x "xabbo.b7c.io/nx/xml"
)

// An index describing the furni visualization and logic types.
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
		Type: xIndex.Type,
		Visualization: xIndex.Visualization,
		Logic: xIndex.Logic,
	}
	return
}
