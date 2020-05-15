package modules

import "github.com/robfig/cron/v3"

func init() {
	Cron = cron.New()
}
