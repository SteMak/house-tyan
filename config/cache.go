package config

import "github.com/SteMak/house-tyan/config/parsers/json"

type ttl struct {
	Username json.Duration `json:"username,omitempty"`
}

type cache struct {
	Path string `json:"path,omitempty"`
	TTL  ttl    `json:"ttl,omitempty"`
}
