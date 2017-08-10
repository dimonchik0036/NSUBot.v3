package vkbot

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/dimonchik0036/nsu-bot/news"
	"github.com/dimonchik0036/vk-api"
	"log"
	"strings"
	"time"
)

func requestHandler(message *vkapi.LPMessage) {
	vkBot.client.MarkMessageAsRead(message.ID)
	user := searchUser(message.FromID)
	log.Print(user.String() + " пишет: " + message.Text)
	messageHandler(user, message.Text)
}

func messageHandler(user *core.User, message string) {
	if message == "" {
		return
	}

	cmd := commandSelect(user, message)
	user.ContinuationCommand = false

	go func() {
		user.QueueMux.Lock()
		defer user.QueueMux.Unlock()
		core.CommandHandler(user, cmd, strings.ToLower, vkCommands)
	}()
}

func commandSelect(user *core.User, input string) *core.Command {
	if !user.ContinuationCommand || input == "отмена" {
		cmd := core.SearchCommand(input, " ")
		return &cmd
	}

	if err := core.ProcessingInputByFieldNames(input, user.CurrentCommand); err != nil {
		vkBot.SendMessage(user.ID, "Произошла непредвиденная ошибка, повторите попытку.")
		log.Printf(user.String(), "error: ", err.Error())
	}

	return user.CurrentCommand
}

func NewsHandler(users *core.Users, news []news.News, title string) {
	vkUsers := users.VkUsers()
	for _, user := range vkUsers {
		for _, n := range news {
			vkBot.SendMessage(user.ID, title+"\n"+n.URL+"\n"+n.Title+"\n"+n.Decryption)
		}
	}
}

func searchUser(id int64) *core.User {
	user := vkUsers.VkUser(id)
	if user == nil {
		return newUser(id)
	}

	return user
}

func newUser(id int64) *core.User {
	var user core.User
	defer vkUsers.SetVkUser(&user)
	user.ID = id
	user.Platform = core.PlatformVk
	user.Permission = core.PermissionUser

	now := time.Now().Unix()
	user.DateCreated = now
	user.DateLastActivities = now

	vkUsers, err := vkBot.client.UsersInfo(vkapi.NewDstFromUserID(id), vkapi.UserFieldDomain)
	if err != nil || len(vkUsers) == 0 {
		log.Print(id, "not found")
		vkBot.SendMessage(vkAdminID, "Новый пользователь:\n"+user.NewUserString("vk.com/"))
		return &user
	}

	user.FirstName = vkUsers[0].FirstName
	user.LastName = vkUsers[0].LastName
	user.Username = vkUsers[0].Domain
	vkBot.SendMessage(vkAdminID, "Новый пользователь:\n"+user.NewUserString("vk.com/"))
	return &user
}
