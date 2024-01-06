package nx

// Represents a furni's special type.
type FurniType int

const (
	FurniTypeNormal FurniType = iota + 1
	FurniTypeWallpaper
	FurniTypeFloor
	FurniTypeLandscape
	FurniTypeSticky
	FurniTypePoster
	FurniTypeTrax
	FurniTypeDisk
	FurniTypeGift
	FurniTypeMysteryBox
	FurniTypeTrophy
	FurniTypeHorseDye FurniType = iota + 2
	FurniTypeHorseHairstyle
	FurniTypeHorseHairdye
	FurniTypeHorseSaddle
	FurniTypeGroup
	FurniTypeSnowWar
	FurniTypeMonsterPlantSeed
	FurniTypeMonsterPlantRevival
	FurniTypeMonsterPlantRebreeding
	FurniTypeMonsterPlantFertiliser
	FurniTypeClothing
)
