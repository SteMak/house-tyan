package cache

import (
	"bytes"
	"encoding/gob"
	"sort"

	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/badger"
)

type Reward struct {
	Amount uint64
	Users  map[string]discordgo.User
}

type Award struct {
	ID      string
	Author  discordgo.User
	Reason  string
	Rewards []Reward
}

type awards struct{}

func (awards) key(id string) []byte {
	return []byte("award." + id)
}

func (table *awards) CreateFromBlank(id string, blank *Blank) error {

	sort.SliceStable(blank.Rewards, func(i, j int) bool {
		return blank.Rewards[i].Amount > blank.Rewards[j].Amount
	})

	award := &Award{
		ID:      id,
		Reason:  blank.Reason,
		Rewards: blank.Rewards,
	}

	err := cache.Update(func(tx *badger.Txn) error {

		values := bytes.NewBufferString("")
		if err := gob.NewEncoder(values).Encode(award); err != nil {
			return err
		}

		return tx.Set(table.key(award.ID), values.Bytes())
	})

	if err != nil {
		return err
	}
	return nil
}

func (table *awards) Set(award *Award) error {
	err := cache.Update(func(tx *badger.Txn) error {

		values := bytes.NewBufferString("")
		if err := gob.NewEncoder(values).Encode(award); err != nil {
			return err
		}

		return tx.Set(table.key(award.ID), values.Bytes())
	})

	if err != nil {
		return err
	}
	return nil
}

func (table *awards) Get(id string) (*Award, error) {
	result := new(Award)

	err := cache.View(func(tx *badger.Txn) error {
		item, err := tx.Get(table.key(id))
		if err != nil {
			return err
		}

		err = item.Value(func(value []byte) error {
			reader := bytes.NewBuffer(value)
			if err := gob.NewDecoder(reader).Decode(result); err != nil {
				return err
			}
			return nil
		})
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (table *awards) Delete(id string) error {
	err := cache.Update(func(tx *badger.Txn) error {
		return tx.Delete(table.key(id))
	})

	if err != nil {
		return err
	}
	return nil
}
