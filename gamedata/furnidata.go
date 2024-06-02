package gamedata

import (
	"encoding/json"

	"github.com/b7c/nx"
	j "github.com/b7c/nx/json"
)

type FurniData map[string]FurniInfo

/*
type FurniData struct {
	identifier map[string]*FurniInfo
	typeKind map[ItemTypeKind]*FurniInfo
}

func (fd *FurniData) Identifier(identifier string) *FurniInfo{
	return fd.identifier[identifier]
}

func (fd *FurniData) Kind(kind ItemTypeKind) *FurniInfo {
	return fd.typeKind[kind]
}
*/

type FurniInfo struct {
	// A numeric identifier for the furni.
	// A floor and wall item may share the same kind.
	// This identifier also differs between hotels.
	Kind int
	// The type of the furni, whether it is a floor or wall item.
	Type nx.ItemType
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
	SpecialType     nx.FurniType
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
