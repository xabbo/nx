package web

type Profile struct {
	User    User     `json:"user"`
	Badges  []Badge  `json:"badges"`
	Groups  []Group  `json:"groups"`
	Rooms   []Room   `json:"rooms"`
	Friends []Friend `json:"friends"`
}
