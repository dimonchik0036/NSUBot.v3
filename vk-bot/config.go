package vkbot

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var AdminID int64
var BotToken string

func loadConfig() {
	data, err := ioutil.ReadFile(".vk_config")
	if err != nil {
		log.Printf("Vk config not found: %s", err.Error())
		return
	}

	tmp := struct {
		ID    int64
		Token string
	}{}

	if err := yaml.Unmarshal(data, &tmp); err != nil {
		log.Printf("Vk config: yaml throw error: %s", err.Error())
		return
	}

	AdminID = tmp.ID
	BotToken = tmp.Token
	return
}
