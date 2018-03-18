package tgbot

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	strCmdMainMenu     = "menu"
	strCmdSiteMenu     = "m_si"
	strCmdSubMenu      = "m_su"
	strCmdOptionMenu   = "m_opt"
	strCmdScheduleMenu = "m_sch"
	strCmdLabelsMenu   = "m_lb"
)
const (
	backButtonText = "« Назад"
	mainButtonText = "« В начало"
)

func mainMenuCommand(user *core.User, command *core.Command) {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Температура", addCommand(strCmdWeather, "")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Расписание", addCommand(strCmdScheduleMenu, "")),
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

func scheduleMenuCommand(user *core.User, command *core.Command) {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("На сегодня", addCommand(strCmdGetSchedule, "0")), tgbotapi.NewInlineKeyboardButtonData("На завтра", addCommand(strCmdGetSchedule, "1"))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("На всю неделю", addCommand(strCmdGetSchedule, "-1"))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButtonText, addCommand(strCmdMainMenu, ""))),
	)

	sendMessage(user, command, "Расписание", &markup)
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
			tgbotapi.NewInlineKeyboardButtonData("VK группы", addBackArg(strCmdShowSite+"*5_0", strCmdSubMenu)),
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
			tgbotapi.NewInlineKeyboardButtonData("Сайт НГУ", addBackArg(strCmdShowSite+"*1_0", strCmdSiteMenu)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сайт ФИТ", addBackArg(strCmdShowSite+"*0_0", strCmdSiteMenu)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сайт ММФ", addBackArg(strCmdShowSite+"*2_0", strCmdSiteMenu)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сайт ФП", addBackArg(strCmdShowSite+"*3_0", strCmdSiteMenu)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сайт ФилФ", addBackArg(strCmdShowSite+"*4_0", strCmdSiteMenu)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButtonText, addCommand(strCmdSubMenu, "")),
			tgbotapi.NewInlineKeyboardButtonData(mainButtonText, addCommand(strCmdMainMenu, "")),
		),
	)
	sendMessage(user, command, "Выберите сайт", &markup)
	return
}
