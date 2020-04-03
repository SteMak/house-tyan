package config

type bot struct {
	GuildID        string  `yaml:"guild_id,omitempty"`
	LogChannel     *string `yaml:"log_channel,omitempty"`
	ErrorsChannel  *string `yaml:"errors_channel,omitempty"`
	ConsoleChannel *string `yaml:"console_channel,omitempty"`
	Templates      string  `yaml:"templates,omitempty"`
}
