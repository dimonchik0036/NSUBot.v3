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

var vkAdminID int64
var vkBot *botSt
var vkUsers *core.Users
var vkWeather *nsuweather.Weather
var vkSites *core.Sites
var vkSchedule *nsuschedule.Schedule

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

	vkAdminID = tmp.ID
	vkBot = newBot(tmp.Token)
	vkBot.client.SetLanguage(vkapi.LangRU)
	return
}

func initConfig(config *core.Config) {
	config.Mux.Lock()
	defer config.Mux.Unlock()
	vkWeather = config.Weather
	vkSites = config.Sites
	vkSchedule = config.Schedule
	vkUsers = config.Users
}
