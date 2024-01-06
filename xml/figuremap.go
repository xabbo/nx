package xml

type FigureMap struct {
	Libraries []FigureMapLib `xml:"lib"`
}

type FigureMapLib struct {
	Id       string          `xml:"id,attr"`
	Revision int             `xml:"revision,attr"`
	Parts    []FigureMapPart `xml:"part"`
}

type FigureMapPart struct {
	Id   string `xml:"id,attr"`
	Type string `xml:"type,attr"`
}
