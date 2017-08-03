package main

import (
	"fmt"
	"github.com/dimonchik0036/nsu-bot/gconfig"
	"github.com/dimonchik0036/nsu-bot/news"
	"github.com/dimonchik0036/nsu-bot/nsuschedule"
	"github.com/dimonchik0036/nsu-bot/nsuweather"
	"github.com/dimonchik0036/nsu-bot/telegram-bot"
	"github.com/dimonchik0036/nsu-bot/vk-bot"
	"time"
)

func main() {
	println("start")
	sites := news.NewSites()

	for _, s := range sites.Sites {
		n, err := s.Update(3)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(s.Title)
		for i := range n {
			//fmt.Println(n[i])
			fmt.Println(time.Unix(n[i].Date, 0).Format(news.TimeLayout), n[i].Title)
		}
		fmt.Println()
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
