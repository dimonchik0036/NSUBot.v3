package gconfig

import (
	"github.com/dimonchik0036/nsu-bot/nsuschedule"
)

type Config struct {
	Schedule *nsuschedule.Schedule
	Weather  *string
}
