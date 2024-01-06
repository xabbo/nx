package nx

type ItemType rune

const (
	ItemFloor  ItemType = 's'
	ItemWall   ItemType = 'i'
	ItemBadge  ItemType = 'b'
	ItemEffect ItemType = 'e'
	ItemBot    ItemType = 'r'
)

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
