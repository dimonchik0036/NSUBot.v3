package news

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
)

type News struct {
	ID         int64
	URL        string
	Title      string
	Decryption string
	Date       int64
}

type Sites struct {
	Sites []*Site
}

func NewSites() Sites {
	return Sites{
		Sites: append(FitNews(), FitChairs()...),
	}
}

type Site struct {
	Mux           sync.Mutex
	Title         string
	OptionalTitle string
	URL           string
	NewsPage      func(href string, count int) ([]News, error)
	LastNews      News
}

func (s *Site) Update(countCheck int) (newNews []News, err error) {
	news, err := s.NewsPage(s.URL, countCheck)
	if err != nil || len(news) == 0 {
		return newNews, err
	}

	s.Mux.Lock()
	defer s.Mux.Unlock()

	if s.LastNews.ID == 0 {
		s.LastNews = news[0]
		return reversNews(news), nil
	}

	for i := range news {
		if (s.LastNews.ID >= news[i].ID) && news[i].ID != -1 || s.LastNews.ID == -1 && news[i].ID == -1 && (news[i].URL == s.LastNews.URL) {
			break
		}

		newNews = append(newNews, news[i])
	}

	s.LastNews = news[0]

	return reversNews(newNews), nil
}

func getNewsPage(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	if res.StatusCode != http.StatusOK {
		return []byte{}, errors.New("Error status: " + res.Status)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	rg, err := regexp.Compile("[\n\t]")
	if err != nil {
		return []byte{}, err
	}

	return rg.ReplaceAll(body, []byte("")), nil
}

func reversNews(news []News) (result []News) {
	for i := range news {
		result = append(result, news[len(news)-i-1])
	}

	return
}
