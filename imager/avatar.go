package imager

import (
	"fmt"
	"image"
	"image/color"
	"slices"
	"strconv"

	"golang.org/x/exp/maps"
	"xabbo.b7c.io/nx"
	gd "xabbo.b7c.io/nx/gamedata"
	"xabbo.b7c.io/nx/res"
)

/*

Rendering requirements
- Figure part type/id, color
- Figure part libraries
- Figure direction, gesture, expression
- Figure part ordering
- Hand item ID
- Effect ID
- Sign ID

Rendering process
- Parse figure into part set type/id & colors
- Lookup figure data
	- Expand into a list of figure parts & colors
	- Exclude hidden part layers specified in the figure data for each part set
- Find and load required figure part libraries from figure map
- Sort figure part layers based on direction
- Resolve assets for each figure part
	- Certain assets fall back to other assets
	  e.g. wav assets only exist for the waving
	  arm/sleeve, then fall back to std for the body etc.
	- Some assets just don't exist for that part e.g. eye/face when facing away
*/

// An Avatar defines the state of a Figure in a room.
type Avatar struct {
	nx.Figure
	Direction     int
	HeadDirection int
	Action        nx.AvatarState
	Expression    nx.AvatarState
	HandItem      int
	Effect        int
	Sign          int
	HeadOnly      bool
}

var bodyLayers = []nx.FigurePartType{
	nx.Body,
	nx.Shoes,
	nx.Legs,
	nx.Chest,
	nx.ChestPrint,
	nx.Waist,
	nx.Coat,
	nx.ChestAcc,
}

var leftArmLayers = []nx.FigurePartType{
	nx.LeftHand,
	nx.LeftSleeve,
	nx.LeftCoat,
}

var rightArmLayers = []nx.FigurePartType{
	nx.RightHand,
	nx.RightSleeve,
	nx.RightCoat,
}

var handItemLayers = []nx.FigurePartType{
	nx.LeftHandItem,
	nx.RightHandItem,
}

var headLayers = []nx.FigurePartType{
	nx.Head,
	nx.Face,
	nx.Eyes,
	nx.Hair,
	nx.HairBelow,
	nx.FaceAcc,
	nx.EyeAcc,
	nx.Hat,
	nx.HeadAcc,
}

var layerOrderUp = [][]nx.FigurePartType{
	handItemLayers,
	leftArmLayers,
	rightArmLayers,
	bodyLayers,
	headLayers,
}

var layerOrderDown = [][]nx.FigurePartType{
	bodyLayers,
	headLayers,
	handItemLayers,
	leftArmLayers,
	rightArmLayers,
}

var layerOrderSide = [][]nx.FigurePartType{
	leftArmLayers,
	bodyLayers,
	headLayers,
	rightArmLayers,
	handItemLayers,
}

func isMirrored(dir int) bool {
	return dir >= 4 && dir <= 6
}

func flipDir(dir int) int {
	return (6 - dir) % 8
}

type avatarImager struct {
	mgr gd.Manager
}

func NewAvatarImager(mgr gd.Manager) AvatarImager {
	return avatarImager{mgr}
}

type AvatarPart struct {
	LibraryName string
	AssetSpec   FigureAssetSpec
	Asset       *res.Asset
	SetType     nx.FigurePartType
	SetId       int
	Type        nx.FigurePartType
	Id          int
	Color       color.Color
	Hidden      bool
}

