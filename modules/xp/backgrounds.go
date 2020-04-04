package xp

import (
	"fmt"
	"time"

	"github.com/SteMak/house-tyan/out"
	"github.com/SteMak/house-tyan/util"

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

	// for channelID, states := range channels {
	// 	out.Debugln("\nChannelID:", channelID)
	// 	for _, st := range states {
	// 		out.Debugf("%+v\n", *st)
	// 	}
	// }

	out.Debugf("\n")

	for _, allStates := range channels {
		var states []*discordgo.VoiceState
		for _, st := range allStates {
			member, err := w.session.State.Member(st.GuildID, st.UserID)
			if err != nil || member.User.Bot {
				continue
			}
			if !st.SelfDeaf && !st.Deaf {
				states = append(states, st)
			}
		}

		if len(states) < 2 {
			break
		}

		roomBoost := 0
		for _, st := range states {
			if !st.SelfMute && !st.Mute {
				roomBoost++
			}
		}

		for _, st := range states {
			member, err := w.session.State.Member(st.GuildID, st.UserID)
			if err != nil {
				out.Err(false, err)
				continue
			}
			if st.SelfDeaf || st.Deaf {
				continue
			}
			if util.EqualAny(w.config.RoleHermit, member.Roles) {
				continue
			}

			fmt.Println(st.UserID, w.config.VoiceFarm.XpForVoice * roomBoost)
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
