package vkbot

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/dimonchik0036/nsu-bot/news"
	"strconv"
	"strings"
)

var vkCommands core.Handlers

func initCommands() {
	vkCommands = core.Handlers{}
	vkCommands.AddHandler(weatherHandler(), "w", "weather", "погода", "температура")
	vkCommands.AddHandler(helpHandler(), "\"help\"", "help", "помощь")
	vkCommands.AddHandler(welcomeHandler(), "\"привет\"", "привет", "хай")
	vkCommands.AddHandler(subscriptionsHandler(), "подписки")
	vkCommands.AddHandler(myHandler(), "мой", "моя", "мои", "моё")
	vkCommands.AddHandler(unSubscriptionsHandler(), "отписка")
	vkCommands.AddHandler(cancelHandler(), "отмена", "\"отмена\"")
}

func myHandler() core.Handler {
	return core.Handler{
		Handler: myCommand,
	}
}

func myCommand(user *core.User, command *core.Command) {
	if len(command.ArgsStr) != 0 {
		switch command.ArgsStr[0] {
		case "подписки":
			command.Command = "отписка"
			command.ArgsStr = []string{}
			unsubscriptionsCommand(user, command)
		}
	}
}

func cancelHandler() core.Handler {
	return core.Handler{
		Handler: cancelCommand,
	}
}

func cancelCommand(user *core.User, command *core.Command) {
	user.ContinuationCommand = false
	vkBot.SendMessage(user.ID, "Прервано...")
}

func welcomeHandler() core.Handler {
	return core.Handler{
		Handler: welcomeCommand,
	}
}

func welcomeCommand(user *core.User, command *core.Command) {
	vkBot.SendMessage(user.ID, "Привет, "+user.FirstName+"!\n"+
		`Я - Помощник, если вы хотите посмотреть доступные команды, то напишите мне слово "help".

		Если хотите узнать что-то дополнительное обо мне, то пишите @dimonchik0036(ему).`)
}

func weatherHandler() core.Handler {
	return core.Handler{
		Handler: weatherCommand,
	}
}

func weatherCommand(user *core.User, command *core.Command) {
	vkBot.SendMessage(user.ID, vkWeather.ShowWeather()+"\n"+
		"Время последнего обновления: "+ vkWeather.ShowTime())
}

func helpHandler() core.Handler {
	return core.Handler{
		PermissionLevel: 0,
		Handler:         helpCommand,
	}
}

func helpCommand(user *core.User, command *core.Command) {
	defaultHelp := `Список команд:
	отмена - Прерывает любую цепочку команд.

	w | weather | погода | температура - Показывает текущую температуру около НГУ.

	подписки - Показывает доступные сайты для подписки на рассылку.

	мои подписки - Упраление подписками.
	По всем вопросам можно обратиться к @dimonchik0036(нему).`
	vkBot.SendMessage(user.ID, defaultHelp)
}

func subscriptionsHandler() core.Handler {
	return core.Handler{
		Handler: subscriptionsCommand,
	}
}

func subscriptionsCommand(user *core.User, command *core.Command) {
	site := command.Args["site"]
	if site == "" {
		selectSite(user, command)
		return
	}

	if len(site) != 1 || !strings.ContainsAny(site, "01234") {
		vkBot.SendMessage(user.ID, "Вне диапазона")
	}

	if len(command.ArgsStr) == 0 {
		selectNews(user, command)
		return
	}

	var str string
	for _, arg := range command.ArgsStr {
		i, err := strconv.Atoi(arg)
		if err != nil {
			continue
		}

		sites := getNews(command.Args["site"])
		if i < 0 || i >= len(sites) {
			continue
		}
		str += sites[i].URL + "\n"
		vkSites.Sub(sites[i].URL, user)
	}
	if str != "" {
		vkBot.SendMessage(user.ID, "Вы были подписаны на:\n"+str)
	} else {
		vkBot.SendMessage(user.ID, "Вы выбрали неверные разделы, повторите попытку.")
		user.ContinuationCommand = true
	}
}

func selectSite(user *core.User, command *core.Command) {
	user.ContinuationCommand = true
	user.CurrentCommand = command

	command.FieldNames = []string{"site"}
	vkBot.SendMessage(user.ID, `Выберите сайт и напиши его цифру:
	0 - Сайт НГУ
	1 - Сайт ФИТ
	2 - Сайт ФП
	3 - Сайт ММФ
	4 - Сайт ФилФ`)
}

func selectNews(user *core.User, command *core.Command) {
	user.ContinuationCommand = true
	user.CurrentCommand = command

	command.MoreArgs = true
	text := "Напишите номера разделов, чтобы подписаться:\n"

	for i, site := range getNews(command.Args["site"]) {
		text += strconv.Itoa(i) + " - " + site.Title + "\n"
	}

	text += "Напишите \"отмена\", если хотите прервать выбор."

	vkBot.SendMessage(user.ID, text)
}

func getNews(number string) []*news.Site {
	switch number {
	case "0":
		return news.NsuNews()
	case "1":
		return news.FitNews()
	case "2":
		return news.FpNews()
	case "3":
		return news.MmfNews()
	case "4":
		return news.PhilosNews()
	default:
		return []*news.Site{}
	}
}

func unSubscriptionsHandler() core.Handler {
	return core.Handler{
		Handler: unsubscriptionsCommand,
	}
}

func unsubscriptionsCommand(user *core.User, command *core.Command) {
	if len(command.ArgsStr) == 0 {
		var str string
		count := 0
		command.Args = map[string]string{}
		for _, s := range vkSites.Sites {
			if s.Users.VkUser(user.ID) != nil {
				count++
				c := strconv.Itoa(count - 1)
				str += c + " - " + s.Site.Title + " " + s.Site.URL + "\n"
				command.Args[c] = s.Site.URL
			}
		}
		command.Args["count"] = strconv.Itoa(count)

		vkBot.SendMessage(user.ID, "Чтобы отписаться от разделов, отправь их номера:\n"+str)
		user.CurrentCommand = command
		user.ContinuationCommand = true
		command.MoreArgs = true
		return
	}

	с := 0
	for _, arg := range command.ArgsStr {
		i, err := strconv.Atoi(arg)
		if err != nil {
			continue
		}

		count, _ := strconv.Atoi(command.Args["count"])
		if i < 0 || i >= count {
			continue
		}

		vkSites.Unsub(command.Args[arg], user)
		с++
	}

	if с != 0 {
		vkBot.SendMessage(user.ID, "Вы были успешно отписаны.")
	} else {
		vkBot.SendMessage(user.ID, "Вы выбрали неверные разделы, повторите попытку.")
		user.ContinuationCommand = true
	}

}
