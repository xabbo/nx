package xml

type AssetContainer struct {
	Assets []Asset `xml:"asset"`
}

type Asset struct {
	Name     string  `xml:"name,attr"`
	MimeType string  `xml:"mimeType,attr"`
	X        int     `xml:"x,attr"`
	Y        int     `xml:"y,attr"`
	FlipH    bool    `xml:"flipH,attr"`
	FlipV    bool    `xml:"flipV,attr"`
	Source   string  `xml:"source,attr"`
	Params   []Param `xml:"param"`
}

type Param struct {
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr"`
}
