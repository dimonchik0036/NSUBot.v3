package nsuschedule

import (
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Schedule struct {
	Mux      sync.RWMutex        `json:"-"`
	Schedule map[string][]string `json:"schedule"`
}

func (s *Schedule) GetGroup(group string) ([]string, bool) {
	s.Mux.RLock()
	defer s.Mux.RUnlock()
	str, ok := s.Schedule[group]
	return str, ok
}

func (s *Schedule) GetDay(group string, day int) (string, bool) {
	g, ok := s.GetGroup(group)
	if !ok {
		return "", ok
	}
	if day < 0 {
		return "", false
	}
	return g[day%7], true
}

type Day struct {
	Lessons []*Lesson
}

const (
	scheduleUrl = "http://table.nsu.ru/json"
)

func getSchedule() (University, error) {
	code, body, err := fasthttp.Get(nil, scheduleUrl)
	if err != nil {
		return University{}, err
	}

	if code != fasthttp.StatusOK {
		return University{}, errors.New("Bad code: " + strconv.Itoa(code))
	}

	var u University
	if err := json.Unmarshal(body, &u); err != nil {
		return University{}, err
	}

	return u, nil
}

func NewSchedule() Schedule {
	u, err := getSchedule()
	if err != nil {
		return Schedule{}
	}

	return Schedule{Schedule: u.Schedule()}
}

func (s *Schedule) Update() {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	u, err := getSchedule()
	if err != nil {
		log.Print("Bad schedule update: ", err)
		time.Sleep(5 * time.Minute)
		u, err = getSchedule()
		if err != nil {
			log.Print("Very bad schedule update: ", err)
			return
		}
	}

	s.Schedule = u.Schedule()
}

type University struct {
	Name      string       `json:"name"`
	Abbr      string       `json:"abbr"`
	Faculties []*Faculties `json:"faculties"`
}

func (u *University) Schedule() map[string][]string {
	var s map[string][]string
	s = make(map[string][]string)
	for _, f := range u.Faculties {
		for _, g := range f.Groups {
			s[g.Name] = g.GetLessons()
		}
	}
	return s
}

type Faculties struct {
	Name   string   `json:"name"`
	Groups []*Group `json:"groups"`
}

type Group struct {
	Name    string    `json:"name"`
	Lessons []*Lesson `json:"lessons"`
}

func (g *Group) GetLessons() []string {
	var days [6]map[string][]*Lesson

	for _, l := range g.Lessons {
		if l.Date.Weekday > 6 || l.Date.Weekday < 1 {
			continue
		}

		if days[l.Date.Weekday-1] == nil {
			days[l.Date.Weekday-1] = make(map[string][]*Lesson)
		}
		if strings.HasPrefix(l.Time.Start, "9:") {
			l.Time.Start = "0" + l.Time.Start
		}
		days[l.Date.Weekday-1][l.Time.Start+"-"+l.Time.End] = append(days[l.Date.Weekday-1][l.Time.Start+"-"+l.Time.End], l)
	}

	var result []string
	for _, day := range days {
		var keys []string
		for k := range day {
			keys = append(keys, k)
		}

		if len(keys) == 0 {
			result = append(result, "Нет пар")
			continue
		}

		sort.Strings(keys)
		var s string
		for _, k := range keys {
			for _, l := range day[k] {
				s += l.String() + "\n"
			}
		}
		result = append(result, s)
	}
	result = append(result, "Нет пар")
	return result
}

type Lesson struct {
	Subject string `json:"subject"`
	Type    string `json:"type"`
	Time    struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"time"`
	Date struct {
		Start   string `json:"start"`
		End     string `json:"end"`
		Weekday int    `json:"weekday"`
		Week    int    `json:"week"`
	} `json:"date"`
	Audiences []struct {
		Name string `json:"name"`
	} `json:"audiences"`
	Teachers []struct {
		Name string `json:"name"`
	} `json:"teachers"`
}

func (l *Lesson) String() string {
	return func() string {
		switch l.Date.Week {
		case 0:
			return ""
		case 1:
			return "Ч "
		case 2:
			return "Н "
		default:
			return ""
		}
	}() + l.Time.Start + "-" + l.Time.End + ": " + l.Type + ", " + l.Subject + ", " + help(l.Audiences) + ", " + help(l.Teachers)
}

func help(s []struct {
	Name string `json:"name"`
}) string {
	var res []string
	for _, n := range s {
		res = append(res, n.Name)
	}
	return strings.Join(res, ", ")
}
