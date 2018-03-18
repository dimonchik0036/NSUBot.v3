package tgbot

import (
	"errors"
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/dimonchik0036/nsu-bot/news"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"sort"
	"strings"
	"sync"
)

type ByID []*news.Site

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }

var vkGroupSites struct {
	Mux    sync.RWMutex
	Groups []*news.Site
}

const (
	strCmdAddVKSite  = "addvksite"
	strCmdDelVKSite  = "delvksite"
	strCmdVkSiteMenu = "vmenu"
)

func initVkSites() {
	tgSites.Mux.RLock()
	defer tgSites.Mux.RUnlock()

	var groups []*news.Site
	for key, site := range tgSites.Sites {
		if strings.HasPrefix(key, news.VkHref) {
			groups = append(groups, site.Site)
		}
	}

	sort.Sort(ByID(groups))
	vkGroupSites.Mux.Lock()
	defer vkGroupSites.Mux.Unlock()
	vkGroupSites.Groups = groups
}

func initVkSiteCommand() {
	tgCommands.AddHandler(core.Handler{
		Handler:         vkSiteMenuCommand,
		PermissionLevel: core.PermissionAdmin,
	}, strCmdVkSiteMenu)

	tgCommands.AddHandler(core.Handler{
		Handler:         addVkSiteCommand,
		PermissionLevel: core.PermissionAdmin,
	}, strCmdAddVKSite)

	tgCommands.AddHandler(core.Handler{
		Handler:         deleteVkSiteCommand,
		PermissionLevel: core.PermissionAdmin,
	}, strCmdDelVKSite)
}

func addVkSite(domain string, title string) error {
	site := core.Site{
		Site: news.NewVkSite(int64(len(vkGroupSites.Groups)), domain, title),
	}

	_, err := site.Site.Update(2)
	if err != nil {
		return errors.New("Домен не найден")
	}

	vkGroupSites.Mux.Lock()
	defer vkGroupSites.Mux.Unlock()
	for _, s := range vkGroupSites.Groups {
		if s.OptionalURL == domain {
			return errors.New("Уже существует")
		}
	}

	vkGroupSites.Groups = append(vkGroupSites.Groups, site.Site)
	tgSites.AddSite(&site)
	return nil
}

func deleteVkSite(domain string) {
	vkGroupSites.Mux.Lock()
	defer vkGroupSites.Mux.Unlock()
	var index int
	var groups []*news.Site
	for _, site := range vkGroupSites.Groups {
		if site.OptionalURL == domain {
			tgSites.DelSite(site.URL)
			break
		}
		groups = append(groups, site)
		index++
	}

	if index+1 < len(vkGroupSites.Groups) {
		for _, site := range vkGroupSites.Groups[index+1:] {
			site.ID--
			groups = append(groups, site)
		}
	}

	vkGroupSites.Groups = groups
}

func vkSiteMenuCommand(user *core.User, command *core.Command) {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить VK группу", addCommand(strCmdAddVKSite, "")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить VK группу", addCommand(strCmdDelVKSite, "")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButtonText, addCommand(strCmdAdminMenu, "")),
		),
	)

	sendMessage(user, command, "Управление VK группами", &markup)
}

func addVkSiteCommand(user *core.User, command *core.Command) {
	if len(command.ArgsStr) < 2 {
		user.ContinuationCommand = true
		user.CurrentCommand = command
		command.FieldNames = []string{"domain", "title"}
		tgBot.Send(tgbotapi.NewMessage(user.ID, "Введите домен и название"))
		return
	}

	if err := addVkSite(command.ArgsStr[0], strings.Join(command.ArgsStr[1:], " ")); err != nil {
		tgBot.Send(tgbotapi.NewMessage(user.ID, err.Error()))
	} else {
		tgBot.Send(tgbotapi.NewMessage(user.ID, "Готово"))
	}
}

func deleteVkSiteCommand(user *core.User, command *core.Command) {
	if len(command.ArgsStr) == 0 {
		user.ContinuationCommand = true
		user.CurrentCommand = command
		command.FieldNames = []string{"domain"}
		sendError(user, command, "Введите домен", false)
		return
	}

	deleteVkSite(command.ArgsStr[0])

	tgBot.Send(tgbotapi.NewMessage(user.ID, "Готово"))
}
