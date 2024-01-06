package web

type Group struct {
	Online         bool   `json:"online"`
	Id             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Type           string `json:"type"`
	RoomId         string `json:"roomId"`
	BadgeCode      string `json:"badgeCode"`
	PrimaryColor   string `json:"primaryColour"`
	SecondaryColor string `json:"secondaryColour"`
	IsAdmin        bool   `json:"isAdmin"`
}
