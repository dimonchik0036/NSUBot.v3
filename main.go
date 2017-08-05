package main

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/dimonchik0036/nsu-bot/news"
	"github.com/dimonchik0036/nsu-bot/nsuweather"
	"github.com/dimonchik0036/nsu-bot/telegram-bot"
	"github.com/dimonchik0036/nsu-bot/vk-bot"
	"log"
	"os"
	"time"
)

const (
	Usage = "Usage: nsu-bot [<PLATFORM>]\n" +
		"PLATFORM is 'vk' or 'tg'"
)

func main() {
	var processing func(*core.Config)
	var newsHandler func(*core.Users, []news.News)
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--help":
			print(Usage)
			return
		case "vk":
			processing = vkbot.Processing
			newsHandler = vkbot.NewsHandler
		case "tg":
			processing = tgbot.Processing
			newsHandler = tgbot.NewsHandler
		default:
			print("This PLATFORM not found\n" + Usage)
			return
		}
	} else {
		newsHandler = func(users *core.Users, news []news.News) {
			go vkbot.NewsHandler(users, news)
			go tgbot.NewsHandler(users, news)
		}

		processing = func(config *core.Config) {
			go vkbot.Processing(config)
			go tgbot.Processing(config)

			for {
				time.Sleep(24 * time.Hour)
			}
		}
	}
	//initLog() //comment while testing
	config := core.LoadConfig()
	go UpdateSection(config, newsHandler)

	processing(config)
}

func initLog() {
	file, err := os.OpenFile("syslog"+time.Now().Format("2006-01-02T15-04-05")+".txt", os.O_CREATE|os.O_WRONLY, os.FileMode(0700))
	if err != nil {
		log.Panic(err)
	}

	log.SetOutput(file)
}

func UpdateSection(config *core.Config, newsHandler func(*core.Users, []news.News)) {
	go weatherUpdate(config.Weather, 2*time.Minute)

	go save(config, 20*time.Second, 5*time.Minute)

	go sitesUpdate(config.Sites, 3*time.Minute, newsHandler)
}

func weatherUpdate(weather *nsuweather.Weather, duration time.Duration) {
	for {
		weather.Update()
		time.Sleep(duration)
	}
}

func save(config *core.Config, delay time.Duration, duration time.Duration) {
	time.Sleep(delay)
	for {
		config.Save()
		time.Sleep(duration)
	}
}

func sitesUpdate(sites *core.Sites, duration time.Duration, handler func(*core.Users, []news.News)) {
	for {
		sites.Update(handler)
		time.Sleep(duration)
	}
}
