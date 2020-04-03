package config

type storage struct {
	Driver     string `yaml:"driver,omitempty"`
	Connection string `yaml:"connection,omitempty"`
}
