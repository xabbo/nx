package web

type Friend struct {
	UniqueId     string `json:"uniqueId"`
	Name         string `json:"name"`
	Motto        string `json:"motto"`
	FigureString string `json:"figureString"`
	Online       bool   `json:"online"`
}
