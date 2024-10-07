package gamedata

import (
	"encoding/json"

	"xabbo.io/nx"
	j "xabbo.io/nx/raw/json"
)

// FurniData maps furniture info by identifier.
type FurniData map[string]*FurniInfo

// FurniInfo defines various information about a furniture.
type FurniInfo struct {
	// A numeric identifier for the furni.
	// A floor and wall item may share the same kind.
	// This identifier may differ between hotels.
	// Also known as the "Id" in the original document.
	// It is named "Kind" to differentiate it from a furni's unique instance ID.
	Kind int `json:"kind"`
	// The type of the furni, which may be a floor or wall item.
	Type nx.ItemType `json:"type"`
	// A unique string identifier for a kind of furniture.
	// This identifier is the same across hotels.
	// Also known as "ClassName" in the original document.
	Identifier      string       `json:"identifier"`
	Revision        int          `json:"revision"`
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	Category        string       `json:"category"`
	Environment     string       `json:"environment"`
	Line            string       `json:"line"`
	DefaultDir      int          `json:"defaultdir"`
	XDim            int          `json:"xdim"`
	YDim            int          `json:"ydim"`
	PartColors      []string     `json:"partcolors"`
	OfferId         int          `json:"offerid"`
	Buyout          bool         `json:"buyout"`
	BC              bool         `json:"bc"`
	ExcludedDynamic bool         `json:"excludeddynamic"`
	CustomParams    string       `json:"customparams"`
	SpecialType     nx.FurniType `json:"specialtype"`
	CanStandOn      bool         `json:"canstandon"`
	CanSitOn        bool         `json:"cansiton"`
	CanLayOn        bool         `json:"canlayon"`
}

// Unmarshals a JSON document as raw bytes into a FurniData.
func (fd *FurniData) UnmarshalBytes(data []byte) (err error) {
	jFurniData := j.FurniData{}
	err = json.Unmarshal(data, &jFurniData)
	if err != nil {
		return
	}

	*fd = FurniData{}
	for i := range jFurniData.FloorItems.Infos {
		jFurniInfo := &jFurniData.FloorItems.Infos[i]
		(*fd)[jFurniInfo.Identifier] = fromJsonFurniInfo(nx.ItemFloor, jFurniInfo)
	}
	for i := range jFurniData.WallItems.Infos {
		jFurniInfo := &jFurniData.WallItems.Infos[i]
		(*fd)[jFurniInfo.Identifier] = fromJsonFurniInfo(nx.ItemWall, jFurniInfo)
	}

	return
}

func fromJsonFurniInfo(furniType nx.ItemType, jfi *j.FurniInfo) *FurniInfo {
	return &FurniInfo{
		Type:            furniType,
		Kind:            jfi.Id,
		Identifier:      jfi.Identifier,
		Revision:        jfi.Revision,
		Name:            jfi.Name,
		Description:     jfi.Description,
		Category:        jfi.Category,
		Environment:     jfi.Environment,
		Line:            jfi.Line,
		DefaultDir:      jfi.DefaultDir,
		XDim:            jfi.XDim,
		YDim:            jfi.YDim,
		PartColors:      jfi.PartColors.Colors,
		OfferId:         jfi.OfferId,
		Buyout:          jfi.Buyout,
		BC:              jfi.BC,
		ExcludedDynamic: jfi.ExcludedDynamic,
		CustomParams:    jfi.CustomParams,
		SpecialType:     nx.FurniType(jfi.SpecialType),
		CanStandOn:      jfi.CanStandOn,
		CanSitOn:        jfi.CanSitOn,
		CanLayOn:        jfi.CanLayOn,
	}
}
