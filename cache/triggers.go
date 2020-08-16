package cache

import (
	"bytes"
	"encoding/gob"

	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/badger"
)

type Trigger struct {
	Name    string
	Answers []discordgo.MessageSend
}

type triggers struct{}

func (triggers) key(id string) []byte {
	return []byte("trigger." + id)
}

func (table *triggers) Set(trigger *Trigger) error {
	err := cache.Update(func(tx *badger.Txn) error {

		values := bytes.NewBufferString("")
		if err := gob.NewEncoder(values).Encode(trigger); err != nil {
			return err
		}

		return tx.Set(table.key(trigger.Name), values.Bytes())
	})

	return err
}

func (table *triggers) Get(id string) (*Trigger, error) {
	result := new(Trigger)

	err := cache.View(func(tx *badger.Txn) error {
		item, err := tx.Get(table.key(id))
		if err != nil {
			return err
		}

		err = item.Value(func(value []byte) error {
			reader := bytes.NewBuffer(value)
			return gob.NewDecoder(reader).Decode(result)
		})

		return err
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (table *triggers) GetList(page, perPage int) (*[]Trigger, error) {
	result := make([]Trigger, perPage)
	index := 0
	i := 0

	err := cache.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := table.key("")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			if index < page*perPage {
				index++
				continue
			}

			item := it.Item()
			err := item.Value(func(value []byte) error {
				reader := bytes.NewBuffer(value)
				err := gob.NewDecoder(reader).Decode(&result[index])
				if err != nil {
					return err
				}

				index++
				return nil
			})

			if err != nil {
				return err
			}

			i++
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (table *triggers) Delete(id string) error {
	err := cache.Update(func(tx *badger.Txn) error {
		return tx.Delete(table.key(id))
	})

	return err
}
