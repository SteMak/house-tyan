package messages

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func normalizeSpaces(in []byte) []byte {
	var result []byte

	for _, line := range bytes.Split(in, []byte("\n")) {
		line = append(bytes.TrimSpace(line), byte('\n'))
		if string(line) == "\n" {
			continue
		}
		result = append(result, line...)
	}

	return result
}

func buildMessage(data *shema) (*Message, error) {
	result := new(Message)

	if content := strings.TrimSpace(data.Content); content != "" {
		result.Content = content
	}

	if embed := data.Embed; embed != nil {
		result.Embed = new(discordgo.MessageEmbed)

		if title := embed.Title; title != nil {
			result.Embed.Title = *title
		}

		if color := embed.Color; color != nil {
			if len([]rune(*color)) != 7 {
				return nil, fmt.Errorf("Unexpected color format: '%s'", *color)
			}

			c, err := strconv.ParseInt((*color)[1:], 16, 32)
			if err != nil {
				return nil, err
			}
			result.Embed.Color = int(c)
		}

		if description := strings.TrimSpace(embed.Description); description != "" {
			result.Embed.Description = description
		}

		if footer := embed.Footer; footer != nil {
			result.Embed.Footer = new(discordgo.MessageEmbedFooter)
			result.Embed.Footer.Text = *footer
		}

		if fields := embed.Fields; fields != nil {
			result.Embed.Fields = make([]*discordgo.MessageEmbedField, 0)
			for _, field := range *fields {
				result.Embed.Fields = append(result.Embed.Fields, &discordgo.MessageEmbedField{
					Inline: field.Inline,
					Name:   strings.TrimSpace(field.Name),
					Value:  strings.TrimSpace(field.Value),
				})
			}
		}
	}

	if reactions := data.Reactions; reactions != nil {
		result.Reactions = make([]string, len(*reactions))
		copy(result.Reactions, *reactions)
	}

	return result, nil
}
