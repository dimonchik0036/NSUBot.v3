package tgbot

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	strCmdMainMenu   = "menu"
	strCmdSiteMenu   = "m_si"
	strCmdSubMenu    = "m_su"
	strCmdOptionMenu = "m_opt"
)
const (
	backButtonText = "« Назад"
	mainButtonText = "« В начало"
)

func mainMenuCommand(user *core.User, command *core.Command) {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Температура", addBackArg(strCmdWeather, strCmdMainMenu)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Управление подписками", addCommand(strCmdSubMenu, "")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Дополнительно", addCommand(strCmdOptionMenu, "")),
		),
	)

	sendMessage(user, command, "Главное меню", &markup)
}

func optionMenuCommand(user *core.User, command *core.Command) {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("История обновлений", addCommand(strCmdBotNewsList, "")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Оставить отзыв", addCommand(strCmdFeedback, "")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButtonText, addCommand(strCmdMainMenu, "")),
		),
	)

	sendMessage(user, command, "Дополнительно", &markup)
}

func subMenuCommand(user *core.User, command *core.Command) {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сайты", addCommand(strCmdSiteMenu, "")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButtonText, addCommand(strCmdMainMenu, "")),
		),
	)

	sendMessage(user, command, "Управление подписками", &markup)
}

func backButton(command *core.Command) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(backButtonText, command.GetArg(core.StrPreviousCmd))
}

func addCommand(command string, args string) string {
	return "c=" + command + "*" + args
}

func addBackArg(target string, previous string) string {
	return core.GenerateCommandString(target, map[string]string{core.StrPreviousCmd: "c=" + previous})
}

func siteMenuCommand(user *core.User, command *core.Command) {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сайт НГУ", addCommand(strCmdShowSite, "10")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сайт ФИТ", addCommand(strCmdShowSite, "00")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сайт ММФ", addCommand(strCmdShowSite, "20")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сайт ФП", addCommand(strCmdShowSite, "30")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сайт ФилФ", addCommand(strCmdShowSite, "40")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButtonText, addCommand(strCmdSubMenu, "")),
			tgbotapi.NewInlineKeyboardButtonData(mainButtonText, addCommand(strCmdMainMenu, "")),
		),
	)
	sendMessage(user, command, "Выберите сайт", &markup)
	return
}
