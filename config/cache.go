package config

import "time"

type ttl struct {
	Username time.Duration `yaml:"username,omitempty"`
}

type cache struct {
	Path string `yaml:"path,omitempty"`
	TTL  ttl    `yaml:"ttl,omitempty"`
}
