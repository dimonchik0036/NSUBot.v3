package tgbot

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/dimonchik0036/nsu-bot/nsuschedule"
	"github.com/dimonchik0036/nsu-bot/nsuweather"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strconv"
)

var tgAdminID int64
var tgBot *tgbotapi.BotAPI
var tgUsers *core.Users
var tgWeather *nsuweather.Weather
var tgSites *core.Sites
var tgSchedule *nsuschedule.Schedule
var globalConfig *core.Config
var tmpUserList []*core.User

func loadTgConfig() {
	data, err := ioutil.ReadFile(".tg_config")
	if err != nil {
		log.Printf("Tg config not found: %s", err.Error())
		return
	}

	tmp := struct {
		ID    int64
		Token string
	}{}

	if err := yaml.Unmarshal(data, &tmp); err != nil {
		log.Printf("Tg config: yaml throw error: %s", err.Error())
		return
	}

	tgAdminID = tmp.ID
	bot, err := tgbotapi.NewBotAPI(tmp.Token)
	if err != nil {
		log.Panicf("Bot is offline: %s", err.Error())
	}

	tgBot = bot
	if _, err := tgBot.Send(tgbotapi.NewMessage(tgAdminID, "Я запущен")); err != nil {
		log.Panicf("Tg error: %s", err.Error())
	}

	return
}

func initConfig(config *core.Config) {
	config.Mux.Lock()
	defer config.Mux.Unlock()
	tgWeather = config.Weather
	tgSites = config.Sites
	tgSchedule = config.Schedule
	tgUsers = config.Users
	globalConfig = config
}

func sendMessage(user *core.User, command *core.Command, text string, markup *tgbotapi.InlineKeyboardMarkup) {
	id := checkCallback(command.Args)
	if id != 0 {
		msg := tgbotapi.NewEditMessageText(user.ID, id, text)
		msg.ReplyMarkup = markup
		tgBot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(user.ID, text)
		msg.ReplyMarkup = markup
		tgBot.Send(msg)
	}
}

func sendMessageInNewMessage(user *core.User, command *core.Command, text string) {
	if command.Args == nil {
		tgBot.Send(tgbotapi.NewMessage(user.ID, text))
		return
	}

	if callbackID := command.Args[strCallbackID]; callbackID != "" {
		msg := tgbotapi.NewEditMessageText(user.ID, int(command.GetArgInt64(strMessageID)), text)
		markup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(backButton(command)))
		msg.ReplyMarkup = &markup
		tgBot.Send(msg)
	} else {
		tgBot.Send(tgbotapi.NewMessage(user.ID, text))
	}
}

func sendError(user *core.User, command *core.Command, text string) {
	callbackID := command.GetArg(strCallbackID)
	if callbackID == "" {
		tgBot.Send(tgbotapi.NewMessage(user.ID, text))
	} else {
		tgBot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(callbackID, text))
	}
}

func checkCallback(args map[string]string) int {
	if args == nil {
		return 0
	}

	if args[strCallbackID] == "" {
		return 0
	}

	id, _ := strconv.Atoi(args[strMessageID])
	return id
}
