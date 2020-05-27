package mafia

import (
	"github.com/SteMak/house-tyan/modules"
	"github.com/sirupsen/logrus"
)

var (
	log     *logrus.Logger
	_module module
)

func init() {
	modules.Register(_module.ID(), &_module)
}
