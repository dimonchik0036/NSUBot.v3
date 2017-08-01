package main

import (
	"fmt"
	"github.com/dimonchik0036/nsu-bot/gconfig"
	"github.com/dimonchik0036/nsu-bot/news"
	"github.com/dimonchik0036/nsu-bot/nsuschedule"
	"github.com/dimonchik0036/nsu-bot/nsuweather"
	"github.com/dimonchik0036/nsu-bot/telegram-bot"
	"github.com/dimonchik0036/nsu-bot/vk-bot"
)

func main() {
	println("start")

	sites := news.NewSites()

	for _, s := range sites.Sites {
		fmt.Println(s.Update(25))
	}

	weather := nsuweather.NewWeather()
	schedule := nsuschedule.NewSchedule()

	config := gconfig.Config{
		Schedule: &schedule,
		Weather:  &weather,
		Sites:    &sites,
	}

	tgbot.Processing(config)
	vkbot.Processing(config)
}
