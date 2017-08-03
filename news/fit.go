package news

import (
	"html"
	"regexp"
	"strconv"
	"time"
)

const (
	FitHref       = "http://fit.nsu.ru"
	FitTimeLayout = "02-01-06"
)

func FitNews() []*Site {
	return []*Site{
		&Site{
			Title:    "События",
			URL:      "/news/news",
			NewsPage: Fit,
		},
		&Site{
			Title:    "Объявления",
			URL:      "/news/announc",
			NewsPage: Fit,
		},
		&Site{
			Title:    "Конференции",
			URL:      "/news/konf",
			NewsPage: Fit,
		},
		&Site{
			Title:    "Конкурсы",
			URL:      "/news/conc",
			NewsPage: Fit,
		},
		&Site{
			Title:    "Вакансии",
			URL:      "/news/vac",
			NewsPage: Fit,
		},
		&Site{
			Title:    "Административные приказы",
			URL:      "/news/administrativnye-prikazy",
			NewsPage: Fit,
		},
	}
}

func FitChairs() []*Site {
	return []*Site{
		&Site{
			Title:    "Объявления кафедры систем информатики",
			URL:      "/chairs/ksi/anksi",
			NewsPage: Fit,
		},
		&Site{
			Title:    "Объявления кафедры компьютерных систем",
			URL:      "/chairs/kks/ankks",
			NewsPage: Fit,
		},
		&Site{
			Title:    "Объявления кафедры общей информатики",
			URL:      "/chairs/koi/koinews",
			NewsPage: Fit,
		},
		&Site{
			Title:    "Объявления кафедры параллельных вычислений",
			URL:      "/chairs/kpv/kpvnews",
			NewsPage: Fit,
		},
		&Site{
			Title:    "Объявления кафедры компьютерных технологий",
			URL:      "/chairs/k-kt/kktnews",
			NewsPage: Fit,
		},
	}
}

func Fit(href string, count int) (news []News, err error) {
	body, err := getNewsPage(FitHref + href + "?limit=" + strconv.Itoa(count))
	if err != nil {
		return []News{}, err
	}

	rg, err := regexp.Compile("<tbody>.*?</tbody>")
	if err != nil {
		return []News{}, err
	}

	body = rg.Find(body)
	hrefs := hrefProcessing(body, count)
	dates := dateProcessing(body, count, "<td class=\"list-date\">", "</td>", FitTimeLayout)
	for i, v := range hrefs {
		news = append(news, News{
			ID:    idScan(string(v[0])),
			Title: html.UnescapeString(string(v[1])),
			URL:   FitHref + string(v[0]),
			Date:  dates[i],
		})
	}

	return
}

func FitNavigation(href string) ([]*Site, error) {
	body, err := getNewsPage(href)
	if err != nil {
		return []*Site{}, err
	}

	rg, err := regexp.Compile("<nav class=\"leftmenu\">.*?</nav>")
	if err != nil {
		return []*Site{}, err
	}

	return fitNavigationProcessing(rg.Find(body), Fit), nil
}

func fitNavigationProcessing(body []byte, f func(string, int) ([]News, error)) (sites []*Site) {
	for _, b := range hrefProcessing(body, -1) {
		sites = append(sites, &Site{
			URL:      string(b[0]),
			Title:    string(b[1]),
			NewsPage: f,
		})
	}

	return
}
