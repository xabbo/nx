package xml

// index.xml

type Index struct {
	Type          string `xml:"type,attr"`
	Visualization string `xml:"visualization,attr"`
	Logic         string `xml:"logic"`
}

// logic.xml

type Logic struct {
	Type  string `xml:"type,attr"`
	Model Model  `xml:"model"`
}

type Model struct {
	Dimensions      Dimensions       `xml:"dimensions"`
	Directions      []Direction      `xml:"directions>direction"`
	ParticleSystems []ParticleSystem `xml:"particlesystems>particlesystem"`
}

type Dimensions struct {
	X int     `xml:"x,attr"`
	Y int     `xml:"y,attr"`
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
	Id     int              `xml:"id,attr"`
	Layers []AnimationLayer `xml:"animationLayer"`
}

type AnimationLayer struct {
	Id             int             `xml:"id,attr"`
	LoopCount      int             `xml:"loopCount,attr"`
	FrameRepeat    int             `xml:"frameRepeat,attr"`
	Random         bool            `xml:"random,attr"`
	FrameSequences []FrameSequence `xml:"frameSequence"`
}

type FrameSequence struct {
	Frames []AnimationFrame `xml:"frame"`
}

type AnimationFrame struct {
	Id int `xml:"id,attr"`
}
