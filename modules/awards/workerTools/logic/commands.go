package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/SteMak/ani-helper/workerTools/config"
	"github.com/SteMak/ani-helper/workerTools/database"
	"github.com/bwmarrin/discordgo"
)

var (
	err error
)

func detectBumpSiup(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == config.UsSiupID &&
		m.Embeds[0].Title == "Сервер Up" &&
		m.Embeds[0].Footer != nil {

		onSiupServer(s, m)
		return
	}

	matched, err := regexp.Match(`Server bumped by <@\d+>`, []byte(m.Embeds[0].Description))
	if err != nil {
		fmt.Println("ERROR Bump make match regular failure:", err)
		return
	}

	if matched && m.Author.ID == config.UsBumpID {
		onBumpServer(s, m)
		return
	}
}

func onSiupServer(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("FOUND S.up")

	for _, user := range chMonitorWriters {
		if user.strify == m.Embeds[0].Footer.Text {
			fmt.Println("FOUND S.up user", user.id)

			sendAndLog(s, user.id, "S.up", 1000)
			return
		}
	}
}

func onBumpServer(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("FOUND Bump")

	userID := strings.Split(strings.Split(m.Embeds[0].Description, "<@")[1], ">")[0]
	if len(userID) == 0 {
		fmt.Println("ERROR Bump get user ID:", m.Embeds[0].Description)
		return
	}

	if strings.HasPrefix(userID, "!") {
		userID = userID[1:]
	}

	fmt.Println("FOUND Bump user", userID)

	sendAndLog(s, userID, "Bump", 1000)
}

func sendAndLog(s *discordgo.Session, userID string, str string, sum int) {
	_, err = config.API.AddToBalance(config.GdHouseID, userID, 0, sum, "for "+str)
	if err != nil {
		fmt.Println("ERROR "+str+" updating user balance:", err)

		_, err = s.ChannelMessageSend(config.ChForLogsID, "Кажись, что-то пошло не так... <@"+userID+"> сделал "+str+", но денег ему не дали(")
		if err != nil {
			fmt.Println("ERROR "+str+" sending wrong report message:", err)
		}
		_, err = s.ChannelMessageSend(config.ChForBumpSiupID, "<@"+userID+">, у нас снова что-то сломалось, но не волнуйтесь - деньги Вам прилетят чуть позже)")
		if err != nil {
			fmt.Println("ERROR "+str+" sending wrong log message:", err)
		}

		return
	}

	_, err = s.ChannelMessageSend(config.ChForLogsID, strconv.Itoa(sum)+"<:AH_AniCoin:579712087224483850> были выданы <@"+userID+">, за то что он сделал "+str)
	if err != nil {
		fmt.Println("ERROR "+str+" sending right report message:", err)
	}

	_, err = s.ChannelMessageSend(config.ChForBumpSiupID, "<@"+userID+">, "+fmt.Sprintf(config.Responces[rand.Intn(len(config.Responces))], str, strconv.Itoa(sum)+"<:AH_AniCoin:579712087224483850>"))
	if err != nil {
		fmt.Println("ERROR "+str+" sending right log message:", err)
	}

	fmt.Println("GUILD "+str+" by  ", userID)
}

func splitReverse1(str string, sep string) ([2]string, error) {
	for i := len(str) - 1; i > 0; i-- {
		if string(str[i]) == sep {
			return [2]string{str[:i], str[i+1:]}, nil
		}
	}

	return [2]string{"", str}, errors.New("")
}

