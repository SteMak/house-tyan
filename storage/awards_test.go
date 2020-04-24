package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/bwmarrin/discordgo"

	"github.com/SteMak/house-tyan/cache"
)

func randomBlank(rewardsCount, usersCount int) *cache.Blank {
	blank := &cache.Blank{
		Author: discordgo.User{
			ID: gofakeit.Numerify("test##############"),
		},
		Reason: gofakeit.Paragraph(1, 1, 5, ""),
	}

	for i := 0; i < rewardsCount; i++ {
		reward := cache.Reward{
			Amount: uint64(gofakeit.Number(20, 35000)),
			Users:  make(map[string]discordgo.User),
		}

		for j := 0; j < usersCount; j++ {
			id := gofakeit.Numerify("test##############")
			reward.Users[id] = discordgo.User{
				ID: id,
			}
		}

		blank.Rewards = append(blank.Rewards, reward)
	}
	return blank
}

func Test_Awards_Create(t *testing.T) {
	tx := Tx()
	id, err := Awards.Create(tx, randomBlank(3, 2))
	assert.NoError(t, err)
	assert.NotZero(t, id)

	tx.Commit()
}

func Test_Awards_SetStatus(t *testing.T) {
	tx := Tx()

	id, err := Awards.Create(tx, randomBlank(3, 2))
	assert.NoError(t, err)
	assert.NotZero(t, id)
	assert.NoError(t, Awards.SetStatus(tx, id, AwardStatusAccept))

	id, err = Awards.Create(tx, randomBlank(3, 2))
	assert.NoError(t, err)
	assert.NotZero(t, id)
	assert.NoError(t, Awards.SetStatus(tx, id, AwardStatusDiscard))

	tx.Commit()
}

func Test_Awards_SetBlankID(t *testing.T) {
	tx := Tx()

	id, err := Awards.Create(tx, randomBlank(3, 2))
	assert.NoError(t, err)
	assert.NotZero(t, id)
	assert.NoError(t, Awards.SetBlankID(tx, id, gofakeit.Numerify("test##############")))

	tx.Commit()
}

func Test_Awards_SetPaid(t *testing.T) {
	tx := Tx()

	blank := randomBlank(3, 2)
	id, err := Awards.Create(tx, blank)
	assert.NoError(t, err)
	assert.NotZero(t, id)

	blankID := gofakeit.Numerify("test##############")
	assert.NoError(t, Awards.SetBlankID(tx, id, blankID))

	for _, r := range blank.Rewards {
		for _, u := range r.Users {
			assert.NoError(t, Awards.SetPaid(tx, id, u.ID))
		}
	}
	tx.Commit()
}

func Test_Awards_Get(t *testing.T) {
	tx := Tx()

	blank := randomBlank(3, 2)
	id, err := Awards.Create(tx, blank)
	assert.NoError(t, err)
	assert.NotZero(t, id)
	assert.NoError(t, Awards.SetStatus(tx, id, AwardStatusAccept))

	blankID := gofakeit.Numerify("test##############")
	assert.NoError(t, Awards.SetBlankID(tx, id, blankID))

	tx.Commit()

	award, err := Awards.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, award.AuthorID, blank.Author.ID)
	assert.Equal(t, award.Reason, blank.Reason)
	assert.Equal(t, award.Status, AwardStatusAccept)
	if assert.NotNil(t, award.BlankID) {
		assert.Equal(t, *award.BlankID, blankID)
	}
}

func Test_Awards_BlankRewadrs(t *testing.T) {
	tx := Tx()

	blank := randomBlank(3, 2)
	id, err := Awards.Create(tx, blank)
	assert.NoError(t, err)
	assert.NotZero(t, id)
	assert.NoError(t, Awards.SetStatus(tx, id, AwardStatusAccept))

	blankID := gofakeit.Numerify("test##############")
	assert.NoError(t, Awards.SetBlankID(tx, id, blankID))

	tx.Commit()

	award, err := Awards.GetBlankRewards(blankID)
	assert.NoError(t, err)
	assert.Equal(t, award.AuthorID, blank.Author.ID)
	assert.Equal(t, award.Reason, blank.Reason)
	assert.Equal(t, award.Status, AwardStatusAccept)
	if assert.NotNil(t, award.BlankID) {
		assert.Equal(t, *award.BlankID, blankID)
	}
}
