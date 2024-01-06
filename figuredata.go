package nx

import (
	"encoding/xml"
	x "github.com/b7c/nx/xml"
)

type FigureData struct {
	Palettes    map[int]FigureColorPalette
	SetPalettes map[FigurePartType]int
	Sets        map[FigurePartType]FigurePartSetMap
}

type FigureColorPalette map[int]FigurePartColorInfo

type FigurePartSetMap map[int]FigurePartSetInfo

type FigurePartColorInfo = x.FigureColor

type FigurePartSetInfo struct {
	Id            int
	Gender        string
	Club          int
	Colorable     bool
	Selectable    bool
	Preselectable bool
	Parts         []FigurePartInfo
	HiddenLayers  []FigurePartType
}

type FigurePartInfo struct {
	Id         int
	Type       FigurePartType
	Colorable  bool
	Index      int
	ColorIndex int
}

func (fd *FigureData) PaletteFor(partType FigurePartType) FigureColorPalette {
	return fd.Palettes[fd.SetPalettes[partType]]
}

func (fd *FigureData) UnmarshalBytes(data []byte) (err error) {
	var xfd *x.FigureData
	err = xml.Unmarshal(data, &xfd)
	if err != nil {
		return
	}

	*fd = FigureData{}
	fd.Palettes = map[int]FigureColorPalette{}
	fd.SetPalettes = map[FigurePartType]int{}
	fd.Sets = map[FigurePartType]FigurePartSetMap{}

	for _, p := range xfd.Palettes {
		palette := FigureColorPalette{}
		for _, c := range p.Colors {
			palette[c.Id] = c
		}
		fd.Palettes[p.Id] = palette
	}

	for _, xSetType := range xfd.Sets {
		partSetType := FigurePartType(xSetType.Type)

		setMap := FigurePartSetMap{}
		for _, xSet := range xSetType.Sets {
			partSet := FigurePartSetInfo{}
			for _, xPart := range xSet.Parts {
				part := FigurePartInfo{
					Id:         xPart.Id,
					Type:       FigurePartType(xPart.Type),
					Colorable:  xPart.Colorable,
					Index:      xPart.Index,
					ColorIndex: xPart.ColorIndex,
				}
				partSet.Parts = append(partSet.Parts, part)
			}
			for _, xLayer := range xSet.HiddenLayers {
				partSet.HiddenLayers = append(partSet.HiddenLayers,
					FigurePartType(xLayer.PartType))
			}
			setMap[xSet.Id] = partSet
		}

		fd.Sets[partSetType] = setMap
		fd.SetPalettes[partSetType] = xSetType.PaletteId
	}

	return
}
