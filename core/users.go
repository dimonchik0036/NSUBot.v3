package core

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	PlatformTg = "tg"
	PlatformVk = "vk"
)

const (
	UserLayout = "2006/01/02 15:04:05"
)

type User struct {
	ID                  int64      `json:"id"`
	Username            string     `json:"username"`
	FirstName           string     `json:"first_name"`
	LastName            string     `json:"last_name"`
	Platform            string     `json:"platform"`
	Permission          int        `json:"permission"`
	DateCreated         int64      `json:"date_created"`
	DateLastActivities  int64      `json:"date_last_activities"`
	ContinuationCommand bool       `json:"command_in_queue"`
	CurrentCommand      *Command   `json:"command"`
	QueueMux            sync.Mutex `json:"-"`
}

func (u *User) String() string {
	return strconv.FormatInt(u.ID, 10) + ", " + u.Username + ", " + u.FirstName + " " + u.LastName
}

func (u *User) NewUserString(usernamePrefix string) string {
	return fmt.Sprintf("ID:%d\n"+
		"Ник: %s\n"+
		"Имя: %s\n"+
		"Фамилия: %s\n"+
		"Дата регистрации: %s", u.ID, usernamePrefix+u.Username, u.FirstName, u.LastName, time.Unix(u.DateCreated, 0).Format(UserLayout))
}

func key(prefix string, id int64) string {
	return prefix + strconv.FormatInt(id, 10)
}

type Users struct {
	Mux   sync.RWMutex     `json:"-"`
	Users map[string]*User `json:"users"`
}

func (u *Users) DelUser(prefix string, id int64) {
	u.Mux.Lock()
	defer u.Mux.Unlock()
	delete(u.Users, key(prefix, id))
}

func (u *Users) User(prefix string, id int64) *User {
	u.Mux.RLock()
	defer u.Mux.RUnlock()
	return u.Users[key(prefix, id)]
}

func (u *Users) VkUser(id int64) *User {
	return u.User(PlatformVk, id)
}

func (u *Users) TgUser(id int64) *User {
	return u.User(PlatformTg, id)
}

func (u *Users) SetUser(prefix string, user *User) {
	u.Mux.Lock()
	defer u.Mux.Unlock()
	if u.Users == nil {
		u.Users = map[string]*User{}
	}

	u.Users[key(prefix, user.ID)] = user
}

func (u *Users) SetVkUser(user *User) {
	u.SetUser(PlatformVk, user)
}

func (u *Users) SetTgUser(user *User) {
	u.SetUser(PlatformTg, user)
}

func (u *Users) PlatformUsers(platform string) (result []*User) {
	u.Mux.RLock()
	defer u.Mux.RUnlock()
	for key, u := range u.Users {
		if strings.Contains(key, platform) {
			result = append(result, u)
		}
	}

	return
}

func (u *Users) TgUsers() []*User {
	return u.PlatformUsers(PlatformTg)
}

func (u *Users) VkUsers() []*User {
	return u.PlatformUsers(PlatformVk)
}
