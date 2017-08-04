package core

import (
	"sync"
)

const (
	PlatformTg = "tg"
	PlatformVk = "vk"
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
	CurrentCommand      Command    `json:"command"`
	QueueMux            sync.Mutex `json:"-"`
	Queue               []Command  `json:"queue"`
}

type Key struct {
	Platform string `json:"platform"`
	ID       int64  `json:"id"`
}

type Users struct {
	Mux   sync.RWMutex  `json:"-"`
	Users map[Key]*User `json:"users"`
}

func (u *Users) SetUser(prefix string, user *User) {
	u.Mux.Lock()
	u.Users[Key{
		Platform: prefix,
		ID:       user.ID,
	}] = user
	u.Mux.Unlock()
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
		if key.Platform == platform {
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
