package awards

import (
	"os"
)

func (bot *module) loadEnv() {
	bot.config.Bank.Token = os.Getenv("BANKER_TOKEN")
}
