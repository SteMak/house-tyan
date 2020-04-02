package cache

import (
	"bytes"
	"encoding/gob"

	"github.com/bwmarrin/discordgo"

	"github.com/SteMak/house-tyan/config"
	"github.com/dgraph-io/badger"
)

type usernames struct{}

func (usernames) key(id string) []byte {
	return []byte("username." + id)
}

func (table *usernames) Set(user *discordgo.User) error {
	err := cache.Update(func(tx *badger.Txn) error {

		values := bytes.NewBufferString("")
		if err := gob.NewEncoder(values).Encode(user); err != nil {
			return err
		}

		entry := badger.NewEntry(table.key(user.String()), values.Bytes()).
			WithTTL(config.Cache.TTL.Username.Duration)

		return tx.SetEntry(entry)
	})

	if err != nil {
		return err
	}
	return nil
}

func (table *usernames) Get(username string) (*discordgo.User, error) {
	result := new(discordgo.User)

	err := cache.View(func(tx *badger.Txn) error {
		item, err := tx.Get(table.key(username))
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

func (table *usernames) Delete(username string) error {
	err := cache.Update(func(tx *badger.Txn) error {
		return tx.Delete(table.key(username))
	})

	if err != nil {
		return err
	}
	return nil
}
