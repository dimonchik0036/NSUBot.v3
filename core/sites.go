package core

import (
	"github.com/dimonchik0036/nsu-bot/news"
	"log"
	"sync"
	"time"
)

type Sites struct {
	Mux   sync.RWMutex     `json:"-"`
	Sites map[string]*Site `json:"sites"`
}

func (s *Sites) AddSite(site *Site) {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	if site == nil || site.Site == nil {
		log.Println("WTF?! Site is a nil pointer")
		return
	}

	s.Sites[site.Site.URL] = site
}

func (s *Sites) DelSite(href string) {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	delete(s.Sites, href)
}

func (s *Sites) ChangeSub(href string, user *User) {
	s.Mux.RLock()
	defer s.Mux.RUnlock()
	site := s.Sites[href]
	if site == nil {
		log.Printf("%s wtf?!?! href %s not found", user.String(), href)
		return
	}

	if site.Users.User(user.Platform, user.ID) == nil {
		site.Users.SetUser(user.Platform, user)
	} else {
		site.Users.DelUser(user.Platform, user.ID)
	}
}

func (s *Sites) Sub(href string, user *User) {
	s.Mux.RLock()
	defer s.Mux.RUnlock()
	site := s.Sites[href]
	if site == nil {
		log.Printf("%s wtf?!?! href %s not found", user.String(), href)
		return
	}

	site.Users.SetUser(user.Platform, user)
}

func (s *Sites) Unsub(href string, user *User) {
	s.Mux.RLock()
	defer s.Mux.RUnlock()
	site := s.Sites[href]
	if site == nil {
		log.Printf("%s wtf?!?! href %s not found", user.String(), href)
		return
	}

	site.Users.DelUser(user.Platform, user.ID)
}

func (s *Sites) CheckUser(href string, user *User) bool {
	s.Mux.RLock()
	defer s.Mux.RUnlock()
	site := s.Sites[href]
	if site == nil {
		log.Printf("%s wtf?!?! href %s not found", user.String(), href)
		return false
	}

	if site.Users.User(user.Platform, user.ID) != nil {
		return true
	}

	return false
}

func (s *Sites) Update(handler func(*Users, []news.News, string)) {
	for _, site := range s.Sites {
		news, err := site.Site.Update(5)
		if err != nil {
			log.Printf("%s error: %s", site.Site.Title, err.Error())
			continue
		}

		if len(news) == 0 {
			continue
		}

		go handler(&site.Users, news, site.Site.Title)
		time.Sleep(250 * time.Millisecond)
	}
}

type Site struct {
	Site  *news.Site `json:"site"`
	Users Users      `json:"users"`
}

func NewSites() (sites Sites) {
	sites.Sites = map[string]*Site{}
	s := news.GetAllSites()
	for _, site := range s {
		sites.Sites[site.URL] = &Site{Site: site, Users: Users{}}
	}
	return
}
