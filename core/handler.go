package core

type Handler struct {
	PermissionLevel int
	Handler         func(user *User, command *Command)
}

type Handlers map[string]Handler

func (h Handlers) AddHandler(handler Handler, key ...string) {
	for _, v := range key {
		h[v] = handler
	}
}

func CommandHandler(user *User, command *Command, commandInterpreter func(string) string, handlers Handlers) {
	handler, ok := handlers[commandInterpreter(command.Command)]
	if !ok || !checkHandler(user, handler) {
		return
	}

	handler.Handler(user, command)
}

func checkHandler(user *User, handler Handler) bool {
	return user.Permission >= handler.PermissionLevel
}
