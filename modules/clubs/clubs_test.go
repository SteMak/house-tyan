package clubs

import (
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/SteMak/house-tyan/cache"
	conf "github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/messages"
	"github.com/SteMak/house-tyan/modules"
	"github.com/SteMak/house-tyan/storage"
)

func TestMain(t *testing.M) {
	conf.Load("./../../cli/bot/config/dev/config.yaml")

	cache.Init()
	defer cache.Close()

	storage.Init()

	conf.Bot.Templates = "./../../cli/bot/assets/templates/"
	messages.Init()

	modules.Run()
	defer modules.Stop()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
