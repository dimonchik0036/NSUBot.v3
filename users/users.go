package users

type User struct {
}

type Users struct {
	VKUsers map[int64]*User
	TgUsers map[int64]*User
}
