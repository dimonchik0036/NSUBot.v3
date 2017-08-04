package main

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/dimonchik0036/nsu-bot/telegram-bot"
	"github.com/dimonchik0036/nsu-bot/vk-bot"
)

func main() {
	println("start")
	config := core.NewConfig()

	/*for _, s := range config.Sites.Sites {
		n, err := s.Site.Update(2)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(s.Site.Title)
		for i := range n {
			fmt.Println(n[i])
			//fmt.Println(time.Unix(n[i].Date, 0).Format(news.TimeLayout), "\n"+n[i].Title+" "+n[i].URL+"\n"+n[i].Decryption)
		}
		fmt.Println()
	}*/

	go vkbot.Processing(config)
	tgbot.Processing(config)
}
