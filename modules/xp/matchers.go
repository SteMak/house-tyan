package xp

func isntRightChannel(realChannelID string, normalChannels []string) bool {
	for _, normalChannelID := range normalChannels {
		if normalChannelID == realChannelID {
			return false
		}
	}

	return true
}
