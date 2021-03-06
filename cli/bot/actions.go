package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/modules"
	"github.com/SteMak/house-tyan/out"
	"github.com/SteMak/house-tyan/storage"
	"github.com/urfave/cli"
)

func run(c *cli.Context) error {
	fmt.Println("Bot is running. Press Ctrl + C to exit.")

	config.Load(c.GlobalString("config"))
	out.SetDebug(c.GlobalBool("debug"))

	storage.Init()

	modules.Run()
	defer modules.Stop()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	return nil
}

func migrate(c *cli.Context) error {
	config.Load(c.GlobalString("config"))
	out.SetDebug(c.GlobalBool("debug"))

	storage.Init()

	return nil
}
