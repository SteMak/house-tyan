package xp

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

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

	content, _ = thinkAboutMathing(content, `\s\s+`, " ", 0)
	content, _ = thinkAboutMathing(content, `\*\*\*(\*)+`, "*", 0)
	content, countOfRunes = thinkAboutMathing(content, `<@!?\d+>`, "", countOfRunes)
	content, countOfRunes = thinkAboutMathing(content, `<@&\d+>`, "", countOfRunes)
	content, countOfRunes = thinkAboutMathing(content, `(https?://)?(\w+\.)+\S*`, "", countOfRunes)
	content, countOfRunes = thinkAboutMathing(content, `<a?:\w+:\d+>`, "", countOfRunes)

	content, _ = thinkAboutMathing(content, `\S{15,}`, "123456789012345", 0)

	countOfCommon, countOfRunes = countOtherSymbols(content, countOfCommon, countOfRunes)

	return countOfCommon, countOfRunes
}

func thinkAboutMathing(content, pattern string, replace string, acc int) (string, int) {
	rx := regexp.MustCompile(pattern)
	foundings := rx.FindAllString(content, -1)

	return rx.ReplaceAllLiteralString(content, replace), acc + len(foundings)
}

func countOtherSymbols(content string, common, runes int) (int, int) {
	contentRunes := []rune(content)
	for _, r := range contentRunes {
		if r == rune(65039) {
			continue
		}
		if r == rune(8205) {
			runes--
			continue
		}
		if isCommon(r) {
			common++
		} else {
			runes++
		}
	}

	return common, runes
}

func isCommon(r rune) bool {
	return !unicode.IsGraphic(r) ||
		unicode.IsPunct(r) ||
		unicode.IsMark(r) ||
		unicode.IsLower(r) ||
		unicode.IsLetter(r) ||
		unicode.IsNumber(r) ||
		unicode.IsSpace(r) ||
		util.EqualAny(string(r), strings.Split(" !\"#$%&'()*+,-./:;<=>?@[]^_`{|}", ""))
}
