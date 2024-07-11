package gamedata

import (
	"reflect"

	"xabbo.b7c.io/nx/res"
)

// Represents a type of game data.
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

// A Manager provides an interface to manage game data.
type Manager interface {
	FigureManager
	FurniManager
	ProductManager
	TextManager
	VariableManager
	// Loads the specified game data types.
	// If none are specified, all game data types are loaded.
	Load(types ...Type) error
	// Gets whether all of the specified game data types are loaded.
	Loaded(types ...Type) bool
}

// Map of hashed game data types.
var hashTypeMap = map[Type]reflect.Type{
	GameDataFurni:     reflect.TypeOf((*FurniData)(nil)).Elem(),
	GameDataFigure:    reflect.TypeOf((*FigureData)(nil)).Elem(),
	GameDataProduct:   reflect.TypeOf((*ProductData)(nil)).Elem(),
	GameDataTexts:     reflect.TypeOf((*ExternalTexts)(nil)).Elem(),
	GameDataVariables: reflect.TypeOf((*ExternalVariables)(nil)).Elem(),
}

// A FurniDataManager provides an interface to get furni data.
type FurniDataManager interface {
	Furni() FurniData // Gets the furni data.
}

// A FurniLibraryManager provides an interface to manage furni libraries.
type FurniLibraryManager interface {
	res.LibraryManager
	LoadFurni(libraries ...string) error // Loads the specified furni libraries by name.
}

// A FurniManager provides an interface to manage furni data and libraries.
type FurniManager interface {
	FurniDataManager
	FurniLibraryManager
}

// A FigureDataManager provides an interface to get figure data, figure map and avatar actions.
type FigureDataManager interface {
	Figure() *FigureData          // Gets the figure data.
	FigureMap() *FigureMap        // Gets the figure map.
	AvatarActions() AvatarActions // Gets the avatar actions.
}

// A FigureLibraryManager provides an interface to manage figure part libraries.
type FigureLibraryManager interface {
	res.LibraryManager
	LoadFigureParts(libraries ...string) error
}

// A FigureManager provides an interface to manage figure data and libraries.
type FigureManager interface {
	FigureDataManager
	FigureLibraryManager
}

// A ProductManager provides an interface to get product data.
type ProductManager interface {
	Products() ProductData // Gets the products data.
}

// A TextManager provides an interface to get external texts.
type TextManager interface {
	Texts() ExternalTexts // Gets the external texts.
}

// A VariableManager provides an interface to get external variables.
type VariableManager interface {
	Variables() ExternalVariables // Gets the external variables.
}
