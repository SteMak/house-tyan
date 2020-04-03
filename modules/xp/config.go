package xp

type channels struct {
	XpFarming []string `yaml:"xp_farming,omitempty"`
}

type xpConfig struct {
	XpBoost   int `yaml:"xp_boost,omitempty"`
	XpMessage int `yaml:"xp_message,omitempty"`
	XpVoice   int `yaml:"xp_voice,omitempty"`
	XpMesLen  int `yaml:"xp_mes_len,omitempty"`
	XpVoiSec  int `yaml:"xp_voi_sec,omitempty"`
}

type roles struct {
	Hermit string `yaml:"hermit,omitempty"`
}

type config struct {
	Channels channels `yaml:"channels,omitempty"`
	Roles    roles    `yaml:"roles,omitempty"`
	XpConfig xpConfig `yaml:"xp_config,omitempty"`
}
