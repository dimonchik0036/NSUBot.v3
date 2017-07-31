package main

import (
	"github.com/dimonchik0036/nsu-bot/gconfig"
	"github.com/dimonchik0036/nsu-bot/nsuschedule"
	"github.com/dimonchik0036/nsu-bot/nsuweather"
	"github.com/dimonchik0036/nsu-bot/telegram-bot"
	"github.com/dimonchik0036/nsu-bot/vk-bot"
)

func main() {
	weather, _ := nsuweather.GetWeather()
	schedule, _ := nsuschedule.GetAllSchedule()

	config := gconfig.Config{
		Schedule: &schedule,
		Weather:  &weather,
	}

	tgbot.Processing(config)
	vkbot.Processing(config)
}
