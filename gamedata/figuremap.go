package gamedata

import (
	"encoding/xml"
	"strconv"

	"xabbo.io/nx"
	x "xabbo.io/nx/raw/xml"
)

// A FigureMap defines mappings between figure part libraries and figure part identifiers.
type FigureMap struct {
	Libs  map[string]*FigureMapLib        // Maps library name -> library.
	Parts map[nx.FigurePart]*FigureMapLib // Maps figure part -> library.
}

// A FigureMapLib defines a figure part library name, revision and the figure parts contained within the library.
type FigureMapLib struct {
	Name     string
	Revision int
	Parts    []nx.FigurePart
}

// Unmarshals an XML document as raw bytes into a FigureMap.
func (fm *FigureMap) UnmarshalBytes(data []byte) (err error) {
	var xfm x.FigureMap
	err = xml.Unmarshal(data, &xfm)
	if err != nil {
		return
	}

	*fm = FigureMap{
		Libs:  make(map[string]*FigureMapLib),
		Parts: make(map[nx.FigurePart]*FigureMapLib),
	}

	for _, xlib := range xfm.Libraries {
		lib := &FigureMapLib{
			Name:     xlib.Id,
			Revision: xlib.Revision,
			Parts:    make([]nx.FigurePart, len(xlib.Parts)),
		}
		fm.Libs[lib.Name] = lib
		for _, xpart := range xlib.Parts {
			// A few parts in the hh_human_fx lib have non-numeric IDs, for now, we are ignoring them.
			if id, err := strconv.Atoi(xpart.Id); err == nil {
				part := nx.FigurePart{
					Type: nx.FigurePartType(xpart.Type),
					Id:   id,
				}
				lib.Parts = append(lib.Parts, part)
				// There is only one instance of duplicate part identifiers between
				// acc_eye_cyeyepiece / acc_eye_U_cyeyepiece, we just ignore it here.
				if _, exist := fm.Parts[part]; !exist {
					fm.Parts[part] = lib
				}
			}
		}
	}

	return
}
