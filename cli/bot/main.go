package main

import (
	"os"

	_ "github.com/SteMak/house-tyan/modules/awards"
	_ "github.com/SteMak/house-tyan/modules/xp"

	"github.com/SteMak/house-tyan/out"
)

func main() {
	if err := commands().Run(os.Args); err != nil {
		out.Fatal(err)
	}
}
