package web

type Badge struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SelectedBadge struct {
	Badge
	Index int `json:"badgeIndex"`
}
