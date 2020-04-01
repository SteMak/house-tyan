package config

type module struct {
	Enabled bool   `json:"enabled,omitempty"`
	Prefix  string `json:"prefix,omitempty"`
	Config  string `json:"config,omitempty"`
}
