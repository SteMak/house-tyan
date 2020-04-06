package awards

type bank struct {
	Token string `yaml:"-"`
}

type channels struct {
	Confirm  string `yaml:"confirm,omitempty"`
	Requests string `yaml:"requests,omitempty"`
	Bump     string `yaml:"bump,omitempty"`
}

type users struct {
	Confirmator string `yaml:"confirmator,omitempty"`
}

type roles struct {
	Requester string `yaml:"requester,omitempty"`
}

type bots struct {
	Bumper string `yaml:"bumper,omitempty"`
	Uper   string `yaml:"uper,omitempty"`
}

type config struct {
	AwardAmount int `yaml:"award_amount,omitempty"`

	Channels channels `yaml:"channels,omitempty"`
	Users    users    `yaml:"users,omitempty"`
	Roles    roles    `yaml:"roles,omitempty"`
	Bots     bots     `yaml:"bots,omitempty"`

	Bank bank `yaml:"bank,omitempty"`
}
