package voices

import "github.com/SteMak/house-tyan/modules"

var _module module

func init() {
	modules.Register(_module.ID(), &_module)
}
