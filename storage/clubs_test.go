package storage

import (
	"testing"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func randomClub() *Club {
	desc := gofakeit.Paragraph(1, 1, 10, "")
	return &Club{
		OwnerID:     gofakeit.Numerify("test##############"),
		RoleID:      gofakeit.Numerify("test##############"),
		Title:       gofakeit.Word(),
		Symbol:      gofakeit.Emoji(),
		Description: &desc,
	}
}

func TestClubsCreate(t *testing.T) {
	tx, err := Tx()
	assert.NoError(t, err)

	var club Club
	club.randomize()
	err = Clubs.Create(tx, &club)
	assert.NoError(t, err)
	assert.NotZero(t, club.ID)

	tx.Commit()
}
