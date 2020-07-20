package app

import (
	"strings"
)

func matchCommand(alias, args string) (bool, string) {
	cut := args
	for len(cut) != 0 {
		if strings.HasPrefix(cut, alias) {
			return true, strings.TrimPrefix(cut, alias)
		}
		cut = cut[1:]
	}
	return false, args
}
