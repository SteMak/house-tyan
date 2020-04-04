package xp

import (
	"time"
)

type messageFarm struct {
	Channels      []string `yaml:"channels,omitempty"`
	XpForMessage  int      `yaml:"xp_for_message,omitempty"`
	MessageLength int      `yaml:"message_length,omitempty"`
}

type voiceFarm struct {
	WaitFor    time.Duration `yaml:"wait_for,omitempty"`
	XpForVoice int           `yaml:"xp_for_voice,omitempty"`
}

type config struct {
	RoleHermit  string      `yaml:"role_hermit,omitempty"`
	MessageFarm messageFarm `yaml:"message_farm,omitempty"`
	VoiceFarm   voiceFarm   `yaml:"voice_farm,omitempty"`
}
