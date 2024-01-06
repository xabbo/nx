package nx

import (
	"fmt"
	"strconv"
	"strings"
)

// Defines the visual representation of a figure.
type Figure struct {
	Gender Gender
	Parts  []FigurePart
}

// Defines the visual representation of an avatar in a room.
type Avatar struct {
	Figure
	Direction     int
	HeadDirection int
	Action        AvatarState
	Expression    AvatarState
	HandItem      int
	Effect        int
	Sign          int
	HeadOnly      bool
}

func (f *Figure) String() string {
	sb := strings.Builder{}
	for i, part := range f.Parts {
		if i > 0 {
			sb.WriteRune('.')
		}
		part.writeBuilder(&sb)
	}
	return sb.String()
}

// Represents a part set identifier & color in a Habbo figure.
type FigurePart struct {
	Type   FigurePartType
	Id     int
	Colors []int
}

func (p *FigurePart) writeBuilder(sb *strings.Builder) {
	sb.WriteString(string(p.Type))
	sb.WriteRune('-')
	sb.WriteString(strconv.Itoa(p.Id))
	for i := range p.Colors {
		sb.WriteRune('-')
		sb.WriteString(strconv.Itoa(p.Colors[i]))
	}
}

func (p *FigurePart) String() string {
	sb := strings.Builder{}
	p.writeBuilder(&sb)
	return sb.String()
}

type FigurePartType string

const (
	Hair          FigurePartType = "hr"  // Hair.
	HairBelow     FigurePartType = "hrb" // Hair below hat.
	Head          FigurePartType = "hd"  // Head.
	Hat           FigurePartType = "ha"  // Hat.
	HeadAcc       FigurePartType = "he"  // Head accessory.
	EyeAcc        FigurePartType = "ea"  // Eye accessory, i.e. glasses.
	FaceAcc       FigurePartType = "fa"  // Face accessory, i.e. masks.
	Eyes          FigurePartType = "ey"  // Eyes.
	Face          FigurePartType = "fc"  // Face.
	Body          FigurePartType = "bd"  // Body.
	LeftHand      FigurePartType = "lh"  // Left hand.
	RightHand     FigurePartType = "rh"  // Right hand.
	Chest         FigurePartType = "ch"  // Chest, i.e. shirts.
	ChestPrint    FigurePartType = "cp"  // Chest print.
	ChestAcc      FigurePartType = "ca"  // Chest accessory, i.e. jewellery.
	LeftSleeve    FigurePartType = "ls"  // Left sleeve.
	RightSleeve   FigurePartType = "rs"  // Right sleeve.
	Legs          FigurePartType = "lg"  // Legs, i.e. trousers.
	Shoes         FigurePartType = "sh"  // Shoes.
	Waist         FigurePartType = "wa"  // Waist, i.e. belts.
	Coat          FigurePartType = "cc"  // Coat/jacket.
	LeftCoat      FigurePartType = "lc"  // Left coat sleeve.
	RightCoat     FigurePartType = "rc"  // Right coat sleeve.
	LeftHandItem  FigurePartType = "li"  // Left hand item.
	RightHandItem FigurePartType = "ri"  // Right hand item.
)

// Gets whether the part type belongs to the head.
func (pt FigurePartType) IsHead() bool {
	switch pt {
	case Hair, HairBelow, Head, Hat, HeadAcc, EyeAcc, FaceAcc, Eyes, Face:
		return true
	default:
		return false
	}
}

// Gets whether the part type belongs to the left arm.
func (pt FigurePartType) IsLeft() bool {
	switch pt {
	case LeftHand, LeftSleeve, LeftCoat, LeftHandItem:
		return true
	default:
		return false
	}
}

// Gets whether the part type belongs to the right arm.
func (pt FigurePartType) IsRight() bool {
	switch pt {
	case RightHand, RightSleeve, RightCoat, RightHandItem:
		return true
	default:
		return false
	}
}

// Flips the part type between left/right arms, if it is an arm part.
func (pt FigurePartType) Flip() FigurePartType {
	switch {
	case pt.IsLeft():
		return FigurePartType("r" + string(pt[1]))
	case pt.IsRight():
		return FigurePartType("l" + string(pt[1]))
	default:
		return pt
	}
}

