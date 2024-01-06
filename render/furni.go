package render

import "github.com/b7c/nx"

type FurniRenderer struct {
	mgr *nx.GamedataManager
}

type Furni struct {
	Identifier string
	Direction  int
	State      int
}

func NewFurniRenderer(mgr *nx.GamedataManager) *FurniRenderer {
	return &FurniRenderer{mgr}
}

func (r *FurniRenderer) Sprites(identifier string) (sprites []Sprite) {
	return
}
