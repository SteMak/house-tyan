package app

import (
	"errors"
	"strings"
)

type scannable interface {
	Scan() error
}

type scanner struct {
	reader strings.Reader
}

func (s *scanner) Scan(pattern string, dest ...interface{}) error {
	if strings.Count(pattern, "?") != len(dest) {
		return errors.New("")
	}
	return nil
}
