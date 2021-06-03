package cache

import (
	"bytes"
	"encoding/gob"

	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/badger"
)

type Voice struct {
	ID          string
	Name        string
	Permissions []*discordgo.PermissionOverwrite
}

type voices struct{}

func (voices) key(id string) []byte {
	return []byte("voice." + id)
}
func (table *voices) Set(voice *Voice) error {
	return cache.Update(func(tx *badger.Txn) error {
		values := bytes.NewBufferString("")
		if err := gob.NewEncoder(values).Encode(voice); err != nil {
			return err
		}

		return tx.Set(table.key(voice.ID), values.Bytes())
	})
}

func (table *voices) Get(id string) (*Voice, error) {
	result := new(Voice)

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

			if err != nil {
				return err
			}

			return err
		})
		return err
	})

	if err != nil {
		return nil, err
	}

	return result, err
}

func (table *voices) Delete(id string) error {
	err := cache.Update(func(tx *badger.Txn) error {
		return tx.Delete(table.key(id))
	})

	if err != nil {
		return err
	}

	return nil
}
