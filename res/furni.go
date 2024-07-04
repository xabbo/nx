package res

import (
	"strconv"

	"xabbo.b7c.io/nx/raw/nitro"
	x "xabbo.b7c.io/nx/raw/xml"
)

type FurniAssetSpec struct {
	Name      string
	Size      int
	Layer     int
	Direction int
	Frame     int
}

func (spec *FurniAssetSpec) String() string {
	return spec.Name + "_" +
		strconv.Itoa(spec.Size) + "_" +
		string('a'+rune(spec.Layer)) + "_" +
		strconv.Itoa(spec.Direction) + "_" +
		strconv.Itoa(spec.Frame)
}

// visualization

type VisualizationData struct {
	Type           string         // The furni type for the visualization.
	Visualizations Visualizations // A map of visualizations by size.
}

type Visualizations map[int]*Visualization

func (visualizationData *VisualizationData) UnmarshalBytes(b []byte) (err error) {
	var xVisData x.VisualizationData
	err = decodeXml(b, &xVisData)
	if err != nil {
		return
	}

	visualizationData.fromXml(&xVisData)
	return
}

func (visualizationData *VisualizationData) fromXml(v *x.VisualizationData) {
	*visualizationData = VisualizationData{
		Type:           v.Type,
		Visualizations: make(map[int]*Visualization, len(v.Graphics.Visualizations)),
	}
	for _, xVisualization := range v.Graphics.Visualizations {
		var visualization Visualization
		visualization.fromXml(&xVisualization)
		visualizationData.Visualizations[visualization.Size] = &visualization
	}
}

func (visualizations Visualizations) fromNitro(v []nitro.Visualization) Visualizations {
	visualizations = Visualizations{}
	for i := range v {
		visualizations[v[i].Size] = new(Visualization).fromNitro(&v[i])
	}
	return visualizations
}

type Visualization struct {
	Size       int
	LayerCount int
	Angle      int
	Directions map[int]struct{} // A map of directions.
	Layers     map[int]*Layer   // Layers mapped by ID.
	Colors     map[int]*Color   // Colors mapped by ID.
	Animations map[int]*Animation
	// TODO postures, gestures
}

func (vis *Visualization) fromXml(v *x.Visualization) {
	vis.Size = v.Size
	vis.LayerCount = v.LayerCount
	vis.Angle = v.Angle

	vis.Directions = make(map[int]struct{}, len(v.Directions))
	for _, dir := range v.Directions {
		vis.Directions[dir.Id] = struct{}{}
	}

	vis.Layers = make(map[int]*Layer, len(v.Layers))
	for i := range v.Layers {
		var layer Layer
		layer.fromXml(&v.Layers[i])
		vis.Layers[layer.Id] = &layer
	}

	vis.Colors = make(map[int]*Color, len(v.Colors))
	for i := range v.Colors {
		var color Color
		color.fromXml(&v.Colors[i])
		vis.Colors[color.Id] = &color
	}

	vis.Animations = make(map[int]*Animation, len(v.Animations))
	transitions := map[int]int{}
	for i := range v.Animations {
		xAnim := &v.Animations[i]
		var anim Animation
		anim.fromXml(xAnim)
		vis.Animations[anim.Id] = &anim
		if xAnim.TransitionTo != nil {
			transitions[xAnim.Id] = *xAnim.TransitionTo
		}
	}
	for from, to := range transitions {
		vis.Animations[from].TransitionTo = vis.Animations[to]
	}
}

func (vis *Visualization) fromNitro(v *nitro.Visualization) *Visualization {
	vis.Size = v.Size
	vis.LayerCount = v.LayerCount
	vis.Angle = v.Angle

	vis.Directions = make(map[int]struct{}, len(v.Directions))
	for dir := range v.Directions {
		vis.Directions[dir] = struct{}{}
	}

	vis.Layers = make(map[int]*Layer, len(v.Layers))
	for id, layer := range v.Layers {
		vis.Layers[id] = new(Layer).fromNitro(id, layer)
	}

	vis.Colors = make(map[int]*Color, len(v.Colors))
	for id, color := range v.Colors {
		vis.Colors[id] = new(Color).fromNitro(id, color)
	}

	vis.Animations = make(map[int]*Animation, len(v.Animations))
	transitions := map[int]int{}
	for id, animation := range v.Animations {
		vis.Animations[id] = new(Animation).fromNitro(id, animation)
		if animation.TransitionTo != nil {
			transitions[id] = *animation.TransitionTo
		}
	}
	for from, to := range transitions {
		vis.Animations[from].TransitionTo = vis.Animations[to]
	}

	return vis
}

