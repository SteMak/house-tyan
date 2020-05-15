package clubs

import (
	"time"
)

type config struct {
	Cost                uint          `yaml:"cost,omitempty"`
	MinimumMembers      uint          `yaml:"minimum_members,omitempty"`
	NotVerifiedLifetime time.Duration `yaml:"not_verified_lifetime,omitempty"`
}