// Converts the specified figure into individual figure parts.
func (imgr avatarImager) Parts(fig nx.Figure) (parts []AvatarPart, err error) {
	figureData := imgr.mgr.Figure()
	if figureData == nil {
		err = fmt.Errorf("figure data not loaded")
		return
	}
	figureMap := imgr.mgr.FigureMap()
	if figureMap == nil {
		err = fmt.Errorf("figure map not loaded")
		return
	}

	// Find all part layers that should be hidden.
	// Certain figure items may cause other layers to be hidden,
	// e.g. a hat may cause certain hair assets to be hidden.
	hiddenLayers := map[nx.FigurePartType]bool{}
	for _, item := range fig.Items {
		setInfo := figureData.Sets[item.Type][item.Id]
		for _, layer := range setInfo.HiddenLayers {
			hiddenLayers[layer] = true
		}
	}

	// Loop over each item (part set) in the figure.
	for _, item := range fig.Items {
		setInfo := figureData.Sets[item.Type][item.Id]
		palette := figureData.PaletteFor(item.Type)

		assumedLibrary := ""
		// Loop over each part in the figure item (part set).
		// Each item is comprised of multiple parts,
		// e.g. a shirt may have sprites for the body, left and right arms.
		for _, partInfo := range setInfo.Parts {
			// Resolve the color for this part.
			var col color.Color = color.White
			if partInfo.Colorable && partInfo.ColorIndex > 0 && partInfo.Type != nx.Eyes {
				if (partInfo.ColorIndex - 1) >= len(item.Colors) {
					err = fmt.Errorf("expected at least %d color(s) for part %s-%d", partInfo.ColorIndex, item.Type, item.Id)
					return
				}
				colorId := item.Colors[partInfo.ColorIndex-1]
				partColor, ok := palette[colorId]
				if !ok {
					err = fmt.Errorf("color %d not found for part %s", colorId, item.String())
					return
				}
				colorValue, err := strconv.ParseInt(partColor.Value, 16, 64)
				if err != nil {
					return nil, err
				}
				col = color.RGBA{
					R: uint8((colorValue >> 16) & 0xff),
					G: uint8((colorValue >> 8) & 0xff),
					B: uint8(colorValue & 0xff),
					A: 0xff,
				}
			}

			part := AvatarPart{
				SetType: item.Type,
				SetId:   item.Id,
				Type:    partInfo.Type,
				Id:      partInfo.Id,
				Color:   col,
				Hidden:  hiddenLayers[partInfo.Type],
			}

			// Resolve the figure library for this part.
			// Parts in the same part set may come from different figure libraries.
			if lib, ok := figureMap.Parts[nx.FigurePart{Type: partInfo.Type, Id: partInfo.Id}]; ok {
				part.LibraryName = lib.Name
				assumedLibrary = lib.Name
			} else {
				// Some parts don't have a mapping for some reason,
				// so we use the library from the previous part in the same set.
				if assumedLibrary == "" {
					err = fmt.Errorf("failed to find library for part %s-%d", partInfo.Type, partInfo.Id)
					return
				}
				part.LibraryName = assumedLibrary
			}

			parts = append(parts, part)
		}
	}

	return
}

// Finds the required figure part libraries given the specified Figure.
func (imgr avatarImager) RequiredLibs(fig nx.Figure) (libs []string, err error) {
	figureData := imgr.mgr.Figure()
	if figureData == nil {
		err = fmt.Errorf("figure data not loaded")
		return
	}
	figureMap := imgr.mgr.FigureMap()
	if figureMap == nil {
		err = fmt.Errorf("figure map not loaded")
		return
	}

	libSet := map[string]struct{}{}
	for _, item := range fig.Items {
		setGroup, ok := figureData.Sets[item.Type]
		if !ok {
			err = fmt.Errorf("no figure part sets found for part type %q", item.Type)
			return
		}

		partSet, ok := setGroup[item.Id]
		if !ok {
			err = fmt.Errorf("no figure part set found for %s-%d", item.Type, item.Id)
		}

		for _, partInfo := range partSet.Parts {
			partLib, ok := figureMap.Parts[nx.FigurePart{Type: partInfo.Type, Id: partInfo.Id}]
			if !ok {
				err = fmt.Errorf("part library not found for %s:%d", partInfo.Type, partInfo.Id)
				return
			}
			libSet[partLib.Name] = struct{}{}
		}
	}

	libs = maps.Keys(libSet)
	return
}

