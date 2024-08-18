package origins

import (
	"errors"
	"strconv"
	"strings"

	"xabbo.b7c.io/nx"
)

type FigureConverter struct {
	figureData *FigureData
	colorMap   ColorMap
	setIds     map[int]FigurePartSet
}

func NewFigureConverter(figureData *FigureData, colorMap ColorMap) *FigureConverter {
	fc := &FigureConverter{
		figureData: figureData,
		colorMap:   colorMap,
		setIds:     map[int]FigurePartSet{},
	}
	for _, genderSet := range []map[nx.FigurePartType]FigurePartSets{figureData.M, figureData.F} {
		for setType, items := range genderSet {
			for _, partSet := range items {
				partSet.Type = setType
				fc.setIds[partSet.Id] = partSet
			}
		}
	}
	return fc
}

// Hardcoded hair -> hat map
var hairToHatMap = map[int]int{
	// m
	120: 1001,
	130: 1010,
	140: 1004,
	150: 1003,
	160: 1004,
	175: 1006,
	176: 1007,
	177: 1008,
	178: 1009,
	800: 1012,
	801: 1011,
	802: 1013,
	// f
	525: 1002,
	535: 1003,
	565: 1004,
	570: 1005,
	580: 1007,
	585: 1006,
	590: 1008,
	595: 1009,
	810: 1012,
	811: 1013,
}

// Convert converts an origins figure string to its modern `nx.Figure` representation.
func (fc *FigureConverter) Convert(originsFigure string) (figure nx.Figure, err error) {
	if len(originsFigure) != 25 {
		err = errors.New("invalid figure string: must be 25 characters in length")
		return
	}

	for _, c := range originsFigure {
		if c < '0' || c > '9' {
			err = errors.New("invalid figure string: must consist only of numbers")
			return
		}
	}

	for i := 0; i < 25; i += 5 {
		setId, _ := strconv.Atoi(originsFigure[i : i+3])
		colorIndex, _ := strconv.Atoi(originsFigure[i+3 : i+5])

		set := fc.setIds[setId]
		figureItem := nx.FigureItem{
			Type: nx.FigurePartType(set.Type),
			Id:   setId,
		}

		partColor := strings.ToLower(set.Colors[colorIndex-1])
		colorId := fc.colorMap[set.Type][partColor]
		figureItem.Colors = append(figureItem.Colors, colorId)

		figure.Items = append(figure.Items, figureItem)

		if figureItem.Type == nx.Hair {
			if hatId, ok := hairToHatMap[figureItem.Id]; ok {
				figure.Items = append(figure.Items, nx.FigureItem{
					Type:   nx.Hat,
					Id:     hatId,
					Colors: []int{colorId},
				})
			}
		}
	}

	return figure, nil
}
