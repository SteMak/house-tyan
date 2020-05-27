package mafia

import (
	"bytes"
	"math/rand"
	"strconv"
	"sync"

	"github.com/SteMak/house-tyan/libs/dgutils"
	"github.com/SteMak/house-tyan/out"
	"github.com/bwmarrin/discordgo"
)

var (
	err      error
	commands = map[string]interface{}{
		"create": &dgutils.Command{
			Function: _module.onMafiaCreate,
			Handlers: []func(*dgutils.MessageContext){
				_module.middlewareIsOwner,
			},
		},
		"start": &dgutils.Command{
			Function: _module.onMafiaStart,
			Handlers: []func(*dgutils.MessageContext){
				_module.middlewareIsOwner,
				_module.middlewareMafia,
				_module.middlewareMafiaIsNotStarted,
			},
		},
		"next": &dgutils.Command{
			Function: _module.onMafiaNext,
			Handlers: []func(*dgutils.MessageContext){
				_module.middlewareIsOwner,
				_module.middlewareMafia,
				_module.middlewareMafiaIsStarted,
				_module.middlewareMafiaVoteIsNotStated,
			},
		},
		"kill": &dgutils.Command{
			Function: _module.onMafiaKill,
			Handlers: []func(*dgutils.MessageContext){
				_module.middlewareIsOwner,
				_module.middlewareMafia,
				_module.middlewareMafiaIsStarted,
				_module.middlewareMafiaVoteIsNotStated,
			},
		},
		"immunity": &dgutils.Command{
			Function: _module.onMafiaImmunity,
			Handlers: []func(*dgutils.MessageContext){
				_module.middlewareIsOwner,
				_module.middlewareMafia,
				_module.middlewareMafiaIsStarted,
				_module.middlewareMafiaVoteIsNotStated,
			},
		},
		"jail": &dgutils.Command{
			Function: _module.onMafiaJail,
			Handlers: []func(*dgutils.MessageContext){
				_module.middlewareIsOwner,
				_module.middlewareMafia,
				_module.middlewareMafiaIsStarted,
				_module.middlewareMafiaVoteIsNotStated,
			},
		},
		"vote": &dgutils.Command{
			Function: _module.onMafiaVote,
			Handlers: []func(*dgutils.MessageContext){
				_module.middlewareIsOwner,
				_module.middlewareMafia,
				_module.middlewareMafiaIsStarted,
				_module.middlewareMafiaVoteIsNotStated,
			},
		},
		"vote.end": &dgutils.Command{
			Function: _module.onMafiaVoteEnd,
			Handlers: []func(*dgutils.MessageContext){
				_module.middlewareIsOwner,
				_module.middlewareMafia,
				_module.middlewareMafiaIsStarted,
				_module.middlewareMafiaVoteIsStated,
			},
		},
		"finish": &dgutils.Command{
			Function: _module.onMafiaFinish,
			Handlers: []func(*dgutils.MessageContext){
				_module.middlewareIsOwner,
				_module.middlewareMafia,
				_module.middlewareMafiaIsStarted,
				_module.middlewareMafiaVoteIsNotStated,
			},
		},
	}
)

func (bot *module) onMafiaCreate(ctx *dgutils.MessageContext) {
	bot.game = NewGame()

	joinMsg, err = bot.session.ChannelMessageSend(ctx.Message.ChannelID, "Что-бы присоединиться к мафии, нажмите на галочку.")
	if err != nil {
		out.Err(true, err)
		return
	}

	err = bot.session.MessageReactionAdd(joinMsg.ChannelID, joinMsg.ID, "✅")
	if err != nil {
		out.Err(true, err)
	}

	ch, err := bot.session.UserChannelCreate(ctx.Message.Author.ID)
	if err != nil {
		out.Err(true, err)
		return
	}
	infoMsg, err = bot.session.ChannelMessageSendEmbed(ch.ID, &discordgo.MessageEmbed{
		Title: "Информация о Мафии",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Участники",
				Value: "[пусто]",
			},
		},
	})
	if err != nil {
		out.Err(true, err)
	}
}