func createEmbed(s *discordgo.Session, item *database.Record, color int) *discordgo.MessageEmbed {
	var users string

	uss := database.Uss{}
	json.Unmarshal([]byte(item.UsersSum), &uss)
	for _, userSum := range uss.Us {
		users += strconv.FormatUint(userSum.Sum, 10) + "\n<@" + strings.Join(userSum.Users, "> <@") + ">\n"

		if color == 255255 {
			for _, user := range userSum.Users {
				_, err := config.API.AddToBalance(config.GdHouseID, user, 0, int(userSum.Sum), item.Reason)
				if err != nil {
					fmt.Println("ERROR "+item.EmbedID+" add money to user balance:", err)
					_, err = s.ChannelMessageSend(config.ChForLogsID, "Кажись, что-то пошло не так... <@"+user+"> не получил денег за "+item.Reason)
					if err != nil {
						fmt.Println("ERROR "+item.EmbedID+" sending wrong report message:", err)
					}
				} else {
					_, err = s.ChannelMessageSend(config.ChForLogsID, strconv.FormatUint(userSum.Sum, 10)+"<:AH_AniCoin:579712087224483850> были выданы <@"+user+">, за "+item.Reason)
					if err != nil {
						fmt.Println("ERROR "+item.EmbedID+" sending right report message:", err)
					}
				}
			}
		}
	}

	return &discordgo.MessageEmbed{
		Title: "Запрос на выдачу денег",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    item.AuthorName,
			IconURL: item.AuthorIcon,
		},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  "Запрос",
				Value: item.Reason,
			},
			&discordgo.MessageEmbedField{
				Name:  "Получатели",
				Value: users,
			},
		},
		Color: color,
	}
}

func hasRole(member *discordgo.Member, id string) bool {
	for _, role := range member.Roles {
		if role == id {
			return true
		}
	}
	return false
}

func detectRequest(s *discordgo.Session, channelID, content string) {
	if content == "-запрос" {
		s.ChannelMessageSend(channelID, config.Usage)
		return
	}

	if !strings.HasPrefix(content, "-запрос ") {
		return
	}

	item, err := queryIntoRecord(s, channelID, content)
	if err != nil {
		fmt.Println("BADQR parsing query:", err)
		return
	}

	err = database.Records.New(item)
	if err != nil {
		fmt.Println("ERROR create new database record", err)
		return
	}

	_, err = s.ChannelMessageSend(channelID, "Запрос отправлен")
	if err != nil {
		fmt.Println("ERROR sending confirm message", err)
	}

	fmt.Println("GUILD " + item.EmbedID + " request sended")
}

func queryIntoRecord(s *discordgo.Session, channelID, content string) (*database.Record, error) {
	item, err := parseQuery(s, channelID, content)
	if err != nil {
		return item, err
	}

	fmt.Println("FOUND request message", content)

	err = resendRequest(s, channelID, content, item)
	if err != nil {
		return item, err
	}

	return item, nil
}

