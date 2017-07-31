package gconfig

import "github.com/dimonchik0036/nsu-bot/schedule"

type Config struct {
	Schedule *schedule.Schedule
	Weather  *string
}