func (bot *module) onMafiaStart(ctx *dgutils.MessageContext) {
	for i, player := range bot.game.Players {
		player.Index = i + 1
	}

	err := bot.game.Random(ctx.Args)
	if err != nil {
		_, err = bot.session.ChannelMessageSend(ctx.Message.ChannelID, err.Error())
		if err != nil {
			out.Err(true, err)
		}
		return
	}

	updateInfoEmbed(bot)

	err = bot.session.ChannelMessageDelete(joinMsg.ChannelID, joinMsg.ID)
	if err != nil {
		out.Err(true, err)
	}
	joinMsg = nil

	_, err = bot.session.ChannelMessageSend(ctx.Message.ChannelID, "Мафия началась.")
	if err != nil {
		out.Err(true, err)
	}

	for _, player := range bot.game.Players {
		go func(player *Player) {
			ch, err := bot.session.UserChannelCreate(player.ID)
			if err != nil {
				out.Err(true, err)
			}

			if img, ok := bot.roles.Load(player.Role); ok {
				img := img.([]*bytes.Buffer)
				_, err = bot.session.ChannelFileSend(ch.ID, player.Role+".png", img[rand.Intn(len(img))])
				if err != nil {
					out.Err(true, err)
				}
				return
			}

			_, err = bot.session.ChannelMessageSend(ch.ID, "Ваша роль: "+player.Role)
			if err != nil {
				out.Err(true, err)
			}
		}(player)
	}
}

func (bot *module) onMafiaNext(ctx *dgutils.MessageContext) {
	_, err = bot.session.ChannelMessageSend(ctx.Message.ChannelID, "День #"+strconv.Itoa(bot.game.Day))
	if err != nil {
		out.Err(true, err)
	}

	killed := bot.game.Next()
	for _, player := range killed {
		go func(player *Player) {
			_, err = bot.session.ChannelMessageSend(ctx.Message.ChannelID, "<@"+player.ID+"> был убит.")
			if err != nil {
				out.Err(true, err)
			}
		}(player)
	}
}

func (bot *module) onMafiaKill(ctx *dgutils.MessageContext) {
	var wg sync.WaitGroup
	for _, user := range ctx.Message.Mentions {
		wg.Add(1)
		go func(user *discordgo.User) {
			defer wg.Done()
			err := bot.game.Kill(user.ID)
			if err != nil {
				_, err := bot.session.ChannelMessageSend(ctx.Message.ChannelID, err.Error()+": "+user.Mention())
				if err != nil {
					out.Err(true, err)
				}
				return
			}
		}(user)
	}

	updateInfoEmbed(bot)

	err := bot.session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, "👍")
	if err != nil {
		out.Err(true, err)
	}
}

func (bot *module) onMafiaImmunity(ctx *dgutils.MessageContext) {
	for _, user := range ctx.Message.Mentions {
		go func(user *discordgo.User) {
			err := bot.game.Immunity(user.ID)
			if err != nil {
				_, err := bot.session.ChannelMessageSend(ctx.Message.ChannelID, err.Error()+": "+user.Mention())
				if err != nil {
					out.Err(true, err)
				}
				return
			}
		}(user)
	}

	err := bot.session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, "👍")
	if err != nil {
		out.Err(true, err)
	}
}

func (bot *module) onMafiaJail(ctx *dgutils.MessageContext) {
	var wg sync.WaitGroup
	for _, user := range ctx.Message.Mentions {
		wg.Add(1)
		go func(user *discordgo.User) {
			defer wg.Done()
			err := bot.game.Jail(user.ID)
			if err != nil {
				_, err := bot.session.ChannelMessageSend(ctx.Message.ChannelID, err.Error()+": "+user.Mention())
				if err != nil {
					out.Err(true, err)
				}
				return
			}
		}(user)
	}

	updateInfoEmbed(bot)

	err := bot.session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, "👍")
	if err != nil {
		out.Err(true, err)
	}
}

