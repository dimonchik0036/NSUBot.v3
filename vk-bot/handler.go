package vkbot

import (
	"fmt"
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/dimonchik0036/nsu-bot/news"
	"github.com/dimonchik0036/vk-api"
	"log"
	"strings"
	"time"
)

func RequestHandler(message *vkapi.LPMessage) {
	user := searchUser(message.FromID)
	log.Print(user.String() + " пишет: " + message.Text)
	MessageHandler(user, message.Text)
}

func MessageHandler(user *core.User, message string) {
	if message == "" {
		return
	}

	cmd := CommandSelect(user, message)
	go func() {
		user.QueueMux.Lock()
		defer user.QueueMux.Unlock()
		core.CommandHandler(user, cmd, strings.ToLower, Commands)
	}()
}

func CommandSelect(user *core.User, input string) core.Command {
	if !user.ContinuationCommand {
		return core.SearchCommand(input, " ")
	}

	if err := core.ProcessingInputByFieldNames(input, user.CurrentCommand); err != nil {
		Bot.SendMessage(user.ID, "Произошла непредвиденная ошибка, повторите попытку.")
		log.Printf(user.String(), "error: ", err.Error())
	}

	return *user.CurrentCommand
}

func NewsHandler(users *core.Users, news []news.News) {
	fmt.Println(news)
}

func searchUser(id int64) *core.User {
	user := Users.VkUser(id)
	if user == nil {
		return newUser(id)
	}

	return user
}

func newUser(id int64) *core.User {
	var user core.User
	defer Users.SetVkUser(&user)
	user.ID = id
	user.Platform = core.PlatformVk
	user.Permission = core.PermissionUser

	now := time.Now().Unix()
	user.DateCreated = now
	user.DateLastActivities = now

	vkUsers, err := Bot.client.UsersInfo(vkapi.NewDstFromUserID(id), vkapi.UserFieldDomain)
	if err != nil || len(vkUsers) == 0 {
		log.Print(id, "not found")
		Bot.SendMessage(AdminID, "Новый пользователь:\n"+user.NewUserString("vk.com/"))
		return &user
	}

	user.FirstName = vkUsers[0].FirstName
	user.LastName = vkUsers[0].LastName
	user.Username = vkUsers[0].Domain
	Bot.SendMessage(AdminID, "Новый пользователь:\n"+user.NewUserString("vk.com/"))
	return &user
}
