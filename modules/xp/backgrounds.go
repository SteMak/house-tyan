package xp

import (
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
	if len(*w.voiceStates) < 2 {
		return
	}

	channels := make(map[string][]*discordgo.VoiceState)
	for _, state := range *w.voiceStates {
		channels[state.ChannelID] = append(channels[state.ChannelID], state)
	}

	out.Debugln()

	for channelID, allStates := range channels {
		guild, err := w.session.State.Guild(conf.Bot.GuildID)
		if err != nil || channelID == guild.AfkChannelID {
			continue
		}

		var voiceMembers []*discordgo.Member
		for _, st := range allStates {
			member, err := w.session.State.Member(st.GuildID, st.UserID)
			if err != nil || member.User.Bot {
				out.Err(false, err)
				continue
			}
			if st.SelfDeaf || st.Deaf {
				continue
			}
			if st.SelfMute || st.Mute {
				continue
			}

			voiceMembers = append(voiceMembers, member)
		}

		if len(voiceMembers) < 2 {
			continue
		}

		roomBoost := 1
		if len(voiceMembers) > w.config.VoiceFarm.MaxRoomBoost {
			roomBoost = roomBoost * w.config.VoiceFarm.MaxRoomBoost
		} else {
			roomBoost = roomBoost * len(voiceMembers)
		}

		for _, member := range voiceMembers {
			if util.EqualAny(w.config.RoleHermit, member.Roles) {
				continue
			}

			out.Debugln(member.User.ID, w.config.VoiceFarm.XpForVoice*roomBoost)
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
