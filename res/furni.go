package res

import (
	"strconv"

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
	Type           string                // The furni type for the visualization.
	Visualizations map[int]*Visualization // A map of visualizations by size.
}

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

type Visualization struct {
	Size       int
	LayerCount int
	Angle      int
	Directions map[int]struct{} // A map of directions.
	Layers     map[int]*Layer    // Layers mapped by ID.
	Colors     map[int]*Color    // Colors mapped by ID.
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

type Layer struct {
	Id          int
	Z           int
	Alpha       int
	Ink         string // Blend mode used for this layer.
	IgnoreMouse bool
}

func (layer *Layer) fromXml(v *x.Layer) {
	layer.Id = v.Id
	layer.Z = v.Z
	layer.Alpha = v.Alpha
	layer.Ink = v.Ink
	layer.IgnoreMouse = v.IgnoreMouse
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

type ColorLayer struct {
	Id    int
	Color string
}

func (colorLayer *ColorLayer) fromXml(v *x.ColorLayer) {
	colorLayer.Id = v.Id
	colorLayer.Color = v.Color
}

type Animation struct {
	Id     int
	TransitionTo *Animation
	Layers map[int]*AnimationLayer
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

// logic

type Logic struct {
	Type            string
	Model           Model
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
		ParticleSystems: make(map[int]*ParticleSystem, len(v.ParticleSystems)),
	}
	for _, xParticleSystem := range v.ParticleSystems {
		var particleSystem ParticleSystem
		particleSystem.fromXml(&xParticleSystem)
		logic.ParticleSystems[particleSystem.Size] = &particleSystem
	}
	logic.Model.fromXml(&v.Model)
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

type Dimensions struct {
	X int
	Y int
	Z float64
}

func (dimensions *Dimensions) fromXml(v *x.Dimensions) {
	*dimensions = Dimensions{v.X, v.Y, v.Z}
}

type ParticleSystem struct {
	Size int
}

func (particleSystem *ParticleSystem) fromXml(v *x.ParticleSystem) {
	*particleSystem = ParticleSystem{Size: v.Size}
}
