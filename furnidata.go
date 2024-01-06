package nx

import (
	"encoding/json"

	j "github.com/b7c/nx/json"
)

type FurniData map[string]FurniInfo

type FurniInfo struct {
	// The kind of the furni, i.e. its numeric identifier.
	// This identifier differs between hotels.
	Kind int
	// The type of the furni, whether it is a floor or wall item.
	Type ItemType
	// The identifier of the furni, also known as its class name.
	// This uniquely identifies furni and is the same across hotels.
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
	SpecialType     FurniType
	CanStandOn      bool
	CanSitOn        bool
	CanLayOn        bool
}

func (fd *FurniData) UnmarshalBytes(data []byte) (err error) {
	jfd := j.FurniData{}
	err = json.Unmarshal(data, &jfd)
	if err != nil {
		return
	}

	*fd = FurniData{}
	for _, jfi := range jfd.FloorItems.Infos {
		(*fd)[jfi.Identifier] = fromJsonFurniInfo(ItemFloor, jfi)
	}
	for _, jfi := range jfd.WallItems.Infos {
		(*fd)[jfi.Identifier] = fromJsonFurniInfo(ItemWall, jfi)
	}

	return
}

func fromJsonFurniInfo(furniType ItemType, jfi j.FurniInfo) FurniInfo {
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
		SpecialType:     FurniType(jfi.SpecialType),
		CanStandOn:      jfi.CanStandOn,
		CanSitOn:        jfi.CanSitOn,
		CanLayOn:        jfi.CanLayOn,
	}
}
