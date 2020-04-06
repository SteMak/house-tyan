package cache

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/SteMak/house-tyan/config"
	"github.com/SteMak/house-tyan/out"
	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/badger"
)

type Blank struct {
	ID        string
	Author    discordgo.User
	Reason    string
	Rewards   []Reward
	ExpiresAt time.Time
}

type blanks struct{}

func (blanks) key(id string) []byte {
	return []byte("blank." + id)
}

func (table *blanks) Create(id, reason string, author *discordgo.User) (*Blank, error) {
	blank := &Blank{
		ID:     id,
		Reason: reason,
		Author: *author,
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
	out.Debugln(blank.ExpiresAt)
	return blank, nil
}

func (table *blanks) Get(id string) (*Blank, error) {
	result := new(Blank)

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
		result.ExpiresAt = time.Unix(int64(item.ExpiresAt()), 0).UTC()
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (table *blanks) Delete(id string) error {
	err := cache.Update(func(tx *badger.Txn) error {
		return tx.Delete(table.key(id))
	})

	if err != nil {
		return err
	}
	return nil
}
