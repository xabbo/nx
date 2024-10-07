package origins

import (
	"strings"

	"xabbo.io/nx"
	"xabbo.io/nx/gamedata"
)

// ColorMap is a mapping from
// Figure Part Type -> Color (lowercase hex) -> Modern Color ID
type ColorMap = map[nx.FigurePartType]map[string]int

// MakeColorMap creates a ColorMap from the specified figure data.
func MakeColorMap(fd *gamedata.FigureData) ColorMap {
	colorMap := map[nx.FigurePartType]map[string]int{}
	for partType, paletteId := range fd.SetPalettes {
		colorMap[partType] = map[string]int{}
		palette := fd.Palettes[paletteId]
		for _, color := range palette {
			colorMap[partType][strings.ToLower(color.Value)] = color.Id
		}
	}
	return colorMap
}
