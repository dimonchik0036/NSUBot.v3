package news

import (
	"html"
	"regexp"
	"strconv"
	"time"
)

const (
	FitHref    = "http://fit.nsu.ru"
	TimeLayout = "02-01-2006"
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
	dates := dateProcessing(body, count)

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

func dateProcessing(body []byte, count int) (dates []int64) {
	begin := "<td class=\"list-date\">"
	end := "</td>"
	rg, err := regexp.Compile(begin + ".*?" + end)
	if err != nil {
		return
	}

	for _, date := range rg.FindAll(body, count) {
		t, err := time.Parse(begin+"02-01-06"+end, string(date))
		if err != nil {
			panic(err)
		}
		dates = append(dates, t.Unix())
	}
	return
}

func idScan(url string) int64 {
	rg, err := regexp.Compile("/[\\d]*?-")
	if err != nil {
		return 0
	}

	idString := rg.FindString(url)
	if len(idString) < 3 {
		return -1
	}

	id, err := strconv.ParseInt(idString[1:len(idString)-1], 10, 64)
	if err != nil {
		return 0
	}

	return id
}
