package awards

type bank struct {
	Token string
}

type channels struct {
	Confirm   string `json:"confirm,omitempty"`
	Responces string `json:"responces,omitempty"`
	Logs      string `json:"logs,omitempty"`
	Bump      string `json:"bump,omitempty"`
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
