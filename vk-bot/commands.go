package vkbot

import "github.com/dimonchik0036/nsu-bot/core"

var Commands core.Handlers

func initCommands() {
	Commands = core.Handlers{}
	Commands.AddHandler(weatherHandler(), "w", "weather", "погода", "температура")
	Commands.AddHandler(helpHandler(), "\"help\"", "help", "помощь")
	Commands.AddHandler(welcomeHandler(), "\"привет\"", "привет", "хай")
}

func welcomeHandler() core.Handler {
	return core.Handler{
		Handler: WelcomeCommand,
	}
}

func WelcomeCommand(user *core.User, command core.Command) {
	Bot.SendMessage(user.ID, "Привет, "+user.FirstName+"!\n"+
		`Я - Помощник, если ты хочешь посмотреть доступные команды, то напиши мне слово "help".

		Если хочешь узнать что-то дополнительное обо мне, то пиши @dimonchik0036(ему).`)
}

func weatherHandler() core.Handler {
	return core.Handler{
		Handler: WeatherCommand,
	}
}

func WeatherCommand(user *core.User, command core.Command) {
	Bot.SendMessage(user.ID, Weather.ShowWeather()+"\n"+
		"Время последнего обновления: "+Weather.ShowTime())
}

func helpHandler() core.Handler {
	return core.Handler{
		PermissionLevel: 0,
		Handler:         HelpCommand,
	}
}

func HelpCommand(user *core.User, command core.Command) {
	defaultHelp := `Список команд:
	w | weather | погода | температура - Показывает текущую температуру около НГУ.

	По всем вопросам можно обратиться к @dimonchik0036(нему).`
	Bot.SendMessage(user.ID, defaultHelp)
}
