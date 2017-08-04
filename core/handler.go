package core

type Handler struct {
	PermissionLevel int
	Handler         func(user *User, command Command)
}

func CommandHandler(user *User, command Command, handlers map[string]Handler) {
	handler, ok := handlers[command.Command]
	if !ok || !checkHandler(user, handler) {
		return
	}

	handler.Handler(user, command)
}

func checkHandler(user *User, handler Handler) bool {
	return user.Permission >= handler.PermissionLevel
}
