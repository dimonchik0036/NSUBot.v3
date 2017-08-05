package core

import (
	"github.com/dimonchik0036/nsu-bot/news"
	"log"
)

type Sites struct {
	Sites []Site `json:"sites"`
}

func (s *Sites) Update(handler func(*Users, []news.News)) {
	for _, site := range s.Sites {
		news, err := site.Site.Update(5)
		if err != nil {
			log.Printf("%s error: %s", site.Site.Title, err.Error())
			continue
		}

		if len(news) == 0 {
			continue
		}

		go handler(&site.Users, news)
	}
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
