package cache

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/pkg/errors"

	"github.com/SteMak/house-tyan/config"
	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/badger"
)

type BlankActions struct {
	SetReason bool
	SetUsers  bool
	SetAmount bool
	Send      bool
	Discard   bool
}

type Blank struct {
	ID        string
	Message   discordgo.Message
	Author    discordgo.User
	Reason    string
	Rewards   []Reward
	ExpiresAt time.Time
	Actions   BlankActions
}

type blanks struct{}

func (blanks) key(id string) []byte {
	return []byte("blank." + id)
}

func (table *blanks) Create(id, reason string, author *discordgo.User, message *discordgo.Message) (*Blank, error) {
	blank := &Blank{
		ID:      id,
		Message: *message,
		Reason:  reason,
		Author:  *author,
		Actions: BlankActions{
			SetReason: true,
			SetUsers:  true,
			Discard:   true,
		},
	}

	err := cache.Update(func(tx *badger.Txn) error {
		values := bytes.NewBufferString("")
		if err := gob.NewEncoder(values).Encode(blank); err != nil {
			return err
		}

		entry := badger.NewEntry(table.key(id), values.Bytes()).
			WithTTL(config.Cache.TTL.Blank)

		blank.ExpiresAt = time.Unix(int64(entry.ExpiresAt), 0).UTC()

		return tx.SetEntry(entry)
	})

	if err != nil {
		return nil, err
	}
	return blank, nil
}

func (table *blanks) Set(blank *Blank) error {
	return cache.Update(func(tx *badger.Txn) error {
		values := bytes.NewBufferString("")
		if err := gob.NewEncoder(values).Encode(blank); err != nil {
			return err
		}

		entry := badger.NewEntry(table.key(blank.ID), values.Bytes()).
			WithTTL(config.Cache.TTL.Blank)

		blank.ExpiresAt = time.Unix(int64(entry.ExpiresAt), 0).UTC()

		return tx.SetEntry(entry)
	})
}

func (table *blanks) Get(id string) (*Blank, bool, error) {
	result := new(Blank)
	exists := false
	err := cache.View(func(tx *badger.Txn) error {
		item, err := tx.Get(table.key(id))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return nil
			}
			return err
		}
		exists = true

		err = item.Value(func(value []byte) error {
			reader := bytes.NewBuffer(value)
			if err := gob.NewDecoder(reader).Decode(result); err != nil {
				return err
			}
			return nil
		})
		result.ExpiresAt = time.Unix(int64(item.ExpiresAt()), 0).UTC()
		return nil
	})
	return result, exists, err
}

func (table *blanks) Delete(id string) error {
	return cache.Update(func(tx *badger.Txn) error {
		return tx.Delete(table.key(id))
	})
}

func (table *blanks) Exists(id string) (bool, error) {
	result := false
	err := cache.View(func(tx *badger.Txn) error {
		_, err := tx.Get(table.key(id))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return nil
			}
			return err
		}
		result = true
		return nil
	})
	return result, err
}
