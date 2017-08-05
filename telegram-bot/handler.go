package tgbot

import (
	"fmt"
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/dimonchik0036/nsu-bot/news"
)

func NewsHandler(users *core.Users, news []news.News) {
	fmt.Println(news)
}
