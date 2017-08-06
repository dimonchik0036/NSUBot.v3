package tgbot

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func Processing(config *core.Config) {
	loadTgConfig()
	initConfig(config)
	initCommands()
	initAdminCommands()
	log.Printf("Телеграм-бот запущен")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := tgBot.GetUpdatesChan(u)
	if err != nil {
		log.Panicf("Tg error: %s", err.Error())
	}

	for update := range updates {
		requestHandler(update)
	}
}
