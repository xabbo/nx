package gamedata

import (
	"encoding/xml"

	"xabbo.b7c.io/nx"
	x "xabbo.b7c.io/nx/xml"
)

// FigureData defines the figure part sets and color palettes used for Habbo avatars.
type FigureData struct {
	// Palettes maps figure color palettes by ID.
	// Each figure part set uses a certain color palette.
	Palettes map[int]FigureColorPalette
	// SetPalettes maps figure part types to color palette IDs.
	SetPalettes map[nx.FigurePartType]int
	Sets        map[nx.FigurePartType]FigurePartSetMap
}

// A FigureColorPalette maps FigurePartColorInfo by ID.
type FigureColorPalette map[int]FigurePartColorInfo

// A FigurePartSetMap maps FigurePartsetInfo by ID.
type FigurePartSetMap map[int]FigurePartSetInfo

// A FigurePartColorInfo defines information used to color figure parts.
type FigurePartColorInfo = x.FigureColor

// A FigurePartSetInfo contains information about a collection of figure parts.
type FigurePartSetInfo struct {
	Id            int
	Gender        string
	Club          int
	Colorable     bool // Whether this part set is colorable.
	Selectable    bool // Whether this part set can be selected.
	Preselectable bool
	Parts         []FigurePartInfo    // The parts contained in this part set.
	HiddenLayers  []nx.FigurePartType // Defines layers to be hidden when this part set is worn.
}

// A FigurePartInfo contains information about a figure part.
type FigurePartInfo struct {
	Id         int
	Type       nx.FigurePartType
	Colorable  bool
	Index      int
	ColorIndex int
}

// PaletteFor finds the color palette for the specified figure part type.
func (fd *FigureData) PaletteFor(partType nx.FigurePartType) FigureColorPalette {
	return fd.Palettes[fd.SetPalettes[partType]]
}

// Unmarshals an XML document as raw bytes into a FigureData.
func (fd *FigureData) UnmarshalBytes(data []byte) (err error) {
	var xfd *x.FigureData
	err = xml.Unmarshal(data, &xfd)
	if err != nil {
		return
	}

	*fd = FigureData{}
	fd.Palettes = map[int]FigureColorPalette{}
	fd.SetPalettes = map[nx.FigurePartType]int{}
	fd.Sets = map[nx.FigurePartType]FigurePartSetMap{}

	for _, p := range xfd.Palettes {
		palette := FigureColorPalette{}
		for _, c := range p.Colors {
			palette[c.Id] = c
		}
		fd.Palettes[p.Id] = palette
	}

	for _, xSetType := range xfd.Sets {
		partSetType := nx.FigurePartType(xSetType.Type)

		setMap := FigurePartSetMap{}
		for _, xSet := range xSetType.Sets {
			partSet := FigurePartSetInfo{}
			for _, xPart := range xSet.Parts {
				part := FigurePartInfo{
					Id:         xPart.Id,
					Type:       nx.FigurePartType(xPart.Type),
					Colorable:  xPart.Colorable,
					Index:      xPart.Index,
					ColorIndex: xPart.ColorIndex,
				}
				partSet.Parts = append(partSet.Parts, part)
			}
			for _, xLayer := range xSet.HiddenLayers {
				partSet.HiddenLayers = append(partSet.HiddenLayers,
					nx.FigurePartType(xLayer.PartType))
			}
			setMap[xSet.Id] = partSet
		}

		fd.Sets[partSetType] = setMap
		fd.SetPalettes[partSetType] = xSetType.PaletteId
	}

	return
}
