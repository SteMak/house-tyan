package cache

import (
	"bytes"
	"encoding/gob"

	"github.com/dgraph-io/badger"
)

type Trigger struct {
	Name    string
	Answers []string
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

	if err != nil {
		return err
	}
	return nil
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

func (table *triggers) Delete(id string) error {
	err := cache.Update(func(tx *badger.Txn) error {
		return tx.Delete(table.key(id))
	})

	if err != nil {
		return err
	}
	return nil
}
