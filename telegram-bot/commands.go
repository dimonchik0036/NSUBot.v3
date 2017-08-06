package tgbot

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	strCallbackID = "callbackID"
	strMessageID  = "messageID"
)

var tgCommands core.Handlers

func initCommands() {
	tgCommands = core.Handlers{}
	tgCommands.AddHandler(helpHandler(), "help", "помощь")
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
	По всем вопросам можно обратиться к @dimonchik0036.`
	if callbackID := command.Args[strCallbackID]; callbackID != "" {
		tgBot.Send(tgbotapi.NewEditMessageText(user.ID, int(command.GetArgInt64(strMessageID)), defaultHelp))
	} else {
		tgBot.Send(tgbotapi.NewMessage(user.ID, defaultHelp))
	}
}
