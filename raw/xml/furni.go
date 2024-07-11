package xml

import (
	"encoding/xml"
	"strconv"
)

// index.xml

type Index struct {
	Type          string `xml:"type,attr"`
	Visualization string `xml:"visualization,attr"`
	Logic         string `xml:"logic"`
}

// logic.xml

type Logic struct {
	Type            string           `xml:"type,attr"`
	Model           Model            `xml:"model"`
	ParticleSystems []ParticleSystem `xml:"particlesystems>particlesystem"`
}

type Model struct {
	Dimensions Dimensions  `xml:"dimensions"`
	Directions []Direction `xml:"directions>direction"`
}

type Dimensions struct {
	X float64 `xml:"x,attr"`
	Y float64 `xml:"y,attr"`
	Z float64 `xml:"z,attr"`
}

type Direction struct {
	Id int `xml:"id,attr"`
}

type ParticleSystem struct {
	Size int `xml:"size,attr"`
}

// visualization.xml

type VisualizationData struct {
	Type     string   `xml:"type,attr"`
	Graphics Graphics `xml:"graphics"`
}

type Graphics struct {
	Visualizations []Visualization `xml:"visualization"`
}

type Visualization struct {
	Size       int         `xml:"size,attr"`
	LayerCount int         `xml:"layerCount,attr"`
	Angle      int         `xml:"angle,attr"`
	Layers     []Layer     `xml:"layers>layer"`
	Directions []Direction `xml:"directions>direction"`
	Colors     []Color     `xml:"colors>color"`
	Animations []Animation `xml:"animations>animation"`
	// Postures
	// Gestures
}

type Layer struct {
	Id          int    `xml:"id,attr"`
	Z           int    `xml:"z,attr"`
	Alpha       int    `xml:"alpha,attr"`
	Ink         string `xml:"ink,attr"`
	IgnoreMouse bool   `xml:"ignoreMouse,attr"`
	Color       int    `xml:"color,attr"`
}

type Color struct {
	Id     int          `xml:"id,attr"`
	Layers []ColorLayer `xml:"colorLayer"`
}

type ColorLayer struct {
	Id    int    `xml:"id,attr"`
	Color string `xml:"color,attr"`
}

type Animation struct {
	Id           int              `xml:"id,attr"`
	TransitionTo *int             `xml:"transitionTo,attr"`
	Layers       []AnimationLayer `xml:"animationLayer"`
}

type AnimationLayer struct {
	Id             int             `xml:"id,attr"`
	LoopCount      int             `xml:"loopCount,attr"`
	FrameRepeat    int             `xml:"frameRepeat,attr"`
	Random         int             `xml:"random,attr"`
	FrameSequences []FrameSequence `xml:"frameSequence"`
}

func (layer *AnimationLayer) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	shim := struct {
		Id             int             `xml:"id,attr"`
		LoopCount      string          `xml:"loopCount,attr"`
		FrameRepeat    float64         `xml:"frameRepeat,attr"`
		Random         int             `xml:"random,attr"`
		FrameSequences []FrameSequence `xml:"frameSequence"`
	}{}

	err := d.DecodeElement(&shim, &start)
	if err != nil {
		return err
	}

	// TODO: investigate
	// -----
	// wf_act_furni_to_furni
	// frameRepeat="1.5"
	// -----
	// arcade_c23_cyberpunk
	// <frame id="NaN"/>
	// -----
	// hween12_duck
	// <animationLayer id="1" loopCount="+">

	loopCount, _ := strconv.Atoi(shim.LoopCount)
	*layer = AnimationLayer{
		Id:             shim.Id,
		LoopCount:      loopCount,
		FrameRepeat:    int(shim.FrameRepeat),
		Random:         shim.Random,
		FrameSequences: shim.FrameSequences,
	}
	return nil
}

type FrameSequence struct {
	Frames []AnimationFrame `xml:"frame"`
}

type AnimationFrame struct {
	Id int `xml:"id,attr"`
}

func (frame *AnimationFrame) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	shim := struct {
		Id string `xml:"id,attr"`
	}{}

	err := d.DecodeElement(&shim, &start)
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(shim.Id)
	*frame = AnimationFrame{Id: id}
	return nil
}
