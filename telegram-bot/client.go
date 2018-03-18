package tgbot

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func Processing(config *core.Config) {
	loadTgConfig()
	initConfig(config)
	initBotNews()
	//initVkSites()
	initCommands()
	for _, u := range tgUsers.TgUsers() {
		_, err := tgBot.Send(tgbotapi.NewMessage(u.ID, "Полный редизайн бота (к сожалению пришлось отказаться от встроенных клавиатур, возможно, что временно).\n" +
			"Завершено расписание и закладки.\n" +
			"Из-за обновления ядра новостей, подписки были сброшены."))
		if err != nil {
			log.Println("Ошибка у ", u.ID)
		}
	}
	return
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
