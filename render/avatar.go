package render

import (
	"fmt"
	"image"
	"image/color"
	"slices"
	"strconv"

	"github.com/b7c/nx"
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

var layers = []nx.FigurePartType{
	nx.LeftHand,
	nx.LeftSleeve,
	nx.LeftCoat,
	nx.Body,
	nx.Shoes,
	nx.Legs,
	nx.Chest,
	nx.Waist,
	nx.Coat,
	nx.ChestAcc,
	nx.ChestPrint,
	nx.Head,
	nx.Face,
	nx.Eyes,
	nx.Hair,
	nx.HairBelow,
	nx.FaceAcc,
	nx.EyeAcc,
	nx.Hat,
	nx.HeadAcc,
	nx.LeftHandItem,
	nx.RightHandItem,
	nx.RightHand,
	nx.RightSleeve,
	nx.RightCoat,
}

type layerGroup []nx.FigurePartType

var bodyLayers = layerGroup{
	nx.Body,
	nx.Shoes,
	nx.Legs,
	nx.Chest,
	nx.ChestPrint,
	nx.Waist,
	nx.Coat,
	nx.ChestAcc,
}

var leftArmLayers = layerGroup{
	nx.LeftHand,
	nx.LeftSleeve,
	nx.LeftCoat,
}

var rightArmLayers = layerGroup{
	nx.RightHand,
	nx.RightSleeve,
	nx.RightCoat,
}

var handItemLayers = layerGroup{
	nx.LeftHandItem,
	nx.RightHandItem,
}

var headLayers = layerGroup{
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

var layerOrderUp = []layerGroup{
	handItemLayers,
	leftArmLayers,
	rightArmLayers,
	bodyLayers,
	headLayers,
}

var layerOrderDown = []layerGroup{
	bodyLayers,
	headLayers,
	handItemLayers,
	leftArmLayers,
	rightArmLayers,
}

var layerOrderSide = []layerGroup{
	leftArmLayers,
	bodyLayers,
	headLayers,
	rightArmLayers,
	handItemLayers,
}

type AvatarRenderer struct {
	mgr *nx.GamedataManager
}

func NewFigureRenderer(mgr *nx.GamedataManager) *AvatarRenderer {
	return &AvatarRenderer{mgr}
}

type AvatarPart struct {
	LibraryName string
	AssetSpec   FigureAssetSpec
	Asset       nx.Asset
	SetType     nx.FigurePartType
	SetId       int
	Type        nx.FigurePartType
	Id          int
	Color       color.Color
	Hidden      bool
}

// Converts the specified figure into individual figure parts.
func (r *AvatarRenderer) Parts(fig nx.Figure) (parts []AvatarPart, err error) {
	if !r.mgr.FigureLoaded() {
		err = fmt.Errorf("figure data not loaded")
		return
	}
	if !r.mgr.FigureMapLoaded() {
		err = fmt.Errorf("figure map not loaded")
		return
	}

	hiddenLayers := map[nx.FigurePartType]bool{}
	for _, partSet := range fig.Parts {
		setInfo := r.mgr.Figure.Sets[partSet.Type][partSet.Id]
		for _, layer := range setInfo.HiddenLayers {
			hiddenLayers[layer] = true
		}
	}

	for _, partSet := range fig.Parts {
		setInfo := r.mgr.Figure.Sets[partSet.Type][partSet.Id]
		palette := r.mgr.Figure.PaletteFor(partSet.Type)

		assumedLibrary := ""
		for _, partInfo := range setInfo.Parts {
			var c color.Color
			c = color.White
			if partInfo.Colorable && partInfo.Type != nx.Eyes {
				if partInfo.ColorIndex > 0 {
					cv, err := strconv.ParseInt(palette[partSet.Colors[partInfo.ColorIndex-1]].Value, 16, 64)
					if err != nil {
						return nil, err
					}
					c = color.RGBA{
						R: uint8((cv >> 16) & 0xff),
						G: uint8((cv >> 8) & 0xff),
						B: uint8(cv & 0xff),
						A: 0xff,
					}
				}
			}

			renderPart := AvatarPart{
				SetType: partSet.Type,
				SetId:   partSet.Id,
				Type:    partInfo.Type,
				Id:      partInfo.Id,
				Color:   c,
				Hidden:  hiddenLayers[partInfo.Type],
			}

			if lib, ok := r.mgr.FigureMap.Parts[nx.FigureMapPart{Type: partInfo.Type, Id: partInfo.Id}]; ok {
				renderPart.LibraryName = lib.Name
				assumedLibrary = lib.Name
			} else {
				// Some parts don't have a mapping for some reason,
				// so we use the library from previous parts in the same set
				if assumedLibrary == "" {
					err = fmt.Errorf("failed to find library for part %s-%d", partInfo.Type, partInfo.Id)
					return
				}
				renderPart.LibraryName = assumedLibrary
			}

			parts = append(parts, renderPart)
		}
	}

	return
}

// Finds the required figure part libraries given the specified Figure.
func (r *AvatarRenderer) RequiredLibs(fig nx.Figure) (libs []string, err error) {
	if !r.mgr.FigureLoaded() {
		err = fmt.Errorf("figure data not loaded")
		return
	}
	if !r.mgr.FigureMapLoaded() {
		err = fmt.Errorf("figure map not loaded")
		return
	}

	known := map[string]struct{}{}

	for _, part := range fig.Parts {
		setGroup, ok := r.mgr.Figure.Sets[part.Type]
		if !ok {
			err = fmt.Errorf("no figure part sets found for part type %q", part.Type)
			return
		}

		set, ok := setGroup[part.Id]
		if !ok {
			err = fmt.Errorf("no figure part set found for %s-%d", part.Type, part.Id)
		}

		for _, part := range set.Parts {
			mapPart := nx.FigureMapPart{
				Type: part.Type,
				Id:   part.Id,
			}
			partLib, ok := r.mgr.FigureMap.Parts[mapPart]
			if !ok {
				err = fmt.Errorf("part library not found for %s:%d", part.Type, part.Id)
				return
			}
			if _, exist := known[partLib.Name]; !exist {
				known[partLib.Name] = struct{}{}
				libs = append(libs, partLib.Name)
			}
		}
	}

	return
}

func isMirrored(dir int) bool {
	return dir >= 4 && dir <= 6
}

func flipDir(dir int) int {
	return (6 - dir) % 8
}

// Renders a figure to a list of sprites.
func (r *AvatarRenderer) Sprites(avatar nx.Avatar) (sprites []Sprite, err error) {
	parts, err := r.Parts(avatar.Figure)
	if err != nil {
		return
	}

	var ordering []layerGroup
	switch avatar.Direction {
	case 7:
		ordering = layerOrderUp
	case 0, 1, 2, 4, 5, 6:
		ordering = layerOrderSide
	case 3:
		ordering = layerOrderDown
	}

	n := 0
	layerOrder := map[nx.FigurePartType]int{}
	for _, group := range ordering {
		for _, layer := range group {
			layerOrder[layer] = n
			n++
		}
	}

	partMap := map[nx.FigurePartType][]AvatarPart{}
	// Groups parts by part type
	for _, part := range parts {
		// part.Order = layerOrder[part.Type]
		partMap[part.Type] = append(partMap[part.Type], part)
	}

	flipAvatar := avatar.Direction >= 4 && avatar.Direction <= 6

	// First pass over parts

	// Find flipped parts and replace if necessary.
	// most assets are just flipped when facing S-NW
	// however some parts have an asset for that direction
	// e.g. left arm wave when facing S-NW has an asset
	// since when facing S-NW, the assets are flipped
	// the left arm is visually the right arm
	// however, since there is a left arm asset
	// there will be 2 left arms, one is the flipped right arm
	// the flipped right arm must be removed

	// the right arm must be removed and replaced with
	// another left arm asset

	// the left arm asset must also be moved to the
	// right arm layer so that it is ordered correctly

	type partExtra struct {
		Spec   FigureAssetSpec
		Asset  nx.Asset
		Order  int
		Offset image.Point
		FlipH  bool
	}
	type partId struct {
		Type nx.FigurePartType
		Id   int
	}
	partExtraData := map[partId]partExtra{}

	for i := range parts {
		part := &parts[i]

		if avatar.HeadOnly && !part.Type.IsHead() {
			part.Hidden = true
		}

		if part.Hidden {
			continue
		}

		if !r.mgr.Assets.LibraryExists(part.LibraryName) {
			err = fmt.Errorf("required part library not loaded: %q", part.LibraryName)
			return
		}
		lib := r.mgr.Assets.Library(part.LibraryName)

		partDir := avatar.Direction
		isHead := part.Type.IsHead()
		if isHead {
			partDir = avatar.HeadDirection
		}

		flipPart := isMirrored(partDir)

		spec := r.ResolveAsset(lib, avatar, *part)
		if spec == nil {
			part.Hidden = true
			continue
		}

		var asset nx.Asset
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

		partExtraData[partId{part.Type, part.Id}] = partExtra{
			Asset:  asset,
			Spec:   *spec,
			Order:  layerOrder[part.Type],
			Offset: offset,
			FlipH:  flipPart,
		}
	}

	slices.SortFunc(parts, func(a, b AvatarPart) int {
		diff := partExtraData[partId{a.Type, a.Id}].Order - partExtraData[partId{b.Type, b.Id}].Order
		if diff == 0 {
			diff = a.Id - b.Id
		}
		return diff
	})

	// Convert parts into sprites
	for _, part := range parts {
		if part.Hidden {
			continue
		}
		extra := partExtraData[partId{part.Type, part.Id}]
		sprites = append(sprites, Sprite{
			Name:   extra.Spec.String(),
			Asset:  extra.Asset,
			Offset: extra.Offset,
			Color:  part.Color,
			FlipH:  extra.FlipH,
		})
	}

	return
}

func (r *AvatarRenderer) ResolveAsset(lib nx.AssetLibrary, avatar nx.Avatar, part AvatarPart) *FigureAssetSpec {
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

func (n FigureAssetSpec) String() string {
	return fmt.Sprintf("h_%s_%s_%d_%d_%d",
		n.State, n.Type, n.Id, n.Dir, n.Frame)
}

func (r *AvatarRenderer) Dependencies(fig nx.Figure) (err error) {
	return nil
}

func (r *AvatarRenderer) CompileAssets(fig nx.Figure) []Asset {
	return nil
}
