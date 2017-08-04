package vkbot

import (
	"github.com/dimonchik0036/vk-api"
	"log"
)

type Bot struct {
	client *vkapi.Client
}

func NewBot() *Bot {
	bot, _ := vkapi.NewClientFromToken(BotToken)
	return &Bot{client: bot}
}

func (b *Bot) SendMessage(userID int64, text string) {
	b.client.SendMessage(vkapi.MessageConfig{Destination: vkapi.Destination{UserID: userID}, Message: text})
}

func (b *Bot) UpdateChan() vkapi.LPChan {
	if err := b.client.InitLongPoll(0, 2); err != nil {
		log.Panic(err)
	}

	updates, _, err := b.client.GetLPUpdatesChan(100, vkapi.LPConfig{25, 0})
	if err != nil {
		log.Panic(err)
	}

	return updates
}
