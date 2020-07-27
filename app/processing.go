package app

import (
	"strings"
)

func matchCommand(alias, args string) (bool, string) {
	cut := args
	for len(cut) != 0 {
		if strings.HasPrefix(cut, alias) {
			return true, strings.TrimSpace(strings.TrimPrefix(cut, alias))
		}
		if cut[0] != ' ' {
			return false, args
		}
		cut = cut[1:]
	}
	return false, args
}
