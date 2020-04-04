package xp

import (
	"time"

	"github.com/SteMak/house-tyan/out"

	conf "github.com/SteMak/house-tyan/config"
	"github.com/pkg/errors"

	"github.com/bwmarrin/discordgo"
)

type voiceXpWorker struct {
	config      *config
	session     *discordgo.Session
	voiceStates *[]*discordgo.VoiceState

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
	if len(*w.voiceStates) == 0 {
		return
	}

	channels := make(map[string][]*discordgo.VoiceState)
	for _, state := range *w.voiceStates {
		channels[state.ChannelID] = append(channels[state.ChannelID], state)
	}

	for channelID, states := range channels {
		out.Debugln("\nChannelID:", channelID)
		for _, st := range states {
			out.Debugf("%+v\n", *st)
		}
	}
}

func (w *voiceXpWorker) start(s *discordgo.Session) error {
	if w.isRunning {
		return nil
	}
	w.session = s

	guild, err := w.session.State.Guild(conf.Bot.GuildID)
	if err != nil {
		return errors.WithStack(err)
	}

	w.voiceStates = &guild.VoiceStates
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
