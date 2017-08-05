package vkbot

import "github.com/dimonchik0036/nsu-bot/core"

var Commands map[string]core.Handler

func initCommands() {
	Commands = map[string]core.Handler{
		"weather":     weatherHandler(),
		"w":           weatherHandler(),
		"погода":      weatherHandler(),
		"температура": weatherHandler(),
	}
}

func weatherHandler() core.Handler {
	return core.Handler{
		PermissionLevel: 0,
		Handler:         WeatherCommand,
	}
}

func WeatherCommand(user *core.User, command core.Command) {
	Bot.SendMessage(user.ID, Weather.ShowWeather()+"\n"+
		"Время последнего обновления: "+Weather.ShowTime())
}
