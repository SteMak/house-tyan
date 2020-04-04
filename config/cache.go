package config

import "time"

type ttl struct {
	Blank    time.Duration `yaml:"blank,omitempty"`
	Username time.Duration `yaml:"username,omitempty"`
}

type cache struct {
	Path string `yaml:"path,omitempty"`
	TTL  ttl    `yaml:"ttl,omitempty"`
}
