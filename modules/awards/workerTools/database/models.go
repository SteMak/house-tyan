package database

// Us additional structure for Record
type Us struct {
	Users []string `json:"users"`
	Sum   uint64   `json:"sum"`
}

// Uss structure for users and sums
type Uss struct {
	Us []Us `json:"us"`
}

// Record for database
type Record struct {
	EmbedID    string `gorm:"pk"`
	AuthorName string
	AuthorIcon string
	Reason     string
	UsersSum   string
}
