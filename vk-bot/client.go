package vkbot

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"log"
)

func Processing(config *core.Config) {
	loadVkConfig()
	initConfig(config)
	initCommands()
	log.Printf("Вк-бот запущен")
	updates := Bot.UpdateChan()

	for update := range updates {
		if update.Message == nil || !update.IsNewMessage() || update.Message.Outbox() {
			continue
		}

		RequestHandler(update.Message)
	}
}