func (bot *module) onMafiaVote(ctx *dgutils.MessageContext) {
	var wg sync.WaitGroup

	count := 0
	for _, user := range ctx.Message.Mentions {
		wg.Add(1)
		go func(user *discordgo.User) {
			defer wg.Done()

			if _, player := bot.game.GetPlayer(user.ID); player != nil {
				votePlayers.Store(player, 0)
				count++
			} else {
				_, err := bot.session.ChannelMessageSend(ctx.Message.ChannelID, user.Mention()+" не играет в мафию.")
				if err != nil {
					out.Err(true, err)
				}
			}
		}(user)
	}

	wg.Wait()

	if count < 2 {
		_, err := bot.session.ChannelMessageSend(ctx.Message.ChannelID, "Недостаточно кандидатов")
		if err != nil {
			out.Err(true, err)
		}
		return
	}

	voteMsg, err = bot.session.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "Голосование",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Напишите в чат номер игрока, против которого вы голосуете",
					Value: votesToString(&votePlayers),
				},
			},
		},
	})
	if err != nil {
		out.Err(true, err)
	}
}

func (bot *module) onMafiaVoteEnd(ctx *dgutils.MessageContext) {
	var prisoners []*Player
	var mostVotes int

	votePlayers.Range(func(_, votes interface{}) bool {
		if votes.(int) > mostVotes {
			mostVotes = votes.(int)
		}
		return true
	})

	votePlayers.Range(func(prisoner, votes interface{}) bool {
		if votes.(int) == mostVotes {
			prisoner := prisoner.(*Player)
			err := bot.game.Jail(prisoner.ID)
			if err != nil {
				_, err := bot.session.ChannelMessageSend(ctx.Message.ChannelID, err.Error()+": <@"+prisoner.ID+">")
				if err != nil {
					out.Err(true, err)
				}
				return true
			}
			prisoners = append(prisoners, prisoner)
		}
		return true
	})

	voteMsg = nil
	votePlayers = sync.Map{}
	votedPlayers = sync.Map{}

	_, err := bot.session.ChannelMessageSendEmbed(ctx.Message.ChannelID, &discordgo.MessageEmbed{
		Title: "Итоги голосования",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Отправились в тюрьму",
				Value: playersToString(prisoners, false),
			},
		},
	})
	if err != nil {
		out.Err(true, err)
	}
}

func (bot *module) onMafiaFinish(ctx *dgutils.MessageContext) {
	_, err = bot.session.ChannelMessageSendEmbed(ctx.Message.ChannelID, &discordgo.MessageEmbed{
		Title: "Мафия закончилась",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Выжившие",
				Value: playersToString(bot.game.Players, true),
			},
		},
	})
	if err != nil {
		out.Err(true, err)
	}
}

func playersToString(players []*Player, showRole bool) (pStr string) {
	if len(players) == 0 {
		pStr = "[пусто]"
	} else {
		for _, player := range players {
			if player.Index != 0 {
				pStr += "`" + strconv.Itoa(player.Index) + "`"
			}
			pStr += "<@" + player.ID + ">"
			if showRole && player.Role != "" {
				pStr += ": " + player.Role
			}
			pStr += "\n"
		}
	}

	return
}

func updateInfoEmbed(bot *module) {
	infoMsg.Embeds[0].Fields[0].Value = playersToString(bot.game.Players, true)

	_, err := bot.session.ChannelMessageEditEmbed(infoMsg.ChannelID, infoMsg.ID, infoMsg.Embeds[0])
	if err != nil {
		out.Err(true, err)
	}
}

func votesToString(players *sync.Map) (pStr string) {
	players.Range(func(p, v interface{}) bool {
		player := p.(*Player)
		votes := v.(int)

		pStr += "`" + strconv.Itoa(player.Index) + "` <@" + player.ID + ">"
		pStr += ": **" + strconv.Itoa(votes) + "**"
		pStr += "\n"

		return true
	})

	return
}

func updateVoteEmbed(bot *module) {
	voteMsg.Embeds[0].Fields[0].Value = votesToString(&votePlayers)

	_, err := bot.session.ChannelMessageEditEmbed(voteMsg.ChannelID, voteMsg.ID, voteMsg.Embeds[0])
	if err != nil {
		out.Err(true, err)
	}
}
