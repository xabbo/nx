package xml

type IndexData struct {
	Type          string `xml:"type,attr"`
	Visualization string `xml:"visualization,attr"`
	Logic         string `xml:"logic"`
}
