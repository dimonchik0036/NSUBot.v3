package news

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
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

const (
	TimeLayout = "02.01.2006"
)

func NewSites() (sites Sites) {
	sites.Sites = append(sites.Sites, MmfNews()...)
	sites.Sites = append(sites.Sites, FitNews()...)
	sites.Sites = append(sites.Sites, FitChairs()...)
	return
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

func hrefProcessing(body []byte, count int) (result [][][]byte) {
	rg, err := regexp.Compile("<a.*?>.*?</a>")
	if err != nil {
		return
	}

	rgHref, err := regexp.Compile("\" ?>")
	if err != nil {
		return
	}

	for _, href := range rg.FindAll(body, count) {
		href = href[9 : len(href)-4]
		begInd := rgHref.FindIndex(href)
		result = append(result, [][]byte{href[:begInd[0]], href[begInd[1]:]})
	}

	return
}

func idScan(url string) int64 {
	rg, err := regexp.Compile("[\\d]+")
	if err != nil {
		return 0
	}

	id, err := strconv.ParseInt(rg.FindString(url), 10, 64)
	if err != nil {
		return 0
	}

	return id
}

func dateProcessing(body []byte, count int, begin string, end string, layout string) (dates []int64) {
	rg, err := regexp.Compile(begin + ".*?" + end)
	if err != nil {
		return
	}

	for _, date := range rg.FindAll(body, count) {
		t, err := time.Parse(begin+layout+end, string(date))
		if err != nil {
			dates = append(dates, 0)
		}
		dates = append(dates, t.Unix())
	}
	return
}
