package tgbot

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"strings"
)

const (
	strCmdReset            = "reset"
	strCmdAdminMenu        = "admin"
	strCmdDelMessage       = "dm"
	strCmdUserList         = "ulist"
	strCmdShowUser         = "suser"
	strCmdSendMessage      = "usend"
	strCmdSendMessageAll   = "sendall"
	strCmdChangePermission = "changeperm"
)

const (
	userOnOnePage = 5
)

func initAdminCommands() {
	tgCommands.AddHandler(core.Handler{
		Handler:         adminMenu,
		PermissionLevel: core.PermissionAdmin,
	}, strCmdAdminMenu)

	tgCommands.AddHandler(core.Handler{
		Handler:         resetCommand,
		PermissionLevel: core.PermissionAdmin,
	}, strCmdReset)

	tgCommands.AddHandler(core.Handler{
		Handler: deleteMessageCommand,
	}, strCmdDelMessage)

	tgCommands.AddHandler(core.Handler{
		Handler: adminIAmGod,
	}, "godmode")

	tgCommands.AddHandler(core.Handler{
		Handler:         adminUserListCommand,
		PermissionLevel: core.PermissionAdmin,
	}, strCmdUserList)

	tgCommands.AddHandler(core.Handler{
		Handler:         adminShowUser,
		PermissionLevel: core.PermissionAdmin,
	}, strCmdShowUser)

	tgCommands.AddHandler(core.Handler{
		Handler:         adminSendMessageUser,
		PermissionLevel: core.PermissionAdmin,
	}, strCmdSendMessage)

	tgCommands.AddHandler(core.Handler{
		Handler:         adminSendMessageAll,
		PermissionLevel: core.PermissionAdmin,
	}, strCmdSendMessageAll)

	tgCommands.AddHandler(core.Handler{
		Handler:         adminChangePerm,
		PermissionLevel: core.PermissionAdmin,
	}, strCmdChangePermission)
}

func adminChangePerm(user *core.User, command *core.Command) {
	strID := command.GetArg("id")
	strPerm := command.GetArg("perm")

	if strID == "" {
		sendError(user, command, "Введите ID и уровень доступа")
		user.ContinuationCommand = true
		user.CurrentCommand = command
		command.FieldNames = []string{"id", "perm"}
		return
	}

	if strPerm == "" {
		sendError(user, command, "Введите уровень доступа")
		user.ContinuationCommand = true
		user.CurrentCommand = command
		command.FieldNames = []string{"perm"}
		return
	}

	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		tgBot.Send(tgbotapi.NewMessage(user.ID, "Неверный формат ID"))
		return
	}

	perm, err := strconv.Atoi(strPerm)
	if err != nil {
		tgBot.Send(tgbotapi.NewMessage(user.ID, "Неверный формат perm"))
		return
	}

	tgUsers.ChangePermission(core.PlatformTg, id, perm)
	tgBot.Send(tgbotapi.NewMessage(user.ID, "Успешно"))
}

func adminIAmGod(user *core.User, command *core.Command) {
	if user.ID == tgAdminID {
		user.Permission = core.PermissionAdmin
		tgBot.Send(tgbotapi.NewMessage(user.ID, "Приветствую тебя, о Достойнейший!"))
		return
	}

	tgBot.Send(tgbotapi.NewMessage(user.ID, "Ты не достоин!"))
}

func adminMenu(user *core.User, command *core.Command) {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Список пользователей", addCommand(strCmdUserList, "f0")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Выключить бота", addCommand(strCmdReset, "")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Закрыть меню", addCommand(strCmdDelMessage, "")),
		),
	)

	sendMessage(user, command, "Панель админа", &markup)
}

func resetCommand(user *core.User, command *core.Command) {
	if command.GetArg("answer") == "yes" || command.GetArg(strCmdArg) == "yes" {
		globalConfig.Reset()
	} else {
		command.FieldNames = []string{"answer"}
		user.ContinuationCommand = true
		user.CurrentCommand = command
		tgBot.Send(tgbotapi.NewMessage(user.ID, "Введите 'yes', чтобы выключить бота"))
	}
}

func deleteMessageCommand(user *core.User, command *core.Command) {
	tgBot.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: user.ID, MessageID: int(command.GetArgInt64(strMessageID))})
	tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(command.GetArg(strCallbackID), "Готово"))
}

