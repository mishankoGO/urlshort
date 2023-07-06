package repository

import (
	"fmt"
	"github.com/boltdb/bolt"
)

type Repo struct {
	DB *bolt.DB
}

func NewBoltRepo() (*Repo, error) {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return &Repo{DB: db}, nil
}

func (r *Repo) Close() error {
	err := r.DB.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) Update(key string, value string) error {
	err := r.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) View(key string) string {
	var v []byte
	r.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		v = b.Get([]byte(key))
		return nil
	})
	return string(v)
}
