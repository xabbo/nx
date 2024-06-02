package render

import (
	"github.com/xabbo/nx/gamedata"
)

type FurniRenderer struct {
	mgr *gamedata.GamedataManager
}

type Furni struct {
	Identifier string
	Direction  int
	State      int
}

func NewFurniRenderer(mgr *gamedata.GamedataManager) *FurniRenderer {
	return &FurniRenderer{mgr}
}

func (r *FurniRenderer) Sprites(identifier string) (sprites []Sprite) {
	return
}
