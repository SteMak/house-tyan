package libs

import (
	"os"

	"github.com/SteMak/house-tyan/libs/unb"
)

var (
	Unb *unb.UnbelievaBoatAPI
)

func Init() {
	Unb = unb.NewUnbelievaBoatAPI(os.Getenv("BANKER_TOKEN"))
}
