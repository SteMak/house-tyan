package config

type bot struct {
	GuildID        string  `json:"guild_id,omitempty"`
	LogChannel     *string `json:"log_channel,omitempty"`
	ErrorsChannel  *string `json:"errors_channel,omitempty"`
	ConsoleChannel *string `json:"console_channel,omitempty"`
	Templates      string  `json:"templates,omitempty"`
}
