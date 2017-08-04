package vkbot

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"log"
)

func Processing(config *core.Config) {
	//initLog() //comment while testing
	loadConfig()
	log.Printf("Вк-бот запущен")
	bot := NewBot()
	updates := bot.UpdateChan()

	for update := range updates {
		if update.Message == nil || !update.IsNewMessage() || update.Message.Outbox() {
			continue
		}

		RequestHandler(update.Message)
	}
}
