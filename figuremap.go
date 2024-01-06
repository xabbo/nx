package nx

import (
	"encoding/xml"
	"strconv"

	x "github.com/b7c/nx/xml"
)

type FigureMap struct {
	Libs  map[string]FigureMapLib        // Maps library name -> library.
	Parts map[FigureMapPart]FigureMapLib // Maps figure part -> library.
}

type FigureMapLib struct {
	Name     string
	Revision int
	Parts    []FigureMapPart
}

type FigureMapPart struct {
	Type FigurePartType
	Id   int
}

func (fm *FigureMap) UnmarshalBytes(data []byte) (err error) {
	var xfm x.FigureMap
	err = xml.Unmarshal(data, &xfm)
	if err != nil {
		return
	}

	*fm = FigureMap{
		Libs:  make(map[string]FigureMapLib),
		Parts: make(map[FigureMapPart]FigureMapLib),
	}

	for _, xlib := range xfm.Libraries {
		lib := FigureMapLib{
			Name:     xlib.Id,
			Revision: xlib.Revision,
			Parts:    make([]FigureMapPart, len(xlib.Parts)),
		}
		fm.Libs[lib.Name] = lib
		for _, xpart := range xlib.Parts {
			// A few parts in the hh_human_fx lib have non-numeric IDs, for now, we are ignoring them.
			// -- as of 3rd January, 2024
			if id, err := strconv.Atoi(xpart.Id); err == nil {
				part := FigureMapPart{
					Type: FigurePartType(xpart.Type),
					Id:   id,
				}
				lib.Parts = append(lib.Parts, part)
				// There is only one instance of duplicate part identifiers between
				// acc_eye_cyeyepiece / acc_eye_U_cyeyepiece, we just ignore it here.
				// -- as of 3rd January, 2024
				if _, exist := fm.Parts[part]; !exist {
					fm.Parts[part] = lib
				}
			}
		}
	}

	return
}