// Gets if the figure part type is valid in a figure string.
func (layer FigurePartType) Wearable() (wearable bool) {
	switch layer {
	case Head, Hair, Hat, HeadAcc, EyeAcc, FaceAcc, Chest, ChestPrint, Coat, ChestAcc, Legs, Shoes, Waist:
		wearable = true
	}
	return
}

type AvatarState string

const (
	ActStand    AvatarState = "std"     // Standing.
	ActWalk     AvatarState = "wlk"     // Walking.
	ActWave     AvatarState = "wav"     // Waving.
	ActLay      AvatarState = "lay"     // Laying.
	ActBlowKiss AvatarState = "blw"     // Blowing a kiss.
	ActCarry    AvatarState = "crr"     // Carrying a hand item.
	ActDrink    AvatarState = "drk"     // Drinking.
	ActRespect  AvatarState = "respect" // Respecting.
	ActSign     AvatarState = "sig"     // Showing a sign.
	ActSit      AvatarState = "sit"     // Sitting.
)

var AllActions = []AvatarState{
	ActStand, ActWalk, ActWave, ActLay, ActBlowKiss,
	ActCarry, ActDrink, ActRespect, ActSign, ActSit,
}

func (state AvatarState) IsAction() bool {
	switch state {
	case ActStand, ActWalk, ActWave, ActLay, ActBlowKiss, ActCarry, ActDrink, ActRespect, ActSign, ActSit:
		return true
	default:
		return false
	}
}

const (
	// Used to specify a neutral / no expression.
	// It is not an official expression, and is not included in AllExpressions.
	ExprNeutral      AvatarState = "ntr"
	ExprSpeak        AvatarState = "spk" // Speaking.
	ExprSleep        AvatarState = "eyb" // Sleeping.
	ExprSad          AvatarState = "sad" // Sad.
	ExprSmile        AvatarState = "sml" // Smiling.
	ExprAngry        AvatarState = "agr" // Angry.
	ExprSurprised    AvatarState = "srp" // Surprised.
	ExprSpeakLay     AvatarState = "lsp" // Speaking / laying.
	ExprSleepLay     AvatarState = "ley" // Sleeping / laying.
	ExprSadLay       AvatarState = "lsa" // Sad / laying.
	ExprSmileLay     AvatarState = "lsm" // Smiling / laying.
	ExprAngryLay     AvatarState = "lag" // Angry / laying.
	ExprSurprisedLay AvatarState = "lsr" // Surprised / laying.
)

var AllExpressions = []AvatarState{
	ExprSpeak, ExprSleep, ExprSad, ExprSmile,
	ExprAngry, ExprSurprised, ExprSpeakLay, ExprSleepLay,
	ExprSadLay, ExprSmileLay, ExprAngryLay, ExprSurprisedLay,
}

func (state AvatarState) IsExpression() bool {
	switch state {
	case ExprNeutral, ExprSpeak, ExprSleep, ExprSad, ExprSmile, ExprAngry, ExprSurprised, ExprSpeakLay, ExprSleepLay, ExprSadLay, ExprSmileLay, ExprAngryLay, ExprSurprisedLay:
		return true
	default:
		return false
	}
}

func (f *Figure) Parse(figure string) (err error) {
	split := strings.Split(figure, ".")
	parts := make([]FigurePart, 0, len(split))
	for _, partStr := range split {
		part := FigurePart{}
		partSplit := strings.Split(partStr, "-")
		if len(partSplit) < 1 || partSplit[0] == "" {
			return fmt.Errorf("empty figure part in figure string")
		}
		part.Type = FigurePartType(partSplit[0])
		if !part.Type.Wearable() {
			return fmt.Errorf("non-wearable figure part type %q", part.Type)
		}
		if len(partSplit) < 2 {
			return fmt.Errorf("unspecified ID for figure part type %q", part.Type)
		}
		part.Id, err = strconv.Atoi(partSplit[1])
		if err != nil {
			return fmt.Errorf("invalid figure part id %q", partSplit[1])
		}
		part.Colors = make([]int, len(partSplit)-2)
		for i, colorStr := range partSplit[2:] {
			color, err := strconv.Atoi(colorStr)
			if err != nil {
				return fmt.Errorf("invalid figure part color %q in %q", colorStr, partStr)
			}
			part.Colors[i] = color
		}
		parts = append(parts, part)
	}
	*f = Figure{Parts: parts, Gender: Unisex}
	return
}