type Layer struct {
	Id          int
	Z           int
	Alpha       int
	Ink         string // Blend mode used for this layer.
	IgnoreMouse bool
	Color       int
}

func (layer *Layer) fromXml(v *x.Layer) {
	layer.Id = v.Id
	layer.Z = v.Z
	layer.Alpha = v.Alpha
	layer.Ink = v.Ink
	layer.IgnoreMouse = v.IgnoreMouse
	layer.Color = v.Color
}

func (layer *Layer) fromNitro(id int, v nitro.Layer) *Layer {
	layer.Id = id
	layer.Z = v.Z
	layer.Alpha = v.Alpha
	layer.Ink = v.Ink
	layer.IgnoreMouse = v.IgnoreMouse
	layer.Color = v.Color
	return layer
}

type Color struct {
	Id     int
	Layers map[int]*ColorLayer
}

func (color *Color) fromXml(v *x.Color) {
	color.Id = v.Id
	color.Layers = make(map[int]*ColorLayer, len(v.Layers))
	for i := range v.Layers {
		var colorLayer ColorLayer
		colorLayer.fromXml(&v.Layers[i])
		color.Layers[colorLayer.Id] = &colorLayer
	}
}

func (color *Color) fromNitro(id int, v nitro.Color) *Color {
	*color = Color{
		Id:     id,
		Layers: make(map[int]*ColorLayer, len(v.Layers)),
	}
	for id, layer := range v.Layers {
		color.Layers[id] = new(ColorLayer).fromNitro(id, layer)
	}
	return color
}

type ColorLayer struct {
	Id    int
	Color string
}

func (colorLayer *ColorLayer) fromXml(v *x.ColorLayer) {
	colorLayer.Id = v.Id
	colorLayer.Color = v.Color
}

func (colorLayer *ColorLayer) fromNitro(id int, v nitro.Layer) *ColorLayer {
	*colorLayer = ColorLayer{Id: id}
	colorLayer.Color = strconv.FormatInt(int64(v.Color), 16)
	return colorLayer
}

type Animation struct {
	Id           int
	TransitionTo *Animation
	Layers       map[int]*AnimationLayer
}

func (anim *Animation) fromXml(v *x.Animation) {
	anim.Id = v.Id
	anim.Layers = make(map[int]*AnimationLayer, len(v.Layers))
	for i := range v.Layers {
		var animLayer AnimationLayer
		animLayer.fromXml(&v.Layers[i])
		anim.Layers[animLayer.Id] = &animLayer
	}
}

func (anim *Animation) fromNitro(id int, v nitro.Animation) *Animation {
	*anim = Animation{
		Id:     id,
		Layers: make(map[int]*AnimationLayer, len(v.Layers)),
	}
	for id, layer := range v.Layers {
		anim.Layers[id] = new(AnimationLayer).fromNitro(id, layer)
	}
	return anim
}

type AnimationLayer struct {
	Id             int
	LoopCount      int
	FrameRepeat    int
	Random         bool
	FrameSequences []FrameSequence // A list of frame sequences.
}

// FrameSequence represents a list of frame IDs for an animation.
type FrameSequence []int

