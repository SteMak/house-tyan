package main

import (
	"os"

	"github.com/SteMak/house-tyan/out"
)

func main() {
	if err := commands().Run(os.Args); err != nil {
		out.Fatal(err)
	}
}
