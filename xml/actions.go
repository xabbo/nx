package xml

type ActionContainer struct {
	Actions []Action
}

type Action struct {
	Id                  string `xml:"id,attr"`
	State               string `xml:"lay,attr"`
	Precedence          int    `xml:"precedence,attr"`
	Main                bool   `xml:"main,attr"`
	IsDefault           bool   `xml:"isdefault,attr"`
	GeometryType        string `xml:"geometrytype,attr"`
	ActivePartSet       string `xml:"activepartset,attr"`
	AssetPartDefinition string `xml:"assetpartdefinition,attr"`
	Prevents            string `xml:"prevents,attr"`
	Animation           bool   `xml:"animation,attr"`
	PreventHeadTurn     bool   `xml:"preventheadturn,attr"`
	StartFromFrameZero  bool   `xml:"startfromframezero,attr"`
	Types               []ActionType
	Params              []ActionParam
}

type ActionType struct {
	Id              int    `xml:"id,attr"`
	Animated        bool   `xml:"animation,attr"`
	Prevents        string `xml:"prevents,attr"`
	PreventHeadTurn bool   `xml:"preventheadturn,attr"`
}

type ActionParam struct {
	Id    string `xml:"id,attr"`
	Value string `xml:"value,attr"`
}
