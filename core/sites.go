package core

import (
	"fmt"
	"github.com/dimonchik0036/nsu-bot/news"
)

type Sites struct {
	Sites []Site `json:"sites"`
}

func (s *Sites) Update() {
	for _, site := range s.Sites {
		n, err := site.Site.Update(2)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(site.Site.Title)
		for i := range n {
			fmt.Println(n[i])
			//fmt.Println(time.Unix(n[i].Date, 0).Format(news.TimeLayout), "\n"+n[i].Title+" "+n[i].URL+"\n"+n[i].Decryption)
		}
		fmt.Println()
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
