package main

import (
	"github.com/dimonchik0036/nsu-bot/gconfig"
	"github.com/dimonchik0036/nsu-bot/schedule"
	"github.com/dimonchik0036/nsu-bot/telegram-bot"
	"github.com/dimonchik0036/nsu-bot/vk-bot"
	"github.com/dimonchik0036/nsu-bot/nsuweather"
)

func main() {
	weather, _ := nsuweather.GetWeather()
	config := gconfig.Config{
		Schedule: &schedule.Schedule{},
		Weather:  &weather,
	}

	tgbot.Processing(config)
	vkbot.Processing(config)
}
