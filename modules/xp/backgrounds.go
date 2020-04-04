package xp

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type voiceXpWorker struct {
	config  *config
	session *discordgo.Session

	stopCh chan struct{}
}

func newVoiceXpWorker(c *config, s *discordgo.Session) *voiceXpWorker {
	return &voiceXpWorker{
		config:  c,
		session: s,
		stopCh:  make(chan struct{}),
	}
}

func (w *voiceXpWorker) onTick() {
	// Тут должна быть логика начисления опыта
}

func (w *voiceXpWorker) start() {
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
}

func (w *voiceXpWorker) stop() {
	w.stopCh <- struct{}{}
}
