package logic

import (
	"fmt"

	"github.com/SteMak/ani-helper/workerTools/config"
	"github.com/bwmarrin/discordgo"
)

// Init initialize discord session
func Init() {
	fmt.Println("1 WORKER started")

	session, err := discordgo.New(config.Token)
	if err != nil {
		fmt.Println("ERROR creating Discord session:", err)
		return
	}

	fmt.Println("2 Discord session created")

	session.AddHandler(messageCreate)
	session.AddHandler(reactionHandler)

	fmt.Println("3 Registred the handlers")

	err = session.Open()
	if err != nil {
		fmt.Println("ERROR opening connection:", err)
		return
	}

	fmt.Println("4 Opened a websocket")
}
