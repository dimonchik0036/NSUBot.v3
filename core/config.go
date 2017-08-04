package core

import (
	"github.com/dimonchik0036/nsu-bot/nsuschedule"
	"github.com/dimonchik0036/nsu-bot/nsuweather"
)

type Config struct {
	Schedule *nsuschedule.Schedule
	Weather  *nsuweather.Weather
	Sites    *Sites
}

func NewConfig() (config Config) {
	schedule := nsuschedule.NewSchedule()
	weather := nsuweather.NewWeather()
	sites := NewSites()
	config.Weather = &weather
	config.Schedule = &schedule
	config.Sites = &sites
	return
}
