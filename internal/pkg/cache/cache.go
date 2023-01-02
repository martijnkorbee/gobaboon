package cache

import (
	"bytes"
	"encoding/gob"
)

type Cache interface {
	// GetConnection returns the underlying connection which can be casted to the respective type
	GetConnection() any

	// Close closes the cache pool
	Close() error

	Has(key string) (bool, error)
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, expiry ...int) error
	Forget(key string) error
	EmptyByMatch(key string) error
	Empty() error
}

type Entry map[string]interface{}

func encode(item Entry) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(item)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func decode(str string) (Entry, error) {
	item := Entry{}
	b := bytes.Buffer{}
	b.Write([]byte(str))
	d := gob.NewDecoder(&b)
	err := d.Decode(&item)
	if err != nil {
		return nil, err
	}
	return item, nil
}
