package awards

import (
	"strconv"
	"strings"

	"github.com/SteMak/house-tyan/cache"
	conf "github.com/SteMak/house-tyan/config"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

func parseRequest(s *discordgo.Session, content string) (*cache.Award, error) {
	args := strings.Split(strings.TrimPrefix(content, "-запрос "), "->")
	if len(args) > 2 {
		return nil, errors.New("Обнаружена лишняя \"->\"")
	}
	if len(args) < 2 {
		return nil, errors.New("Где \"->\"?")
	}

	reason := strings.TrimSpace(args[0])
	if len(reason) == 0 {
		return nil, errors.New("Причина не должна быть пустой")
	}

	item := &cache.Award{
		Reason: reason,
	}

	if len(args[1]) == 0 {
		return nil, errors.New("Не указаны юзвери и их деньги")
	}

	args[1] = makeParsingBetter(strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(args[1]), ","), ","))

	pairsUsersSum := strings.Split(args[1], ",")
	for _, pairUsersSum := range pairsUsersSum {
		if len(pairUsersSum) == 0 {
			return nil, errors.New("Не указаны юзвери и их деньги")
		}

		usersSum, err := splitReverse1(strings.TrimSpace(pairUsersSum), " ")

		if len(usersSum[0]) == 0 {
			return nil, errors.New("В указании юзверей или суммы содержится ошибка")
		}

		sum, err := strconv.ParseUint(strings.TrimSpace(usersSum[1]), 10, 64)
		if err != nil {
			return nil, errors.New("В указании юзверей или суммы содержится ошибка")
		}

		users := strings.Split(usersSum[0], " ")

		for i := 0; i < len(users); i++ {
			if !strings.HasPrefix(users[i], "<@") || !strings.HasSuffix(users[i], ">") {
				return nil, errors.New("В юзверях затесался шпион")
			}

			users[i] = strings.TrimPrefix(users[i], "<@")
			users[i] = strings.TrimPrefix(users[i], "!")
			users[i] = strings.TrimSuffix(users[i], ">")

			_, err = s.GuildMember(conf.Bot.GuildID, users[i])
			if err != nil {
				return nil, errors.New("В юзверях затесался шпион")
			}
		}

		for _, id := range users {
			item.Users = append(item.Users, cache.User{
				ID:     id,
				Amount: sum,
			})
		}
	}
	return item, nil
}

func makeParsingBetter(str string) string {
	var FairyReplacement = [][2]string{
		[2]string{"  ", " "},
		[2]string{",,", ","},
		[2]string{", ,", ","},
		[2]string{"><", "> <"},
		[2]string{">1", "> 1"},
		[2]string{">2", "> 2"},
		[2]string{">3", "> 3"},
		[2]string{">4", "> 4"},
		[2]string{">5", "> 5"},
		[2]string{">6", "> 6"},
		[2]string{">7", "> 7"},
		[2]string{">8", "> 8"},
		[2]string{">9", "> 9"},
		[2]string{">0", "> 0"},
	}
	result := str
	for _, rep := range FairyReplacement {
		result = strings.Join(strings.Split(result, rep[0]), rep[1])
	}

	if result != str {
		return makeParsingBetter(result)
	}

	return result
}

func splitReverse1(str string, sep string) ([2]string, error) {
	for i := len(str) - 1; i > 0; i-- {
		if string(str[i]) == sep {
			return [2]string{str[:i], str[i+1:]}, nil
		}
	}

	return [2]string{"", str}, errors.New("")
}
