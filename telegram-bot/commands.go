package tgbot

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/dimonchik0036/nsu-bot/news"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
	"time"
)

const (
	strInNewMessage = "in_new_message"
	strCmdArg       = "c_a"
)

const (
	siteOnOnePage = 5
)

const (
	strCallbackID = "cID"
	strMessageID  = "mID"
)

const (
	strCmdWeather     = "weather"
	strCmdShowSite    = "sh"
	strCmdFeedback    = "feedback"
	strCmdStart       = "start"
	strCmdHelp        = "help"
	strCmdSchedule    = "schedule"
	strCmdAddLb       = "addlb"
	strCmdGetSchedule = "day"
)

var tgCommands core.Handlers

func initCommands() {
	tgCommands = core.Handlers{}
	tgCommands.AddHandler(core.Handler{Handler: helpCommand}, strCmdHelp, "помощь")
	tgCommands.AddHandler(core.Handler{Handler: startCommand}, strCmdStart)
	tgCommands.AddHandler(core.Handler{Handler: mainMenuCommand}, strCmdMainMenu)
	tgCommands.AddHandler(core.Handler{Handler: optionMenuCommand}, strCmdOptionMenu)

	tgCommands.AddHandler(core.Handler{Handler: weatherCommand}, strCmdWeather)

	tgCommands.AddHandler(core.Handler{Handler: subMenuCommand}, strCmdSubMenu)
	tgCommands.AddHandler(core.Handler{Handler: siteMenuCommand}, strCmdSiteMenu)
	tgCommands.AddHandler(core.Handler{Handler: showSiteCommand}, strCmdShowSite)
	tgCommands.AddHandler(core.Handler{Handler: feedbackCommand}, strCmdFeedback)
	tgCommands.AddHandler(core.Handler{Handler: scheduleMenuCommand}, strCmdScheduleMenu)
	tgCommands.AddHandler(core.Handler{Handler: getScheduleCommand}, strCmdGetSchedule)

	initAdminCommands()
	initVipCommands()
	initBotNewsCommand()
	initVkSiteCommand()
}

func getScheduleCommand(user *core.User, command *core.Command) {
	if command.GetArg("day") == "" {
		day, _ := strconv.Atoi(command.GetArg(strCmdArg))
		command.SetArg("day", strconv.Itoa(day))
	}

	group := command.GetArg("group")
	if group == "" {
		tgBot.Send(tgbotapi.NewMessage(user.ID, "Введите группу"))
		user.ContinuationCommand = true
		user.CurrentCommand = command
		command.FieldNames = []string{"group"}
		return
	}
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButtonText, addCommand(strCmdScheduleMenu, "")),
			tgbotapi.NewInlineKeyboardButtonData(mainButtonText, addCommand(strCmdMainMenu, "")),
		),
	)
	day, _ := strconv.Atoi(command.GetArg("day"))
	var flag bool
	var text string
	if day < 0 {
		group, ok := tgSchedule.GetGroup(group)
		var data [6]string
		if ok && len(group) > 5 {
			flag = ok
			data[0] = "Понедельник.\n" + group[0]
			data[1] = "Вторник.\n" + group[1]
			data[2] = "Среда.\n" + group[2]
			data[3] = "Четверг.\n" + group[3]
			data[4] = "Пятница.\n" + group[4]
			data[5] = "Суббота.\n" + group[5]

			for _, g := range data {
				tgBot.Send(tgbotapi.NewMessage(user.ID, g))
			}
			msg := tgbotapi.NewMessage(user.ID, "Готово")
			msg.ReplyMarkup = markup
			tgBot.Send(msg)
		}
	} else {
		t := int(time.Now().Weekday()) + day + 6
		text, flag = tgSchedule.GetDay(group, t)
		switch t {
		case 0:
			text = "Понедельник.\n" + text
		case 1:
			text = "Вторник.\n" + text
		case 2:
			text = "Среда.\n" + text
		case 3:
			text = "Четверг.\n" + text
		case 4:
			text = "Пятница.\n" + text
		case 5:
			text = "Суббота.\n" + text
		case 6:
			text = "Воскресенье.\n" + text
		}
	}

	if !flag {
		tgBot.Send(tgbotapi.NewMessage(user.ID, "Неверная группа, введите другую группу"))
		user.ContinuationCommand = true
		user.CurrentCommand = command
		command.FieldNames = []string{"group"}
		return
	}

	if day >= 0 {
		msg := tgbotapi.NewMessage(user.ID, text)
		msg.ReplyMarkup = markup
		tgBot.Send(msg)
	}
}

func startCommand(user *core.User, command *core.Command) {
	tgBot.Send(tgbotapi.NewMessage(user.ID, "Приветствую!\n"+
		"Теперь я - ваш помощник.\n"+
		"Я позволяю получить быстрый доступ к температуре воздуха или же вы можете подписаться на рассылку новостей с различных сайтов и групп.\n"+
		"\n"+
		"Возможно будет полезным посмотреть /help, чтобы узнать все команды.\n"+
		"\n"+
		"При возникновении вопросов можно оставить /feedback или обратиться напрямую к @dimonchik0036.\n"))

	mainMenuCommand(user, command)
}

