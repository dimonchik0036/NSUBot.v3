package tgbot

import (
	"github.com/dimonchik0036/nsu-bot/core"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type BotNews struct {
	Date int64
	Text string
}

func (b *BotNews) String() string {
	return "Дата изменения: " + time.Unix(b.Date, 0).Format("15:04 02.01.2006") + "\n\n" + b.Text
}

type ByDate []BotNews

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date > a[j].Date }

const (
	tgBotNewsFilename   = "tgbotnews.txt"
	strCmdAddBotNews    = "addbotnews"
	strCmdReloadBotNews = "reloadbotnews"
	strCmdBotNewsMenu   = "botnewsmenu"
	strCmdBotNewsList   = "bnlist"
)

var tgBotNews struct {
	Mux  sync.RWMutex
	News []BotNews
}

func initBotNewsCommand() {
	tgCommands.AddHandler(core.Handler{
		Handler:         reloadBotNewsCommand,
		PermissionLevel: core.PermissionAdmin,
	}, strCmdReloadBotNews)

	tgCommands.AddHandler(core.Handler{
		Handler:         addBotNewsCommand,
		PermissionLevel: core.PermissionAdmin,
	}, strCmdAddBotNews)

	tgCommands.AddHandler(core.Handler{
		Handler:         botNewsMenuCommand,
		PermissionLevel: core.PermissionAdmin,
	}, strCmdBotNewsMenu)

	tgCommands.AddHandler(core.Handler{Handler: botNewsListCommand}, strCmdBotNewsList)
}

func initBotNews() {
	botNews, err := loadBotNews(tgBotNewsFilename)
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}

	tgBotNews.News = botNews
}

func loadBotNews(filename string) ([]BotNews, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Loading %s is failed. Err: %s", filename, err.Error())
		return []BotNews{}, err
	}

	var botNews []BotNews
	if err := yaml.Unmarshal(data, &botNews); err != nil {
		log.Printf("Loading %s is failed. Err: %s", filename, err.Error())
		return []BotNews{}, err
	}

	sort.Sort(ByDate(botNews))

	return botNews, nil
}

func saveBotNews(filename string, botNews []BotNews) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.FileMode(0700))
	if err != nil {
		log.Printf("Saving %s is failed. Err: %s", filename, err.Error())
		return err
	}

	data, err := yaml.Marshal(botNews)
	if err != nil {
		log.Printf("Saving %s is failed. Err: %s", filename, err.Error())
		return err
	}

	if _, err := file.Write(data); err != nil {
		log.Printf("Saving %s is failed. Err: %s", filename, err.Error())
		return err
	}

	return file.Close()
}

func appendBotNews(botNews []BotNews, newNews BotNews) []BotNews {
	botNews = append(botNews, newNews)
	sort.Sort(ByDate(botNews))
	if err := saveBotNews(tgBotNewsFilename, botNews); err != nil {
		log.Printf("%s", err.Error())
	}

	return botNews
}

func addBotNewsCommand(user *core.User, command *core.Command) {
	if len(command.ArgsStr) == 0 {
		user.ContinuationCommand = true
		user.CurrentCommand = command
		command.FieldNames = []string{"text"}
		sendError(user, command, "Наберите текст обновления", false)
		return
	}

	tgBotNews.Mux.Lock()
	defer tgBotNews.Mux.Unlock()
	tgBotNews.News = appendBotNews(tgBotNews.News, BotNews{
		Date: time.Now().Unix(),
		Text: strings.Join(command.ArgsStr, " "),
	})

	tgBot.Send(tgbotapi.NewMessage(user.ID, "Готово"))
}

func reloadBotNewsCommand(user *core.User, command *core.Command) {
	botNews, err := loadBotNews(tgBotNewsFilename)
	if err != nil {
		sendError(user, command, err.Error(), true)
		return
	}

	tgBotNews.Mux.Lock()
	defer tgBotNews.Mux.Unlock()
	tgBotNews.News = botNews
	sendError(user, command, "Готово", false)
}

func botNewsMenuCommand(user *core.User, command *core.Command) {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Перезагрузка новостей", addCommand(strCmdReloadBotNews, "")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить новость", addCommand(strCmdAddBotNews, "")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButtonText, addCommand(strCmdAdminMenu, "")),
		),
	)

	sendMessage(user, command, "Управление новостями", &markup)
}

func botNewsListCommand(user *core.User, command *core.Command) {
	tgBotNews.Mux.RLock()
	defer tgBotNews.Mux.RUnlock()
	var pageNumber int
	args := command.GetArg(strCmdArg)
	if len(args) != 0 {
		pageNumber, _ = strconv.Atoi(args)
	}

	if pageNumber >= len(tgBotNews.News) {
		pageNumber = len(tgBotNews.News) - 1
	}

	if pageNumber < 0 {
		pageNumber = 0
	}

	var text string = "Обновление: " + strconv.Itoa(pageNumber+1) + "/" + strconv.Itoa(len(tgBotNews.News)) + "\n"
	if len(tgBotNews.News) == 0 {
		text = "Нет никаких новостей."
	} else {
		text += tgBotNews.News[pageNumber].String()
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("«", addCommand(strCmdBotNewsList, strconv.Itoa(pageNumber-1))),
		tgbotapi.NewInlineKeyboardButtonData("Назад", addCommand(strCmdOptionMenu, "")),
		tgbotapi.NewInlineKeyboardButtonData("»", addCommand(strCmdBotNewsList, strconv.Itoa(pageNumber+1))),
	))

	sendMessage(user, command, text, &markup)
}
