package nx

// ItemType represents the type of an item.
// May be floor, wall, badge, effect or bot.
type ItemType rune

const (
	// Represents a floor item type.
	ItemFloor ItemType = 's'
	// Represents a wall item type.
	ItemWall ItemType = 'i'
	// Represents a badge item type.
	ItemBadge ItemType = 'b'
	// Represents an effect item type.
	ItemEffect ItemType = 'e'
	// Represents a bot item type.
	ItemBot ItemType = 'r'
)

// String returns the name of the item type.
func (t ItemType) String() string {
	switch t {
	case ItemFloor:
		return "Floor"
	case ItemWall:
		return "Wall"
	case ItemBadge:
		return "Badge"
	case ItemEffect:
		return "Effect"
	case ItemBot:
		return "Bot"
	default:
		return string(t)
	}
}