func parseQuery(s *discordgo.Session, channelID, content string) (*database.Record, error) {
	resultErr := errors.New("parse failure")

	args := strings.Split(strings.TrimPrefix(content, "-запрос "), "->")
	if len(args) > 2 {
		s.ChannelMessageSend(channelID, "Обнаружена лишняя \"->\"")
		return nil, resultErr
	}
	if len(args) < 2 {
		s.ChannelMessageSend(channelID, "Где \"->\"?")
		return nil, resultErr
	}

	reason := strings.TrimSpace(args[0])
	if len(reason) == 0 {
		s.ChannelMessageSend(channelID, "Причина не должна быть пустой")
		return nil, resultErr
	}

	item := &database.Record{}
	item.Reason = reason

	if len(args[1]) == 0 {
		s.ChannelMessageSend(channelID, "Не указаны юзвери и их деньги")
		return nil, resultErr
	}

	args[1] = makeParsingBetter(strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(args[1]), ","), ","))

	var usersSumDB database.Uss
	pairsUsersSum := strings.Split(args[1], ",")
	for _, pairUsersSum := range pairsUsersSum {
		if len(pairUsersSum) == 0 {
			s.ChannelMessageSend(channelID, "Не указаны юзвери и их деньги")
			return nil, resultErr
		}

		usersSum, err := splitReverse1(strings.TrimSpace(pairUsersSum), " ")

		if len(usersSum[0]) == 0 {
			s.ChannelMessageSend(channelID, "В указании юзверей или суммы содержится ошибка")
			return nil, resultErr
		}

		sum, err := strconv.ParseUint(strings.TrimSpace(usersSum[1]), 10, 64)
		if err != nil {
			s.ChannelMessageSend(channelID, "В указании юзверей или суммы содержится ошибка")
			return nil, resultErr
		}

		users := strings.Split(usersSum[0], " ")

		for i := 0; i < len(users); i++ {
			if !strings.HasPrefix(users[i], "<@") || !strings.HasSuffix(users[i], ">") {
				s.ChannelMessageSend(channelID, "В юзверях затесался шпион")
				return nil, resultErr
			}
			users[i] = strings.TrimPrefix(users[i], "<@")
			users[i] = strings.TrimPrefix(users[i], "!")
			users[i] = strings.TrimSuffix(users[i], ">")

			_, err = s.GuildMember(config.GdHouseID, users[i])
			if err != nil {
				s.ChannelMessageSend(channelID, "В юзверях затесался шпион")
				return nil, resultErr
			}
		}

		usersSumDB.Us = append(usersSumDB.Us, database.Us{
			Users: users,
			Sum:   sum,
		})
	}

	usersSumJSON, err := json.Marshal(usersSumDB)
	if err != nil {
		fmt.Println("ERROR marshling to JSON:", err)
		return nil, resultErr
	}

	item.UsersSum = string(usersSumJSON)

	return item, nil
}

func resendRequest(s *discordgo.Session, channelID, content string, item *database.Record) error {
	resultErr := errors.New("resending failure")

	message, err := s.ChannelMessageSendEmbed(config.ChForRequestID, createEmbed(s, item, 225225))
	if err != nil {
		s.ChannelMessageSend(channelID, "Не удалось отправить запрос")
		fmt.Println("ERROR sending request", err)
		return resultErr
	}
	item.EmbedID = message.ID

	err = s.MessageReactionAdd(config.ChForRequestID, item.EmbedID, "✅")
	if err != nil {
		fmt.Println("ERROR adding reaction ✅", err)
		return resultErr
	}

	err = s.MessageReactionAdd(config.ChForRequestID, item.EmbedID, "🇽")
	if err != nil {
		fmt.Println("ERROR adding reaction 🇽", err)
		return resultErr
	}

	return nil
}

func emojiOnRequest(s *discordgo.Session, r *discordgo.MessageReactionAdd, item *database.Record) {
	switch r.Emoji.Name {
	case "✅":
		_, err = s.ChannelMessageEditEmbed(config.ChForRequestID, item.EmbedID, createEmbed(s, item, 255255))
		if err != nil {
			fmt.Println("ERROR edit embed on ✅", err)
			return
		}

	case "🇽":
		_, err = s.ChannelMessageEditEmbed(config.ChForRequestID, item.EmbedID, createEmbed(s, item, 15158332))
		if err != nil {
			fmt.Println("ERROR edit embed on 🇽", err)
			return
		}
	default:
		return
	}

	err = database.Records.Delete(item.EmbedID)
	if err != nil {
		fmt.Println("ERROR delete record", err)
		return
	}
}

func makeParsingBetter(str string) string {
	result := str
	for _, rep := range config.FairyReplacement {
		result = strings.Join(strings.Split(result, rep[0]), rep[1])
	}

	if result != str {
		return makeParsingBetter(result)
	}

	return result
}

func checkRequest(str string) (bool, string) {
	if strings.HasPrefix(str, "-запрос") {
		return true, strings.TrimPrefix(str, "-запрос")
	}
	if strings.HasPrefix(str, "- запрос") {
		return true, strings.TrimPrefix(str, "- запрос")
	}
	if strings.HasPrefix(str, "-  запрос") {
		return true, strings.TrimPrefix(str, "-  запрос")
	}

	return false, ""
}
