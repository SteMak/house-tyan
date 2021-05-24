package main

import (
	"os"

	_ "github.com/SteMak/house-tyan/modules/awards"
	_ "github.com/SteMak/house-tyan/modules/triggers"
	_ "github.com/SteMak/house-tyan/modules/xp"
	_ "github.com/SteMak/house-tyan/modules/voices"

	"github.com/SteMak/house-tyan/out"
)

func main() {
	if err := commands().Run(os.Args); err != nil {
		out.Fatal(err)
	}
}
