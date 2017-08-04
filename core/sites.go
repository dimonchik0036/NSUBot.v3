package core

import "github.com/dimonchik0036/nsu-bot/news"

type Sites struct {
	Sites []Site `json:"sites"`
}

type Site struct {
	Site  *news.Site `json:"site"`
	Users Users      `json:"users"`
}

func NewSites() (sites Sites) {
	s := news.GetAllSites()
	for _, site := range s {
		sites.Sites = append(sites.Sites, Site{Site: site, Users: Users{}})
	}
	return
}
