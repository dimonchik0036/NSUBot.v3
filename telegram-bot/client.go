package tgbot

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"log"
)

func Processing(config *core.Config) {
	loadConfig()
	log.Printf("Телеграм-бот запущен")

}
