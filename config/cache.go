package config

import "github.com/SteMak/house-tyan/config/parsers/json"

type ttl struct {
	Blank    time.Duration `yaml:"blank,omitempty"`
	Username time.Duration `yaml:"username,omitempty"`
}

type cache struct {
	Path string `json:"path,omitempty"`
	TTL  ttl    `json:"ttl,omitempty"`
}
