package storage

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func randomClub() *Club {
	desc := gofakeit.Paragraph(1, 1, 10, "")
	return &Club{
		OwnerID:     gofakeit.Numerify("test##############"),
		Title:       gofakeit.Word(),
		Symbol:      gofakeit.Emoji(),
		Description: &desc,
	}
}

func TestClubsCreate(t *testing.T) {
	tx, err := pgxconn.Begin()
	assert.NoError(t, err)

	var club Club
	club.randomize()
	err = Clubs.Create(tx, &club)
	assert.NoError(t, err)
	assert.NotZero(t, club.ID)

	tx.Commit()
}

func TestGetClubByUser(t *testing.T) {
	tx, err := pgxconn.Begin()
	assert.NoError(t, err)

	var club Club
	club.randomize()
	err = Clubs.Create(tx, &club)
	assert.NoError(t, err)
	assert.NotZero(t, club.ID)

	tx.Commit()

	c, err := Clubs.GetClubByUser(club.OwnerID)
	assert.NoError(t, err)
	assert.Equal(t, club.ID, c.ID)

	c, err = Clubs.GetClubByUser("not an id")
	assert.NoError(t, err)
	assert.Nil(t, c)
}

func TestGetExpired(t *testing.T) {
	clubs, err := Clubs.GetExpired()
	assert.NoError(t, err)
	fmt.Println(len(clubs))
}

func TestRemoveExpired(t *testing.T) {
	clubs, err := Clubs.RemoveExpired()
	assert.NoError(t, err)
	fmt.Println(len(clubs))
}
