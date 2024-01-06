package nx

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

// Represents a unique Habbo resource identifier.
type HabboId struct {
	// The type of the identifier.
	Kind HabboIdKind
	// The related hotel identifier.
	Hotel string
	// A 128-bit unique identifier.
	Uid [16]byte
}

type HabboIdKind int

const (
	HabboIdKindUser HabboIdKind = iota
	HabboIdKindGroup
	HabboIdKindRoom
)

var rgxHabboId = regexp.MustCompile(`^(?:([rg])-)?hh([a-z]{2})-([0-9a-f]{32})$`)

func (t HabboIdKind) Prefix() string {
	switch t {
	case HabboIdKindUser:
		return ""
	case HabboIdKindGroup:
		return "g-"
	case HabboIdKindRoom:
		return "r-"
	default:
		return "?-"
	}
}

func (id *HabboId) String() string {
	return strings.Join([]string{
		id.Kind.Prefix(),
		"hh",
		id.Hotel,
		"-",
		hex.EncodeToString(id.Uid[:]),
	}, "")
}

func (id *HabboId) Parse(s string) (err error) {
	groups := rgxHabboId.FindStringSubmatch(s)
	if groups == nil || len(groups) != 4 {
		err = fmt.Errorf("invalid format: %q", s)
		return
	}
	var t HabboIdKind
	switch groups[1] {
	case "":
		t = HabboIdKindUser
	case "g":
		t = HabboIdKindGroup
	case "r":
		t = HabboIdKindRoom
	default:
		err = fmt.Errorf("unknown HabboIdKind: %q", groups[1])
		return
	}
	hotel := groups[2]
	bytes, err := hex.DecodeString(groups[3])
	if err != nil {
		return
	}
	hash := [16]byte(bytes)
	*id = HabboId{t, hotel, hash}
	return
}
