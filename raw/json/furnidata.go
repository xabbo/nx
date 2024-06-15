package json

type FurniData struct {
	FloorItems FurniInfos `json:"roomitemtypes"`
	WallItems  FurniInfos `json:"wallitemtypes"`
}

type FurniInfos struct {
	Infos []FurniInfo `json:"furnitype"`
}

type FurniInfo struct {
	Id              int        `json:"id"`
	Identifier      string     `json:"classname"`
	Revision        int        `json:"revision"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	Category        string     `json:"category"`
	Environment     string     `json:"environment"`
	Line            string     `json:"furniline"`
	DefaultDir      int        `json:"defaultdir"`
	XDim            int        `json:"xdim"`
	YDim            int        `json:"ydim"`
	PartColors      PartColors `json:"partcolors"`
	AdUrl           string     `json:"adurl"` // Obsolete.
	OfferId         int        `json:"offerid"`
	Buyout          bool       `json:"buyout"`
	RentOfferId     int        `json:"rentofferid"` // Obsolete.
	RentBuyout      bool       `json:"rentbuyout"`  // Obsolete.
	BC              bool       `json:"bc"`
	ExcludedDynamic bool       `json:"excludeddynamic"`
	CustomParams    string     `json:"customparams"`
	SpecialType     int        `json:"specialtype"`
	CanStandOn      bool       `json:"canstandon"`
	CanSitOn        bool       `json:"cansiton"`
	CanLayOn        bool       `json:"canlayon"`
	Rare            bool       `json:"rare"` // Obsolete.
}

type PartColors struct {
	Colors []string `json:"color"`
}
