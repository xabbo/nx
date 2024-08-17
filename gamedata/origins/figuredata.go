package origins

import (
	"encoding/json"
	"fmt"

	"xabbo.b7c.io/nx"
)

type FigureData struct {
	M, F map[nx.FigurePartType]FigurePartSets
}

type FigurePartSets []FigurePartSet

type FigurePartSet struct {
	Type   nx.FigurePartType `json:"-"`
	Id     int               `json:"s"`
	Parts  map[string]string `json:"p"`
	Colors []string          `json:"c"`
}

// ParseFigureData parses the origins FigureData from the specified byte array.
func ParseFigureData(b []byte) (fd *FigureData, err error) {
	err = fixFigureData(b)
	if err == nil {
		fd = &FigureData{}
		err = json.Unmarshal(b, fd)
	}
	if err != nil {
		err = fmt.Errorf("weirdness in figure data!!! %w", err)
	}
	return
}

// Fixes the origins figure data to valid JSON.
func fixFigureData(b []byte) (err error) {
	sp := -1
	stack := [16]struct {
		i      int
		object bool
	}{}

	for i := range b {
		// assuming these characters don't appear inside any strings
		switch b[i] {
		case '[':
			sp++
			if sp >= len(stack) {
				return fmt.Errorf("overflow in fixFigureData")
			}
			stack[sp].i = i
			stack[sp].object = false
		case ':':
			stack[sp].object = true
		case ']':
			if stack[sp].object {
				b[stack[sp].i] = '{'
				b[i] = '}'
			}
			sp--
			if sp < -1 {
				return fmt.Errorf("underflow in fixFigureData")
			}
		}
	}
	return
}
