package xp

import "github.com/SteMak/house-tyan/modules"

var _module module

func init() {
	modules.Event.XpEvents = &_module
	modules.Register(_module.ID(), &_module)
}
