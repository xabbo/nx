package gamedata

import (
	"reflect"

	"xabbo.b7c.io/nx/res"
)

type Type string

const (
	GameDataHashes    Type = "hashes"
	GameDataFurni     Type = "furnidata"
	GameDataProduct   Type = "productdata"
	GameDataVariables Type = "external_variables"
	GameDataTexts     Type = "external_texts"
	GameDataFigure    Type = "figurepartlist"
	GameDataFigureMap Type = "figuremap"
	GameDataAvatar    Type = "HabboAvatarActions"

	keyFlashClientUrl          = "flash.client.url"
	habboAvatarActionsFilename = "HabboAvatarActions.xml"
)

type Manager interface {
	FigureManager
	FurniManager
	ProductManager
	TextManager
	VariableManager
	// Loads the specified game data types.
	// If no types are specified, all types are loaded.
	Load(types... Type) error
	// Gets whether all of the specified game data types are loaded.
	Loaded(types... Type) bool
}

// Map of hashed game data types.
var hashTypeMap = map[Type]reflect.Type{
	GameDataFurni:     reflect.TypeOf((*FurniData)(nil)).Elem(),
	GameDataFigure:    reflect.TypeOf((*FigureData)(nil)).Elem(),
	GameDataProduct:   reflect.TypeOf((*ProductData)(nil)).Elem(),
	GameDataTexts:     reflect.TypeOf((*ExternalTexts)(nil)).Elem(),
	GameDataVariables: reflect.TypeOf((*ExternalVariables)(nil)).Elem(),
}

type FurniDataManager interface {
	Furni() FurniData // Gets the furni data.
}

type FurniLibraryManager interface {
	res.LibraryManager
	LoadFurni(libraries ...string) error
}

type FurniManager interface {
	FurniDataManager
	FurniLibraryManager
}

type FigureDataManager interface {
	Figure() *FigureData // Gets the figure data.
	FigureMap() *FigureMap // Gets the figure map.
	AvatarActions() AvatarActions // Gets the avatar actions.
}

type FigureLibraryManager interface {
	res.LibraryManager
	LoadFigureParts(libraries ...string) error
}

type FigureManager interface {
	FigureDataManager
	FigureLibraryManager
}

type ProductManager interface {
	Products() ProductData // Gets the products data.
}

type TextManager interface {
	Texts() ExternalTexts // Gets the external texts.
}

type VariableManager interface {
	Variables() ExternalVariables // Gets the external variables.
}

