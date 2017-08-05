package vkbot

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/dimonchik0036/nsu-bot/nsuschedule"
	"github.com/dimonchik0036/nsu-bot/nsuweather"
	"github.com/dimonchik0036/vk-api"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var AdminID int64
var Bot *BotSt
var Users *core.Users
var Weather *nsuweather.Weather
var Sites *core.Sites
var Schedule *nsuschedule.Schedule

func loadVkConfig() {
	data, err := ioutil.ReadFile(".vk_config")
	if err != nil {
		log.Panicf("Vk config not found: %s", err.Error())
		return
	}

	tmp := struct {
		ID    int64
		Token string
	}{}

	if err := yaml.Unmarshal(data, &tmp); err != nil {
		log.Panicf("Vk config: yaml throw error: %s", err.Error())
		return
	}

	AdminID = tmp.ID
	Bot = NewBot(tmp.Token)
	Bot.client.SetLanguage(vkapi.LangRU)
	return
}

func initConfig(config *core.Config) {
	config.Mux.Lock()
	defer config.Mux.Unlock()
	Weather = config.Weather
	Sites = config.Sites
	Schedule = config.Schedule
	Users = config.Users
}
