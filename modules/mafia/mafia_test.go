package mafia

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/dustin/go-humanize"

	conf "github.com/SteMak/house-tyan/config"
)

func TestMain(t *testing.M) {
	conf.Load("./../../cli/bot/config/dev/config.yaml")
	_module.LoadConfig("./../../cli/bot/config/dev/modules/mafia.yaml")

	os.Exit(t.Run())
}

func TestLoadImages(t *testing.T) {
	_module.loadImages()
	memstat := new(runtime.MemStats)
	runtime.ReadMemStats(memstat)
	fmt.Println(humanize.Bytes(memstat.TotalAlloc))
}
