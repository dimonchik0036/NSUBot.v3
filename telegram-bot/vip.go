package tgbot

import (
	"TelegramBot/jokes"
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	strCmdVIPMenu = "vip"
	strCmdJoke    = "joke"
	strCmdVIPHelp = "viph"
)

func initVipCommands() {
	tgCommands.AddHandler(core.Handler{
		Handler:         vipMenuCommand,
		PermissionLevel: core.PermissionVIP,
	}, strCmdVIPMenu)

	tgCommands.AddHandler(core.Handler{
		Handler:         vipJokeCommand,
		PermissionLevel: core.PermissionVIP,
	}, strCmdJoke)

	tgCommands.AddHandler(core.Handler{
		Handler:         vipHelpCommand,
		PermissionLevel: core.PermissionVIP,
	}, strCmdVIPHelp)
}

func vipMenuCommand(user *core.User, command *core.Command) {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Шутка", addCommand(strCmdJoke, "")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подсказки", addBackArg(strCmdVIPHelp, strCmdVIPMenu)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Закрыть меню", addCommand(strCmdDelMessage, "")),
		),
	)

	sendMessage(user, command, "VIP панель", &markup)
}

func vipJokeCommand(user *core.User, command *core.Command) {
	joke, err := jokes.GetJokes()
	if err != nil {
		sendError(user, command, "Произошла ошибка, повторите попытку.", true)
		return
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButtonText, addCommand(strCmdVIPMenu, "")),
			tgbotapi.NewInlineKeyboardButtonData("Обновить", addCommand(strCmdJoke, "")),
		),
	)

	sendMessage(user, command, joke, &markup)
}

func vipHelpCommand(user *core.User, command *core.Command) {
	sendMessageInNewMessage(user, command, "Вы VIP, поздравляю!\n"+
		"Вам доступна команда /joke.")
}
