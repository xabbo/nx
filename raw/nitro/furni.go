package nitro

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
	Source string `json:"source"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	FlipH  bool   `json:"flipH"`
	FlipV  bool   `json:"flipV"`
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
	X int     `json:"x"`
	Y int     `json:"y"`
	Z float64 `json:"z"`
}

type ParticleSystem struct {
	Size int `json:"size"`
}

type Visualization struct {
	Angle      int                  `json:"angle"`
	LayerCount int                  `json:"layerCount"`
	Size       int                  `json:"size"`
	Layers     map[string]Layer     `json:"layers"`
	Directions map[string]struct{}  `json:"directions"`
	Animations map[string]Animation `json:"animations"`
}

type Layer struct {
	Z     int    `json:"z"`
	Ink   string `json:"ink"`
	Alpha int    `json:"alpha"`
}

type Animation struct {
	Layers map[string]AnimationLayer `json:"layers"`
}

type AnimationLayer struct {
	FrameRepeat    int                      `json:"frameRepeat"`
	Random         int                      `json:"random"`
	FrameSequences map[string]FrameSequence `json:"frameSequences"`
}

type FrameSequence struct {
	Frames map[string]AnimationFrame `json:"frames"`
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
	Image  string `json:"image"`
	Format string `json:"format"`
	Size   Size   `json:"size"`
	Scale  int    `json:"scale"`
}
