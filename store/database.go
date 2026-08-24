package store

import (
	"errors"
	"os"
	"path/filepath"

	"go.etcd.io/bbolt"
)

var ErrNotFound = errors.New("record not found")

type Database struct {
	path string
	db   *bbolt.DB
}

func Open(path string) (*Database, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	result := &Database{path: path, db: db}
	if err := result.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return result, nil
}

func (d *Database) initialize() error {
	return d.db.Update(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists([]byte(name)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (d *Database) Close() error {
	if d == nil || d.db == nil {
		return nil
	}
	return d.db.Close()
}

func (d *Database) Path() string {
	if d == nil {
		return ""
	}
	return d.path
}

func (d *Database) View(fn func(*bbolt.Tx) error) error {
	if d == nil || d.db == nil {
		return errors.New("database is closed")
	}
	return d.db.View(fn)
}

func (d *Database) Update(fn func(*bbolt.Tx) error) error {
	if d == nil || d.db == nil {
		return errors.New("database is closed")
	}
	return d.db.Update(fn)
}