func adminUserListCommand(user *core.User, command *core.Command) {
	var pageNumber int
	args := command.GetArg(strCmdArg)
	if len(args) != 0 {
		if args[:1] == "f" {
			tmpUserList = tgUsers.TgUsers()
			args = args[1:]
		}

		pageNumber, _ = strconv.Atoi(args)
	}

	if pageNumber < 0 {
		pageNumber = 0
	}

	if pageNumber*userOnOnePage > len(tmpUserList) {
		pageNumber = len(tmpUserList) / userOnOnePage
	}

	var markup tgbotapi.InlineKeyboardMarkup

	for i, user := range tmpUserList[pageNumber*userOnOnePage:] {
		if i == userOnOnePage {
			break
		}

		markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(func() string {
			if user.Username != "" {
				return user.Username
			} else {
				return user.FirstName + " " + user.LastName
			}
		}(),
			core.GenerateCommandString(strCmdShowUser, map[string]string{core.StrPreviousCmd: addCommand(strCmdUserList, args),
				strCmdArg: strconv.FormatInt(user.ID, 10)}))))
	}

	markup.InlineKeyboard = append(markup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("«", addCommand(strCmdUserList, strconv.Itoa(func() int {
			if pageNumber-1 < 0 {
				return 0
			}
			return pageNumber - 1
		}()))),
		tgbotapi.NewInlineKeyboardButtonData("Назад", addCommand(strCmdAdminMenu, "")),
		tgbotapi.NewInlineKeyboardButtonData("»", addCommand(strCmdUserList, strconv.Itoa(pageNumber+1))),
	))

	sendMessage(user, command, "Всего "+strconv.Itoa(len(tmpUserList))+" пользователей", &markup)
}

func adminShowUser(user *core.User, command *core.Command) {
	log.Printf("arg %s", command.GetArg(core.StrPreviousCmd))
	id, err := strconv.ParseInt(command.GetArg(strCmdArg), 10, 64)
	if err != nil {
		sendError(user, command, "Некорректный ввод ID")
		return
	}

	u := tgUsers.TgUser(id)
	if u == nil {
		sendError(user, command, "Данный ID не найден")
		return
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить привелегии", core.GenerateCommandString(strCmdChangePermission, map[string]string{"id": strconv.FormatInt(id, 10)})),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отправить сообщение", core.GenerateCommandString(strCmdSendMessage, map[string]string{"id": strconv.FormatInt(id, 10)})),
		),
		tgbotapi.NewInlineKeyboardRow(
			backButton(command),
			tgbotapi.NewInlineKeyboardButtonData(mainButtonText, addCommand(strCmdAdminMenu, "")),
		),
	)

	sendMessage(user, command, u.FullString("@"), &markup)
}

func adminSendMessageUser(user *core.User, command *core.Command) {
	args := strings.Split(command.GetArg(strCmdArg), " ")
	if len(args) == 0 {
		user.ContinuationCommand = true
		user.CurrentCommand = command
		command.FieldNames = []string{"id", "text"}
		sendError(user, command, "Введите ID и текст")
		return
	}

	id, err := strconv.ParseInt(command.GetArg("id"), 10, 64)
	if err != nil {
		if len(args) != 0 {
			id, err = strconv.ParseInt(args[0], 10, 64)
			command.SetArg("id", args[0])
			args = command.ArgsStr[1:]
		}

		if err != nil {
			sendError(user, command, "Некорректный ID")
			return
		}
	}

	if len(command.ArgsStr) == 0 {
		user.ContinuationCommand = true
		user.CurrentCommand = command
		command.FieldNames = []string{"text"}
		sendError(user, command, "Наберите текст")
		return
	}

	u := tgUsers.TgUser(id)
	if u == nil {
		sendError(user, command, "Пользователь не найден")
		return
	}

	if _, err := tgBot.Send(tgbotapi.NewMessage(u.ID, strings.Join(command.ArgsStr, " "))); err != nil {
		tgBot.Send(tgbotapi.NewMessage(user.ID, "Ошибка отправки"))
	} else {
		tgBot.Send(tgbotapi.NewMessage(user.ID, "Успешно отправлено"))
	}
}

func adminSendMessageAll(user *core.User, command *core.Command) {
	text := strings.Join(command.ArgsStr, " ")
	var count = 0
	for _, u := range tgUsers.TgUsers() {
		if _, err := tgBot.Send(tgbotapi.NewMessage(u.ID, text)); err != nil {
			count++
			log.Printf("%s %s", u.String(), err.Error())
		}
	}

	sendError(user, command, "Готово, ошибок при отправлении: "+strconv.Itoa(count))
}
