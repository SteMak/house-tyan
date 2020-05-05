package xp

import (
	"github.com/SteMak/house-tyan/modules"
)

func (bot *module) OnXp(fn modules.HandlerXP) {
	bot.xpHandlers = append(bot.xpHandlers, fn)
}
