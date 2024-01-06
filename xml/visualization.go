package xml

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
	Color      []Color     `xml:"colors>color"`
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
	Id          int              `xml:"id,attr"`
	LoopCount   int              `xml:"loopCount"`
	FrameRepeat int              `xml:"frameRepeat"`
	Frames      []AnimationFrame `xml:"frameSequence>frame"`
}

type AnimationFrame struct {
	Id int `xml:"id,attr"`
}
