package tgbot

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/dimonchik0036/nsu-bot/nsuschedule"
	"github.com/dimonchik0036/nsu-bot/nsuweather"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var tgAdminID int64
var tgBot *tgbotapi.BotAPI
var tgUsers *core.Users
var tgWeather *nsuweather.Weather
var tgSites *core.Sites
var tgSchedule *nsuschedule.Schedule

func loadTgConfig() {
	data, err := ioutil.ReadFile(".tg_config")
	if err != nil {
		log.Printf("Tg config not found: %s", err.Error())
		return
	}

	tmp := struct {
		ID    int64
		Token string
	}{}

	if err := yaml.Unmarshal(data, &tmp); err != nil {
		log.Printf("Tg config: yaml throw error: %s", err.Error())
		return
	}

	tgAdminID = tmp.ID
	bot, err := tgbotapi.NewBotAPI(tmp.Token)
	if err != nil {
		log.Panicf("Bot is offline: %s", err.Error())
	}

	tgBot = bot
	if _, err := tgBot.Send(tgbotapi.NewMessage(tgAdminID, "Я запущен")); err != nil {
		log.Panicf("Tg error: %s", err.Error())
	}

	return
}

func initConfig(config *core.Config) {
	config.Mux.Lock()
	defer config.Mux.Unlock()
	tgWeather = config.Weather
	tgSites = config.Sites
	tgSchedule = config.Schedule
	tgUsers = config.Users
}
