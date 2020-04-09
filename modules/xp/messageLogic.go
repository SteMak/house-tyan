package xp

import (
	"errors"
	"regexp"

	"github.com/SteMak/house-tyan/util"

	"github.com/bwmarrin/discordgo"
)

func xpMessageChecks(
	channelID string,
	conf config,
	guildID string,
	authorID string,
	getMember func(string, string) (*discordgo.Member, error),
) error {
	if !util.EqualAny(channelID, conf.MessageFarm.Channels) {
		return errors.New("isn't a message fam channel")
	}

	member, err := getMember(guildID, authorID)
	if err != nil || member.User.Bot {
		return errors.New("maybe bot")
	}
	if util.EqualAny(conf.RoleHermit, member.Roles) {
		return errors.New("it's a hermit")
	}

	return nil
}

func howMuchXp(content string, messageFarm messageFarm) int {
	if len(content) == 0 {
		return int(messageFarm.XpForEmpty)
	}

	countOfCommon, countOfRunes := countSymbols(content)

	charXP := messageFarm.XpForChar * float32(countOfCommon)
	runeXP := messageFarm.XpForRune * float32(countOfRunes)

	return int(charXP) + int(runeXP)
}

func countSymbols(content string) (int, int) {
	var (
		countOfRunes  int
		countOfCommon int
	)

	content, countOfRunes = thinkAboutMathing(content, `<@!?\d+>`, countOfRunes)
	content, countOfRunes = thinkAboutMathing(content, `<a?:\w+:\d+>`, countOfRunes)

	countOfCommon, countOfRunes = countOtherSymbols(content, countOfCommon, countOfRunes)

	return countOfCommon, countOfRunes
}

func thinkAboutMathing(content, pattern string, acc int) (string, int) {
	rx := regexp.MustCompile(pattern)
	foundings := rx.FindAllString(content, -1)

	return rx.ReplaceAllLiteralString(content, ""), acc + len(foundings)
}

func countOtherSymbols(content string, common, runes int) (int, int) {
	contentRunes := []rune(content)
	for _, r := range contentRunes {
		if len(string(r)) != 1 {
			runes++
		} else {
			common++
		}
	}

	return common, runes
}
