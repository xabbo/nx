package web

import "time"

type Room struct {
	Id              int       `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	CreationTime    time.Time `json:"creationTime"`
	HabboGroupId    string    `json:"habboGroupId"`
	Tags            []string  `json:"tags"`
	MaximumVisitors int       `json:"maximumVisitors"`
	ShowOwnerName   bool      `json:"showOwnerName"`
	OwnerNamne      string    `json:"ownerName"`
	OwnerUniqueId   string    `json:"ownerUniqueId"`
	Categories      []string  `json:"categories"`
	ThumbnailUrl    string    `json:"thumbnailUrl"`
	ImageUrl        string    `json:"imageUrl"`
	Rating          int       `json:"rating"`
	UniqueId        string    `json:"uniqueId"`
}
