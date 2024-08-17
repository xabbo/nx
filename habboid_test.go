package nx

import (
	"bytes"
	"fmt"
	"testing"
)

const hhidString = "g-hhus-00112233445566778899aabbccddeeff"
const hhidString2 = "g-hhus-00112233445566778899aabbcc00ffee"

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
	for kind, prefix := range map[HabboIdKind]string{
		HabboIdKindUser:  "",
		HabboIdKindGroup: "g-",
		HabboIdKindRoom:  "r-",
	} {
		for _, hotel := range []string{"us", "nl", "ous"} {
			strId := fmt.Sprintf("%shh%s-00112233445566778899aabbccddeeff", prefix, hotel)
			t.Run(strId, func(t *testing.T) {
				var id HabboId
				err := id.Parse(strId)
				if err != nil {
					t.Fatal(err)
				}
				if id.Kind != kind {
					t.Fatalf("Kind is %s (expected %s)", id.Kind.Prefix(), kind.Prefix())
				}
				if id.Hotel != hotel {
					t.Fatalf("Hotel is %q (expected %q)", id.Hotel, hotel)
				}
				if !bytes.Equal(hhidBytes[:], id.Uid[:]) {
					t.Fatalf("Uid is %v (expected %v)", id.Uid, hhidBytes)
				}
			})
		}
	}
}

func TestHabboIdString(t *testing.T) {
	for kind, prefix := range map[HabboIdKind]string{
		HabboIdKindUser:  "",
		HabboIdKindGroup: "g-",
		HabboIdKindRoom:  "r-",
	} {
		for _, hotel := range []string{"us", "nl", "ous"} {
			expected := fmt.Sprintf("%shh%s-00112233445566778899aabbccddeeff", prefix, hotel)
			t.Run(expected, func(t *testing.T) {
				actual := (&HabboId{kind, hotel, hhidBytes}).String()
				if actual != expected {
					t.Fatalf("actual: %q expected: %q", actual, expected)
				}
			})
		}
	}
}
