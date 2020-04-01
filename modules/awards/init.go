package awards

import "github.com/SteMak/house-tyan/modules"

func init() {
	m := new(module)
	modules.Register(m.ID(), m)
}
