package gamedata

import (
	"encoding/json"

	"xabbo.b7c.io/nx"
	j "xabbo.b7c.io/nx/json"
)

// FurniData maps furniture info by identifier.
type FurniData map[string]FurniInfo

// FurniInfo defines various information about a furniture.
type FurniInfo struct {
	// A numeric identifier for the furni.
	// A floor and wall item may share the same kind.
	// This identifier may differ between hotels.
	// Also known as the "Id" in the original document.
	// It is named "Kind" to differentiate it from a furni's unique instance ID.
	Kind int
	// The type of the furni, which may be a floor or wall item.
	Type nx.ItemType
	// A unique string identifier for a kind of furniture.
	// This identifier is the same across hotels.
	// Also known as "ClassName" in the original document.
	Identifier      string
	Revision        int
	Name            string
	Description     string
	Category        string
	Environment     string
	Line            string
	DefaultDir      int
	XDim            int
	YDim            int
	PartColors      []string
	OfferId         int
	Buyout          bool
	BC              bool
	ExcludedDynamic bool
	CustomParams    string
	SpecialType     nx.FurniType
	CanStandOn      bool
	CanSitOn        bool
	CanLayOn        bool
}

// Unmarshals a JSON document as raw bytes into a FurniData.
func (fd *FurniData) UnmarshalBytes(data []byte) (err error) {
	jfd := j.FurniData{}
	err = json.Unmarshal(data, &jfd)
	if err != nil {
		return
	}

	*fd = FurniData{}
	for _, jfi := range jfd.FloorItems.Infos {
		(*fd)[jfi.Identifier] = fromJsonFurniInfo(nx.ItemFloor, jfi)
	}
	for _, jfi := range jfd.WallItems.Infos {
		(*fd)[jfi.Identifier] = fromJsonFurniInfo(nx.ItemWall, jfi)
	}

	return
}

func fromJsonFurniInfo(furniType nx.ItemType, jfi j.FurniInfo) FurniInfo {
	return FurniInfo{
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
