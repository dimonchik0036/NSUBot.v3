package gconfig

import (
	"github.com/dimonchik0036/nsu-bot/news"
	"github.com/dimonchik0036/nsu-bot/nsuschedule"
	"github.com/dimonchik0036/nsu-bot/nsuweather"
)

type Config struct {
	Schedule *nsuschedule.Schedule
	Weather  *nsuweather.Weather
	Sites    *news.Sites
}
