package xml

type FigureData struct {
	Palettes []FigurePalette  `xml:"colors>palette"`
	Sets     []FigurePartSets `xml:"sets>settype"`
}

type FigurePalette struct {
	Id     int           `xml:"id,attr"`
	Colors []FigureColor `xml:"color"`
}

type FigureColor struct {
	Id         int    `xml:"id,attr"`
	Index      int    `xml:"index,attr"`
	Club       int    `xml:"club,attr"`
	Selectable bool   `xml:"selectable,attr"`
	Value      string `xml:",chardata"`
}

type FigurePartSets struct {
	Type      string          `xml:"type,attr"`
	PaletteId int             `xml:"paletteid,attr"`
	MandM0    bool            `xml:"mand_m_0,attr"`
	MandF0    bool            `xml:"mand_f_0,attr"`
	MandM1    bool            `xml:"mand_m_1,attr"`
	MandF1    bool            `xml:"mand_f_1,attr"`
	Sets      []FigurePartSet `xml:"set"`
}

type FigurePartSet struct {
	Id            int           `xml:"id,attr"`
	Gender        string        `xml:"gender,attr"`
	Club          int           `xml:"club,attr"`
	Colorable     bool          `xml:"colorable,attr"`
	Selectable    bool          `xml:"selectable,attr"`
	Preselectable bool          `xml:"preselectable,attr"`
	Sellable      bool          `xml:"sellable,attr"`
	Parts         []FigurePart  `xml:"part"`
	HiddenLayers  []FigureLayer `xml:"hiddenlayers>layer"`
}

type FigurePart struct {
	Id         int    `xml:"id,attr"`
	Type       string `xml:"type,attr"`
	Colorable  bool   `xml:"colorable,attr"`
	Index      int    `xml:"index,attr"`
	ColorIndex int    `xml:"colorindex,attr"`
}

type FigureLayer struct {
	PartType string `xml:"parttype,attr"`
}
