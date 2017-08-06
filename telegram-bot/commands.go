package tgbot

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/dimonchik0036/nsu-bot/news"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
)

const (
	strInNewMessage = "in_new_message"
	strCmdArg       = "c_a"
)

const (
	strCallbackID = "cID"
	strMessageID  = "mID"
)

const (
	strCmdWeather  = "weather"
	strCmdShowSite = "sh"
	strCmdFeedback = "feedback"
)

var tgCommands core.Handlers

func initCommands() {
	tgCommands = core.Handlers{}
	tgCommands.AddHandler(core.Handler{Handler: helpCommand}, "help", "помощь")

	tgCommands.AddHandler(core.Handler{Handler: mainMenuCommand}, strCmdMainMenu, "start")

	tgCommands.AddHandler(core.Handler{Handler: weatherCommand}, strCmdWeather)

	tgCommands.AddHandler(core.Handler{Handler: subMenuCommand}, strCmdSubMenu)
	tgCommands.AddHandler(core.Handler{Handler: siteMenuCommand}, strCmdSiteMenu)
	tgCommands.AddHandler(core.Handler{Handler: showSiteCommand}, strCmdShowSite)
	tgCommands.AddHandler(core.Handler{Handler: feedbackCommand}, strCmdFeedback)
}

func helpCommand(user *core.User, command *core.Command) {
	defaultHelp := `Список команд:
	/menu - Вызывает главное меню.

	/cancel - Прерывает любую цепочку команд.


	По всем вопросам можно обратиться к @dimonchik0036.`

	sendMessage(user, command, defaultHelp, nil)
}

func feedbackCommand(user *core.User, command *core.Command) {
	if len(command.ArgsStr) == 0 {
		user.ContinuationCommand = true
		user.CurrentCommand = command
		command.FieldNames = []string{"text"}
		sendError(user, command, "Наберите свой отзыв")
		return
	}

	tgBot.Send(tgbotapi.NewMessage(tgAdminID, "Получен отзыв от\n"+user.FullString("@")+"\n\n"+strings.Join(command.ArgsStr, " ")))
	tgBot.Send(tgbotapi.NewMessage(user.ID, "Спасибо за отзыв!"))
}

func weatherCommand(user *core.User, command *core.Command) {
	sendMessageInNewMessage(user, command, tgWeather.ShowWeather()+"\n"+
		"Время последнего обновления: "+tgWeather.ShowTime())
}

func showSiteCommand(user *core.User, command *core.Command) {
	args := command.GetArg(strCmdArg)
	if len(args) < 2 {
		sendError(user, command, "Мало аргументов")
		return
	}

	siteNumber, err := strconv.Atoi(args[:1])
	if err != nil || siteNumber < 0 || siteNumber > 5 {
		sendError(user, command, "Диапазон 0-5")
		return
	}

	pageNumber, err := strconv.Atoi(args[1:2])
	if err != nil {
		sendError(user, command, "Номер страницы неверный")
		return
	}

	siteList := news.GetSite(siteNumber)
	if pageNumber < 0 {
		pageNumber = 0
	}

	if pageNumber*5 > len(siteList) {
		pageNumber = len(siteList) / 5
	}
	siteList = siteList[pageNumber*5:]

	if len(args) > 2 {
		subID, err := strconv.Atoi(args[2:3])
		if err == nil && subID < len(siteList) {
			tgSites.ChangeSub(siteList[subID].URL, user)
		}
	}

	var markup tgbotapi.InlineKeyboardMarkup

	for i, site := range siteList {
		if i == 5 {
			break
		}

		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(checkSite(site.URL, user)+site.Title, addCommand(strCmdShowSite, args[:2]+strconv.Itoa(i)))))
	}

	markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("«", addCommand(strCmdShowSite, args[:1]+strconv.Itoa(func() int {
			if pageNumber-1 < 0 {
				return 0
			}
			return pageNumber - 1
		}()))),
		tgbotapi.NewInlineKeyboardButtonData("Назад", addCommand(strCmdSiteMenu, "")),
		tgbotapi.NewInlineKeyboardButtonData("»", addCommand(strCmdShowSite, args[:1]+strconv.Itoa(pageNumber+1))),
	))

	sendMessage(user, command, "Выберете подписки", &markup)
}

func checkSite(url string, user *core.User) string {
	if tgSites.CheckUser(url, user) {
		return "☑️ "
	}

	return "❌"
}
