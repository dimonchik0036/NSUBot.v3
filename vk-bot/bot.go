package vkbot

import (
	"github.com/dimonchik0036/vk-api"
	"log"
)

type BotSt struct {
	client *vkapi.Client
}

func NewBot(token string) *BotSt {
	bot, _ := vkapi.NewClientFromToken(token)
	return &BotSt{client: bot}
}

func (b *BotSt) SendMessage(userID int64, text string) {
	b.client.SendMessage(vkapi.MessageConfig{Destination: vkapi.Destination{UserID: userID}, Message: text})
}

func (b *BotSt) UpdateChan() vkapi.LPChan {
	if err := b.client.InitLongPoll(0, 2); err != nil {
		log.Panic(err)
	}

	updates, _, err := b.client.GetLPUpdatesChan(100, vkapi.LPConfig{25, 0})
	if err != nil {
		log.Panic(err)
	}

	return updates
}