package tgbot

import (
	"fmt"
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/dimonchik0036/nsu-bot/news"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"time"
)

func NewsHandler(users *core.Users, news []news.News) {
	fmt.Println(news)
}

func requestHandler(update tgbotapi.Update) {
	user := searchUser(getFrom(update))
	log.Print(user.String() + " пишет: " + getMessage(update))
	messageHandler(user, update)
}

func getFrom(update tgbotapi.Update) *tgbotapi.User {
	switch {
	case update.Message != nil:
		return update.Message.From
	case update.CallbackQuery != nil:
		return update.CallbackQuery.From
	default:
		return &tgbotapi.User{}
	}
}

func getMessage(update tgbotapi.Update) string {
	switch {
	case update.Message != nil:
		return update.Message.Text
	case update.CallbackQuery != nil:
		return update.CallbackQuery.Data
	default:
		return ""
	}
}

func messageHandler(user *core.User, update tgbotapi.Update) {
	cmd := commandSelect(user, update)
	user.ContinuationCommand = false

	go func() {
		user.QueueMux.Lock()
		defer user.QueueMux.Unlock()
		core.CommandHandler(user, cmd, strings.ToLower, tgCommands)
	}()
}

func commandSelect(user *core.User, update tgbotapi.Update) *core.Command {
	switch {
	case update.Message != nil:
		return commandSelectMessage(user, update.Message)
	case update.CallbackQuery != nil:
		return commandSelectCallback(user, update.CallbackQuery)
	default:
		return &core.Command{}
	}
}

func commandSelectMessage(user *core.User, message *tgbotapi.Message) *core.Command {
	if message.IsCommand() {
		cmd := core.SearchCommand(message.Text[1:], " ")
		return &cmd
	}

	if user.ContinuationCommand {
		if err := core.ProcessingInputByFieldNames(message.Text, user.CurrentCommand); err != nil {
			log.Printf("%s error: %s", user.String(), err.Error())
		}
		return user.CurrentCommand
	}

	return &core.Command{}
}

func commandSelectCallback(user *core.User, callback *tgbotapi.CallbackQuery) *core.Command {
	cmd := core.ProcessingInput(callback.Data, ";")
	return &cmd
}

func searchUser(from *tgbotapi.User) *core.User {
	user := tgUsers.TgUser(int64(from.ID))
	if user == nil {
		return newUser(from)
	}

	updateUser(user, from)
	return user
}

func newUser(from *tgbotapi.User) *core.User {
	var user core.User
	defer tgUsers.SetTgUser(&user)
	user.ID = int64(from.ID)
	user.Platform = core.PlatformTg
	user.Permission = core.PermissionUser

	now := time.Now().Unix()
	user.DateCreated = now
	user.DateLastActivities = now
	user.FirstName = from.FirstName
	user.LastName = from.LastName
	user.Username = from.UserName

	tgBot.Send(tgbotapi.NewMessage(tgAdminID, "Новый пользователь:\n"+user.NewUserString("@")))
	return &user
}

func updateUser(user *core.User, from *tgbotapi.User) {
	user.DateLastActivities = time.Now().Unix()
	if user.FirstName != from.FirstName {
		log.Printf("%s -> %s", user.FirstName, from.FirstName)
		user.FirstName = from.FirstName
	}

	if user.LastName != from.LastName {
		log.Printf("%s -> %s", user.LastName, from.LastName)
		user.LastName = from.LastName
	}

	if user.Username != from.UserName {
		log.Printf("%s -> %s", user.Username, from.UserName)
		user.Username = from.UserName
	}

}
