package news

import (
	"errors"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type News struct {
	ID         int64  `json:"id"`
	URL        string `json:"url"`
	Title      string `json:"title"`
	Decryption string `json:"decryption"`
	Date       int64  `json:"date"`
}

type Sites struct {
	Sites map[string]*Site `json:"sites"`
}

const (
	TimeLayout = "02.01.2006"
)

func GetAllSites() (sites Sites) {
	sites.Sites = make(map[string]*Site)

	n := NsuNews()
	n = append(n, PhilosNews()...)
	n = append(n, FpNews()...)
	n = append(n, MmfNews()...)
	n = append(n, FitNews()...)
	n = append(n, FitChairs()...)

	for _, s := range n {
		sites.Sites[s.URL] = s
	}

	return
}

type Site struct {
	Mux           sync.Mutex                                   `json:"-"`
	Title         string                                       `json:"title"`
	OptionalTitle string                                       `json:"optional_title"`
	URL           string                                       `json:"url"`
	NewsFunc      func(href string, count int) ([]News, error) `json:"-"`
	NewsFuncName  string                                       `json:"news_func_name"`
	LastNews      News                                         `json:"last_news"`
}

func (s *Site) Update(countCheck int) (newNews []News, err error) {
	news, err := s.NewsFunc(s.URL, countCheck)
	if err != nil || len(news) == 0 {
		return newNews, err
	}

	s.Mux.Lock()
	defer s.Mux.Unlock()

	for i := range news {
		if (s.LastNews.ID > news[i].ID) && news[i].ID != 0 || (news[i].URL == s.LastNews.URL) && (news[i].Date == s.LastNews.Date) && (news[i].Title == s.LastNews.Title) {
			break
		}

		newNews = append(newNews, news[i])
	}

	s.LastNews = news[0]

	return reversNews(newNews), nil
}

func (s *Site) InitFunc() {
	switch s.NewsFuncName {
	case NsuFacFuncName:
		s.NewsFunc = NsuFac
	case NsuFuncName:
		s.NewsFunc = Nsu
	case FitFuncName:
		s.NewsFunc = Fit
	case PhilosFuncName:
		s.NewsFunc = Philos
	case MmfFuncName:
		s.NewsFunc = Mmf
	case FpFuncName:
		s.NewsFunc = Fp
	default:
		panic("WTF?!")
	}
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

	utf8, err := charset.NewReader(res.Body, res.Header.Get("Content-Type"))
	if err != nil {
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(utf8)
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

	rgTitle, err := regexp.Compile(">")
	if err != nil {
		return
	}

	rgHref, err := regexp.Compile("\"")
	if err != nil {
		return
	}

	for _, href := range rg.FindAll(body, count) {
		href = href[9 : len(href)-4]
		titleInd := rgTitle.FindIndex(href)
		hrefInd := rgHref.FindIndex(href)
		result = append(result, [][]byte{href[:hrefInd[0]], href[titleInd[1]:]})
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