func (animLayer *AnimationLayer) fromXml(v *x.AnimationLayer) {
	animLayer.Id = v.Id
	animLayer.LoopCount = v.LoopCount
	animLayer.FrameRepeat = v.FrameRepeat
	animLayer.Random = v.Random
	animLayer.FrameSequences = make([]FrameSequence, 0, len(v.FrameSequences))
	for _, xSequence := range v.FrameSequences {
		sequence := make(FrameSequence, 0, len(xSequence.Frames))
		for _, xFrame := range xSequence.Frames {
			sequence = append(sequence, xFrame.Id)
		}
		animLayer.FrameSequences = append(animLayer.FrameSequences, sequence)
	}
}

func (layer *AnimationLayer) fromNitro(id int, v nitro.AnimationLayer) *AnimationLayer {
	*layer = AnimationLayer{
		Id:             id,
		LoopCount:      v.LoopCount,
		FrameRepeat:    v.FrameRepeat,
		Random:         v.Random != 0,
		FrameSequences: make([]FrameSequence, 0, len(v.FrameSequences)),
	}
	for _, srcSeq := range v.FrameSequences {
		frameSequence := make(FrameSequence, 0, len(srcSeq.Frames))
		for _, srcFrame := range srcSeq.Frames {
			frameSequence = append(frameSequence, srcFrame.Id)
		}
		layer.FrameSequences = append(layer.FrameSequences, frameSequence)
	}
	return layer
}

// logic

type Logic struct {
	Type            string
	Model           *Model
	ParticleSystems map[int]*ParticleSystem // Particle systems mapped by size.
}

func (logic *Logic) UnmarshalBytes(b []byte) (err error) {
	var xLogic x.Logic
	err = decodeXml(b, &xLogic)
	if err != nil {
		return
	}

	logic.fromXml(&xLogic)
	return
}

func (logic *Logic) fromXml(v *x.Logic) {
	*logic = Logic{
		Type:            v.Type,
		Model:           new(Model).fromXml(&v.Model),
		ParticleSystems: make(map[int]*ParticleSystem, len(v.ParticleSystems)),
	}
	for _, xParticleSystem := range v.ParticleSystems {
		var particleSystem ParticleSystem
		particleSystem.fromXml(&xParticleSystem)
		logic.ParticleSystems[particleSystem.Size] = &particleSystem
	}
}

func (logic *Logic) fromNitro(v *nitro.Logic) *Logic {
	*logic = Logic{
		Model:           new(Model).fromNitro(v.Model),
		ParticleSystems: make(map[int]*ParticleSystem, len(v.ParticleSystems)),
	}
	for _, srcParticleSystem := range v.ParticleSystems {
		logic.ParticleSystems[srcParticleSystem.Size] = new(ParticleSystem).fromNitro(srcParticleSystem)
	}
	return logic
}

type ParticleSystem struct {
	Size int
}

func (particleSystem *ParticleSystem) fromXml(v *x.ParticleSystem) {
	*particleSystem = ParticleSystem{Size: v.Size}
}

func (particleSystem *ParticleSystem) fromNitro(v nitro.ParticleSystem) *ParticleSystem {
	*particleSystem = ParticleSystem{Size: v.Size}
	return particleSystem
}

type Model struct {
	Dimensions Dimensions
	Directions []int
}

func (model *Model) fromXml(v *x.Model) *Model {
	*model = Model{
		Directions: make([]int, 0, len(v.Directions)),
	}
	model.Dimensions.fromXml(&v.Dimensions)
	for _, xDir := range v.Directions {
		model.Directions = append(model.Directions, xDir.Id)
	}
	return model
}

func (model *Model) fromNitro(v nitro.Model) *Model {
	*model = Model{
		Dimensions: *new(Dimensions).fromNitro(v.Dimensions),
		Directions: make([]int, 0, len(v.Directions)),
	}
	for _, srcDir := range v.Directions {
		model.Directions = append(model.Directions, srcDir)
	}
	return model
}

type Dimensions struct {
	X int
	Y int
	Z float64
}

func (dimensions *Dimensions) fromXml(v *x.Dimensions) {
	*dimensions = Dimensions{v.X, v.Y, v.Z}
}

func (dimensions *Dimensions) fromNitro(v nitro.Dimensions) *Dimensions {
	*dimensions = Dimensions{v.X, v.Y, v.Z}
	return dimensions
}
