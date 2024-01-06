package xml

type LogicData struct {
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
	Id int `xml:"id"`
}

type ParticleSystem struct {
	Size int `xml:"size,attr"`
}
