package xp

import (
	"time"

	"github.com/SteMak/house-tyan/out"

	conf "github.com/SteMak/house-tyan/config"
	"github.com/pkg/errors"

	"github.com/bwmarrin/discordgo"
)

type voiceXpWorker struct {
	config *config
	state  *discordgo.State
	guild  *discordgo.Guild

	isRunning bool

	stopCh chan struct{}
}

func newVoiceXpWorker(c *config) *voiceXpWorker {
	return &voiceXpWorker{
		config:    c,
		isRunning: false,
		stopCh:    make(chan struct{}),
	}
}

func (w *voiceXpWorker) onTick() {
	if len(w.guild.VoiceStates) < 2 {
		return
	}

	out.Debugln()

	addXpUsers(
		w.guild.AfkChannelID,
		w.config,
		w.guild.VoiceStates,
		w.state.Member,
		func(userID string, xp int) {
			out.Debugln(userID, xp)
		},
	)
}

func (w *voiceXpWorker) start(s *discordgo.Session) error {
	if w.isRunning {
		return nil
	}
	w.state = s.State

	guild, err := w.state.Guild(conf.Bot.GuildID)
	if err != nil {
		return errors.WithStack(err)
	}

	w.guild = guild

	t := time.Tick(w.config.VoiceFarm.WaitFor)

	go func() {
		for {
			select {
			case <-t:
				go w.onTick()
			case <-w.stopCh:
				return
			}
		}
	}()
	w.isRunning = true
	return nil
}

func (w *voiceXpWorker) stop() {
	if !w.isRunning {
		return
	}
	w.stopCh <- struct{}{}
	w.isRunning = false
}
