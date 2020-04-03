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
	Confirmator string `json:"confirmator,omitempty"`
}

type roles struct {
	Requester string `json:"requester,omitempty"`
}

type bots struct {
	Bumper string `json:"bumper,omitempty"`
	Uper   string `json:"uper,omitempty"`
}

type config struct {
	AwardAmount int `json:"award_amount,omitempty"`

	Channels channels `json:"channels,omitempty"`
	Users    users    `json:"users,omitempty"`
	Roles    roles    `json:"roles,omitempty"`
	Bots     bots     `json:"bots,omitempty"`

	Bank bank `json:"bank,omitempty"`
}
