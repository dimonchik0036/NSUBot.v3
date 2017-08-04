package vkbot

import (
	"github.com/dimonchik0036/vk-api"
	"log"
)

func RequestHandler(message *vkapi.LPMessage) {
	log.Print(message.String())
}
