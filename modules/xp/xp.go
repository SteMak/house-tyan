package xp

func (bot *module) handleXp(userID string, xp uint64) {
	for _, fn := range bot.xpHandlers {
		go fn(userID, xp)
	}
}
