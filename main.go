package main

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/dimonchik0036/nsu-bot/telegram-bot"
	"github.com/dimonchik0036/nsu-bot/vk-bot"
	"os"
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

	var processing func(core.Config)
	switch os.Args[1] {
	case "vk":
		processing = vkbot.Processing
	case "tg":
		processing = tgbot.Processing
	default:
		print("This PLATFORM not found\n" + Usage)
		return
	}
	config := core.NewConfig()
	checkConfig(config)

	processing(config)
}

func checkConfig(config core.Config) {
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
}
