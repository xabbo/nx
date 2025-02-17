package nitro

import (
	"encoding/json"
	"strconv"
)

type Furni struct {
	Name              string           `json:"name"`
	LogicType         string           `json:"logicType"`
	VisualizationType string           `json:"visualizationType"`
	Assets            map[string]Asset `json:"assets"`
	Logic             Logic            `json:"logic"`
	Visualizations    []Visualization  `json:"visualizations"`
	Spritesheet       Spritesheet      `json:"spritesheet"`
}

type Asset struct {
	Source string  `json:"source"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	FlipH  bool    `json:"flipH"`
	FlipV  bool    `json:"flipV"`
}

type Logic struct {
	Model           Model            `json:"model"`
	ParticleSystems []ParticleSystem `json:"particleSystems"`
}

type Model struct {
	Dimensions Dimensions `json:"dimensions"`
	Directions []int      `json:"directions"`
}

type Dimensions struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type ParticleSystem struct {
	Size int `json:"size"`
}

type Visualization struct {
	Angle      int               `json:"angle"`
	LayerCount int               `json:"layerCount"`
	Size       int               `json:"size"`
	Layers     map[int]Layer     `json:"layers"`
	Directions map[int]Direction `json:"directions"`
	Colors     map[int]Color     `json:"colors"`
	Animations map[int]Animation `json:"animations"`
}

type Layer struct {
	Z           float64 `json:"z"`
	Ink         string  `json:"ink"`
	Alpha       int     `json:"alpha"`
	Color       int     `json:"color"`
	IgnoreMouse bool    `json:"ignoreMouse"`
}

type Direction struct {
	Layers map[int]Layer `json:"layers"`
}

type Color struct {
	Layers map[int]Layer `json:"layers"`
}

type Animation struct {
	Layers       map[int]AnimationLayer `json:"layers"`
	TransitionTo *int
}

type AnimationLayer struct {
	LoopCount      float64               `json:"loopCount"`
	FrameRepeat    float64               `json:"frameRepeat"`
	Random         float64               `json:"random"`
	FrameSequences map[int]FrameSequence `json:"frameSequences"`
}

type FrameSequence struct {
	Frames map[int]AnimationFrame `json:"frames"`
}

type AnimationFrame struct {
	Id int `json:"id"`
}

type Spritesheet struct {
	Frames map[string]SpriteFrame `json:"frames"`
	Meta   Meta                   `json:"meta"`
}

type SpriteFrame struct {
	Frame            Size  `json:"frame"`
	Rotated          bool  `json:"rotated"`
	Trimmed          bool  `json:"trimmed"`
	SpriteSourceSize Size  `json:"spriteSourceSize"`
	SourceSize       Size  `json:"sourceSize"`
	Pivot            Pivot `json:"pivot"`
}

type Size struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type Pivot struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Meta struct {
	Image  string  `json:"image"`
	Format string  `json:"format"`
	Size   Size    `json:"size"`
	Scale  float64 `json:"scale"`
}

func (m *Meta) UnmarshalJSON(data []byte) (err error) {
	shim := struct {
		Image  string `json:"image"`
		Format string `json:"format"`
		Size   Size   `json:"size"`
		Scale  any    `json:"scale"`
	}{}

	err = json.Unmarshal(data, &shim)
	if err == nil {
		*m = Meta{
			Image:  shim.Image,
			Format: shim.Format,
			Size:   shim.Size,
		}
		switch scale := shim.Scale.(type) {
		case float64:
			m.Scale = scale
		case string:
			scalef64, err := strconv.ParseFloat(scale, 64)
			if err != nil {
				return err
			}
			m.Scale = scalef64
		}
	}
	return
}
