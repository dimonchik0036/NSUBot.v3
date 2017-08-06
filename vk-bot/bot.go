package vkbot

import (
	"github.com/dimonchik0036/vk-api"
	"log"
)

type botSt struct {
	client *vkapi.Client
}

func newBot(token string) *botSt {
	bot, _ := vkapi.NewClientFromToken(token)
	return &botSt{client: bot}
}

func (b *botSt) SendMessage(userID int64, text string) {
	b.client.SendMessage(vkapi.MessageConfig{Destination: vkapi.Destination{UserID: userID}, Message: text})
}

func (b *botSt) UpdateChan() vkapi.LPChan {
	if err := b.client.InitLongPoll(0, 2); err != nil {
		log.Panic(err)
	}

	updates, _, err := b.client.GetLPUpdatesChan(100, vkapi.LPConfig{25, 0})
	if err != nil {
		log.Panic(err)
	}

	return updates
}
