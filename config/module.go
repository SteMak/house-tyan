package config

type module struct {
	Enabled bool    `yaml:"enabled,omitempty"`
	Prefix  string  `yaml:"prefix,omitempty"`
	Config  *string `yaml:"config,omitempty"`
	Log     *Log    `yaml:"log,omitempty"`
}
