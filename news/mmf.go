package news

import (
	"html"
	"regexp"
)

const (
	MmfHref       = "http://mmf.nsu.ru"
	MmfTimeLayout = "02.01.2006"
)

func MmfNews() []*Site {
	return []*Site{
		&Site{
			Title:    "Новости",
			URL:      "/news/index",
			NewsPage: Mmf,
		},
		&Site{
			Title:    "Объявления",
			URL:      "/advert/index",
			NewsPage: Mmf,
		},
		&Site{
			Title:    "Объявления студентам",
			URL:      "/students/advert",
			NewsPage: Mmf,
		},
	}
}

func Mmf(href string, count int) (news []News, err error) {
	body, err := getNewsPage(MmfHref + href)
	if err != nil {
		return []News{}, err
	}

	rg, err := regexp.Compile("<div class=\"views-field views-field-title\">.*?</div>")
	if err != nil {
		return []News{}, err
	}

	dates := dateProcessing(body, count, "<span class=\"date-display-single\">", "</span>", MmfTimeLayout)

	for i, b := range rg.FindAll(body, count) {
		for _, v := range hrefProcessing(b, 1) {
			news = append(news, News{
				ID:    idScan(string(v[0])),
				Title: html.UnescapeString(string(v[1])),
				URL:   MmfHref + string(v[0]),
				Date:  dates[i],
			})
		}
	}

	return
}