// Compose composes an avatar into an animation.
func (imgr avatarImager) Compose(avatar Avatar) (anim Animation, err error) {
	parts, err := imgr.Parts(avatar.Figure)
	if err != nil {
		return
	}

	// Choose a layer ordering based on figure direction.
	var ordering [][]nx.FigurePartType
	switch avatar.Direction {
	case 0, 1, 2, 4, 5, 6:
		ordering = layerOrderSide
	case 3:
		ordering = layerOrderDown
	case 7:
		ordering = layerOrderUp
	}

	// Map layer order by figure part type.
	n := 0
	layerOrder := map[nx.FigurePartType]int{}
	for _, group := range ordering {
		for _, layer := range group {
			layerOrder[layer] = n
			n++
		}
	}

	// Groups parts by part type.
	partMap := map[nx.FigurePartType][]AvatarPart{}
	for _, part := range parts {
		partMap[part.Type] = append(partMap[part.Type], part)
	}

	flipAvatar := avatar.Direction >= 4 && avatar.Direction <= 6

	// First pass over parts

	// Find flipped parts and replace if necessary.
	// Most assets are just flipped when facing left.
	// However, some parts have an asset for that direction
	// e.g. left arm wave has an asset for each direction.
	// The asset cannot simply be flipped as the waving arm would turn into the right arm.
	// However, there will now be 2 left arms, with one being a right arm asset.
	// The flipped right arm must be removed and replaced with a flipped left arm asset.
	// The left arm asset must also be moved to the right arm layer so that it is ordered correctly.

	type partExtra struct {
		Spec   FigureAssetSpec
		Asset  *res.Asset
		Order  int
		Offset image.Point
		FlipH  bool
	}
	partExtraData := map[nx.FigurePart]partExtra{}

	for i := range parts {
		part := &parts[i]

		if avatar.HeadOnly && !part.Type.IsHead() {
			part.Hidden = true
		}

		if part.Hidden {
			continue
		}

		if !imgr.mgr.LibraryExists(part.LibraryName) {
			err = fmt.Errorf("required part library not loaded: %q", part.LibraryName)
			return
		}
		lib := imgr.mgr.Library(part.LibraryName)

		partDir := avatar.Direction
		isHead := part.Type.IsHead()
		if isHead {
			partDir = avatar.HeadDirection
		}

		flipPart := isMirrored(partDir)

		spec := imgr.ResolveAsset(lib, avatar, *part)
		if spec == nil {
			part.Hidden = true
			continue
		}

		var asset *res.Asset
		asset, err = lib.Asset(spec.String())
		if err != nil {
			return
		}

		offset := asset.Offset
		if flipPart {
			offset.X = offset.X*-1 + asset.Image.Bounds().Dx() - 64
			if !flipAvatar && isHead {
				offset.X -= 3
			}
		} else if flipAvatar && isHead {
			offset.X += 3
		}

		partExtraData[nx.FigurePart{Type: part.Type, Id: part.Id}] = partExtra{
			Asset:  asset,
			Spec:   *spec,
			Order:  layerOrder[part.Type],
			Offset: offset,
			FlipH:  flipPart,
		}
	}

	slices.SortFunc(parts, func(a, b AvatarPart) int {
		partA := nx.FigurePart{Type: a.Type, Id: a.Id}
		partB := nx.FigurePart{Type: b.Type, Id: b.Id}
		diff := partExtraData[partA].Order - partExtraData[partB].Order
		if diff == 0 {
			diff = a.Id - b.Id
		}
		return diff
	})

	// Convert parts into sprites
	anim = Animation{
		Layers: map[int]AnimationLayer{},
	}

	layerId := 0
	for _, part := range parts {
		if part.Hidden {
			continue
		}
		extra := partExtraData[nx.FigurePart{Type: part.Type, Id: part.Id}]

		anim.Layers[layerId] = AnimationLayer{
			Frames: map[int]Frame{
				0: {
					Sprite{
						Asset:  extra.Asset,
						Offset: extra.Offset,
						Color:  part.Color,
						FlipH:  extra.FlipH,
						Alpha:  255,
					},
				},
			},
		}
		layerId++
	}

	return
}

func (r *avatarImager) ResolveAsset(lib res.AssetLibrary, avatar Avatar, part AvatarPart) *FigureAssetSpec {
	direction := avatar.Direction
	if part.Type.IsHead() {
		direction = avatar.HeadDirection
	}
	action := avatar.Action
	expression := avatar.Expression

	directions := []int{direction}
	if isMirrored(direction) {
		directions = append(directions, flipDir(direction))
	}

	states := []nx.AvatarState{}

	if part.Type.IsHead() {
		states = append(states, expression)
		if action == nx.ActLay {
			states = append(states, nx.ActLay)
		} else {
			states = append(states, nx.ActStand)
		}
	} else {
		states = append(states, action)
		switch action {
		case nx.ActWalk, nx.ActWave, nx.ActSit, nx.ActDrink:
			states = append(states, nx.ActStand)
		case nx.ActRespect, nx.ActCarry:
			states = append(states, nx.ActWave, nx.ActStand)
		case nx.ActBlowKiss:
			states = append(states, nx.ActDrink, nx.ActStand)
		case nx.ActSign:
			states = append(states, nx.ActWave, nx.ActStand)
		}
	}

	for _, d := range directions {
		for _, a := range states {
			spec := FigureAssetSpec{a, part.Type, part.Id, d, 0}
			assetName := spec.String()
			if !lib.AssetExists(assetName) {
				continue
			}
			return &spec
		}
	}

	return nil
}

type FigureAssetSpec struct {
	State nx.AvatarState
	Type  nx.FigurePartType
	Id    int
	Dir   int
	Frame int
}

func (spec FigureAssetSpec) String() string {
	return "h_" +
		string(spec.State) + "_" +
		string(spec.Type) + "_" +
		strconv.Itoa(spec.Id) + "_" +
		strconv.Itoa(spec.Dir) + "_" +
		strconv.Itoa(spec.Frame)
}