func helpCommand(user *core.User, command *core.Command) {
	defaultHelp := "Список команд:\n" +
		"/menu - Вызывает главное меню.\n" +
		"\n" +
		"/cancel - Прерывает любую цепочку команд.\n" +
		"\n" +
		"Через меню можно подписаться на рассылку различных новостей.\n" +
		"Если вы подписаны на какой-то раздел новостей, то как только появится новая публикация на сайте, вам придёт сообщение об этом.\n" +
		"\n" +
		"По всем вопросам можно обратиться к @dimonchik0036."

	tgBot.Send(tgbotapi.NewMessage(user.ID, defaultHelp))
}

func feedbackCommand(user *core.User, command *core.Command) {
	if len(command.ArgsStr) == 0 {
		user.ContinuationCommand = true
		user.CurrentCommand = command
		command.FieldNames = []string{"text"}
		sendError(user, command, "Наберите свой отзыв", true)
		return
	}

	tgBot.Send(tgbotapi.NewMessage(tgAdminID, "Получен отзыв от\n"+user.FullString("@")+"\n\n"+strings.Join(command.ArgsStr, " ")))
	tgBot.Send(tgbotapi.NewMessage(user.ID, "Спасибо за отзыв!"))
}

func weatherCommand(user *core.User, command *core.Command) {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButtonText, addCommand(strCmdMainMenu, "")),
			tgbotapi.NewInlineKeyboardButtonData("Обновить", addCommand(strCmdWeather, "")),
		),
	)

	sendMessage(user, command, tgWeather.ShowWeather()+"\n"+
		"Время последнего обновления: "+tgWeather.ShowTime(), &markup)
}

func showSiteCommand(user *core.User, command *core.Command) {
	args := strings.Split(command.GetArg(strCmdArg), "_")
	if len(args) < 2 {
		sendError(user, command, "Мало аргументов. Попробуйте вернуться назад и повторить попытку.", true)
		return
	}

	siteNumber, err := strconv.Atoi(args[0])
	if err != nil || siteNumber < 0 || siteNumber > siteOnOnePage {
		sendError(user, command, "Диапазон 0-"+strconv.Itoa(siteOnOnePage), true)
		return
	}

	pageNumber, err := strconv.Atoi(args[1])
	if err != nil {
		sendError(user, command, "Номер страницы неверный", true)
		return
	}

	var siteList []*news.Site
	if siteNumber == 5 {
		vkGroupSites.Mux.RLock()
		defer vkGroupSites.Mux.RUnlock()
		siteList = vkGroupSites.Groups
	} else {
		siteList = news.GetSite(siteNumber)
	}
	if pageNumber < 0 {
		pageNumber = 0
	}

	if pageNumber*siteOnOnePage > len(siteList) {
		pageNumber = len(siteList) / siteOnOnePage
	}

	if len(args) > 2 {
		subID, err := strconv.Atoi(args[2])
		if err == nil && (subID+pageNumber*siteOnOnePage < len(siteList)) {
			tgSites.ChangeSub(siteList[pageNumber*siteOnOnePage+subID].URL, user)
		}
	}

	backCmd := command.GetArg(core.StrPreviousCmd)
	if strings.HasPrefix(backCmd, "c=") {
		backCmd = backCmd[2:]
	} else {
		backCmd = strCmdSubMenu
	}

	var markup tgbotapi.InlineKeyboardMarkup

	for i, site := range siteList[pageNumber*siteOnOnePage:] {
		if i == siteOnOnePage {
			break
		}

		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(checkSite(site.URL, user)+site.Title, addBackArg(strCmdShowSite+"*"+args[0]+"_"+args[1]+"_"+strconv.Itoa(i), backCmd))))
	}

	markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("«", addBackArg(strCmdShowSite+"*"+args[0]+"_"+strconv.Itoa(pageNumber-1), backCmd)),
		tgbotapi.NewInlineKeyboardButtonData("Назад", addCommand(backCmd, "")),
		tgbotapi.NewInlineKeyboardButtonData("»", addBackArg(strCmdShowSite+"*"+args[0]+"_"+strconv.Itoa(pageNumber+1), backCmd)),
	))

	sendMessage(user, command, "Страница: "+strconv.Itoa(pageNumber+1)+"/"+strconv.Itoa(len(siteList)/siteOnOnePage+func() int {
		if len(siteList) == 0 {
			return 1
		} else if len(siteList)%siteOnOnePage == 0 {
			return 0
		} else {
			return 1
		}
	}())+"\nВыберете подписки", &markup)
}

func checkSite(url string, user *core.User) string {
	if tgSites.CheckUser(url, user) {
		return "☑️ "
	}

	return "❌"
}
