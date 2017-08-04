package main

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/dimonchik0036/nsu-bot/telegram-bot"
	"github.com/dimonchik0036/nsu-bot/vk-bot"
	"os"
	"time"
)

const (
	Usage = "Usage: nsu-bot <PLATFORM>\n" +
		"PLATFORM is 'vk' or 'tg'"
)

func main() {
	if len(os.Args) < 2 {
		print(Usage)
		return
	}

	var processing func(*core.Config)
	switch os.Args[1] {
	case "vk":
		processing = vkbot.Processing
	case "tg":
		processing = tgbot.Processing
	default:
		print("This PLATFORM not found\n" + Usage)
		return
	}

	config := core.LoadConfig()
	go UpdateSection(config)
	processing(config)
}

func UpdateSection(config *core.Config) {
	go func() {
		for {
			config.Weather.Update()
			time.Sleep(2 * time.Minute)
		}
	}()

	go func() {
		time.Sleep(20 * time.Second)
		for {
			config.Save()
			time.Sleep(5 * time.Minute)
		}
	}()
	//config.Sites.Update()
}
