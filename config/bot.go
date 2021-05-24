package config

type channels struct {
	Console  string `yaml:"console,omitempty"`
	Errors   string `yaml:"errors,omitempty"`
	Logs     string `yaml:"logs,omitempty"`
	Cards    string `yaml:"cards,omitempty"`
	Category string `yaml:"category,omitempty"`
}
type bot struct {
	GuildID   string   `yaml:"guild_id,omitempty"`
	Channels  channels `yaml:"channels,omitempty"`
	Templates string   `yaml:"templates,omitempty"`
}
