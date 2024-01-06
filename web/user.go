package web

import (
	"encoding/json"
	"time"
)

type User struct {
	UniqueId                    string          `json:"uniqueId"`
	Name                        string          `json:"name"`
	FigureString                string          `json:"figureString"`
	Motto                       string          `json:"motto"`
	Online                      bool            `json:"online"`
	LastAccessTime              *time.Time      `json:"lastAccessTime"`
	MemberSince                 time.Time       `json:"memberSince"`
	ProfileVisible              bool            `json:"profileVisible"`
	CurrentLevel                int             `json:"currentLevel"`
	CurrentLevelCompletePercent int             `json:"currentLevelCompeltePercent"`
	TotalExperience             int             `json:"totalExperience"`
	StarGemCount                int             `json:"starGemCount"`
	SelectedBadges              []SelectedBadge `json:"selectedBadges"`
}

func (u *User) UnmarshalJSON(data []byte) (err error) {
	type Alias User
	shim := &struct {
		Alias
		LastAccessTime string `json:"lastAccessTime"`
		MemberSince    string `json:"memberSince"`
	}{}
	err = json.Unmarshal(data, &shim)
	if err != nil {
		return err
	}
	*u = User(shim.Alias)
	var t time.Time
	// 2024-01-06T13:41:24.000+0000
	tf := "2006-01-02T15:04:05.000-0700"
	if shim.LastAccessTime != "" {
		t, err = time.Parse(tf, shim.LastAccessTime)
		if err != nil {
			return
		}
		u.LastAccessTime = &t
	}
	u.MemberSince, err = time.Parse(tf, shim.MemberSince)
	return
}
