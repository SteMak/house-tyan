package tyan

import "github.com/urfave/cli"

var (
	// Vesion текущяя версия (устанавливается через -ldflags "-X vanilla.Version=")
	Vesion = "debug"

	GlobalCtx *cli.Context
)
