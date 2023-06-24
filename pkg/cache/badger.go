package cache

import (
	"strings"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/martijnkorbee/gobaboon/pkg/logger"
	"github.com/martijnkorbee/gobaboon/pkg/util"
)

// BadgerCache
//
// TODO: add own logger
type BadgerCache struct {
	Connection *badger.DB
	Prefix     string
}

type BadgerConfig struct {
	Prefix string
	Path   string // Path is the full path to badger
}

// BadgerLogger replaces the default badger logger with a zerolog logger
type BadgerLogger struct {
	*logger.Logger
}

func (bl *BadgerLogger) Errorf(msg string, v ...interface{}) {
	msg = strings.ReplaceAll(msg, "\n", "")
	bl.Error().Msgf(msg, v...)
}

func (bl *BadgerLogger) Warningf(msg string, v ...interface{}) {
	msg = strings.ReplaceAll(msg, "\n", "")
	bl.Warn().Msgf(msg, v...)
}

func (bl *BadgerLogger) Infof(msg string, v ...interface{}) {
	msg = strings.ReplaceAll(msg, "\n", "")
	bl.Info().Msgf(msg, v...)
}

func (bl *BadgerLogger) Debugf(msg string, v ...interface{}) {
	msg = strings.ReplaceAll(msg, "\n", "")
	bl.Debug().Msgf(msg, v...)
}

func CreateBadgerCache(c BadgerConfig, log ...*logger.Logger) (*BadgerCache, error) {
	// create bader directory if it doesn't exist
	err := util.CreateDirIfNotExists(c.Path)
	if err != nil {
		return nil, err
	}

	opts := badger.DefaultOptions(c.Path)

	// create badger logger
	if len(log) > 0 {
		opts.Logger = &BadgerLogger{log[0]}
	}

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &BadgerCache{
		Connection: db,
		Prefix:     c.Prefix,
	}, nil
}

func (b *BadgerCache) GetConnection() any {
	return b.Connection
}

func (b *BadgerCache) Close() error {
	return b.Connection.Close()
}

func (b *BadgerCache) Has(str string) (bool, error) {
	_, err := b.Get(str)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (b *BadgerCache) Get(key string) (interface{}, error) {
	var fromCache []byte

	err := b.Connection.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			fromCache = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	decoded, err := decode(string(fromCache))
	if err != nil {
		return nil, err
	}

	item := decoded[key]

	return item, nil
}

func (b *BadgerCache) Set(key string, value interface{}, expires ...int) error {
	entry := Entry{}

	entry[key] = value
	encoded, err := encode(entry)
	if err != nil {
		return err
	}

	if len(expires) > 0 {
		err = b.Connection.Update(func(txn *badger.Txn) error {
			e := badger.NewEntry([]byte(key), encoded).WithTTL(time.Second * time.Duration(expires[0]))
			err = txn.SetEntry(e)
			return err
		})
	} else {
		err = b.Connection.Update(func(txn *badger.Txn) error {
			e := badger.NewEntry([]byte(key), encoded)
			err = txn.SetEntry(e)
			return err
		})
	}

	return nil
}

func (b *BadgerCache) Forget(key string) error {
	err := b.Connection.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		return err
	})

	return err
}

func (b *BadgerCache) EmptyByMatch(key string) error {
	return b.emptyByMatch(key)
}

func (b *BadgerCache) Empty() error {
	return b.emptyByMatch("")
}

func (b *BadgerCache) emptyByMatch(key string) error {
	deleteKeys := func(keysForDelete [][]byte) error {
		if err := b.Connection.Update(func(txn *badger.Txn) error {
			for _, x := range keysForDelete {
				if err := txn.Delete(x); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	}

	collectSize := 100000

	err := b.Connection.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.AllVersions = false
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		keysForDelete := make([][]byte, 0, collectSize)
		keysCollected := 0

		for it.Seek([]byte(key)); it.ValidForPrefix([]byte(key)); it.Next() {
			x := it.Item().KeyCopy(nil)
			keysForDelete = append(keysForDelete, x)
			keysCollected++
			if keysCollected == collectSize {
				if err := deleteKeys(keysForDelete); err != nil {
					return err
				}
			}
		}

		if keysCollected > 0 {
			if err := deleteKeys(keysForDelete); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
