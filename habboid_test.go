package nx

import (
	"bytes"
	"testing"
)

const hhidString = "g-hhus-00112233445566778899aabbccddeeff"
const hhidString2 = "g-hhus-00112233445566778899abbbccddeeff"
const hhidKind = HabboIdKindGroup
const hhidHotel = "us"

var hhidBytes = [16]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}

func TestHabboIdEqual(t *testing.T) {
	var a, b, c HabboId
	a.Parse(hhidString)
	b.Parse(hhidString)
	c.Parse(hhidString2)

	if a != b {
		t.Fatalf("%v should equal %v", a, b)
	}

	if a == c {
		t.Fatalf("%v should not equal %v", a, c)
	}
}

func TestHabboIdParse(t *testing.T) {
	var id HabboId
	err := id.Parse(hhidString)
	if err != nil {
		t.Fatal(err)
	}
	if id.Kind != hhidKind {
		t.Fatalf("Kind is %d (expected %d)", id.Kind, hhidKind)
	}
	if id.Hotel != hhidHotel {
		t.Fatalf("Hotel is %q (expected %q)", id.Hotel, hhidHotel)
	}
	if !bytes.Equal(hhidBytes[:], id.Uid[:]) {
		t.Fatalf("Uid is %v (expected %v)", id.Uid, hhidBytes)
	}
}

func TestHabboIdString(t *testing.T) {
	var id = HabboId{
		Kind:  HabboIdKindGroup,
		Hotel: "us",
		Uid:   hhidBytes,
	}
	s := id.String()
	if s != hhidString {
		t.Fatalf("%q should equal %q", s, hhidString)
	}
}
