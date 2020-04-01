package config

type storage struct {
	Driver     string `json:"driver,omitempty"`
	Connection string `json:"connection,omitempty"`
}
