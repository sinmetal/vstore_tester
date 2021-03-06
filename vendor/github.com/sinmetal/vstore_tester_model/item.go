package vstore_tester_model

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"go.mercari.io/datastore"
	"go.mercari.io/datastore/boom"
)

type Contents []string

type Item struct {
	// FIXME panic !!! Parent datastore.Key `boom:"parent"`
	ID                int64     `boom:"id" datastore:"-"`
	Kind              string    `boom:"kind" datastore:"-"`
	Lot               string    `json:"lot"`
	Index             int       `json:"index"`
	Contents          Contents  `json:"contents"`
	ContentsOrg       []string  `json:"contentsOrg"`
	CryptKey          string    `json:"cryptKey"`          // SampleのためにJson出力する
	EncryptedContents string    `json:"encryptedContents"` // SampleのためにJson出力する
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	SchemaVersion     int       `json:"-"`
}

type ItemStore struct{}

func (item *Item) Load(ctx context.Context, ps []datastore.Property) error {
	err := datastore.LoadStruct(ctx, item, ps)
	if err != nil {
		return err
	}

	return nil
}

func (item *Item) Save(ctx context.Context) ([]datastore.Property, error) {
	item.SchemaVersion = 4
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now()
	}
	item.UpdatedAt = time.Now()

	return datastore.SaveStruct(ctx, item)
}

func (store *ItemStore) AllocatedID(bm *boom.Boom, item *Item) (datastore.Key, error) {
	k, err := bm.AllocateID(item)
	if err != nil {
		return nil, errors.Wrap(err, "datastore.AllocateID")
	}
	return k, nil
}

func (store *ItemStore) Put(bm *boom.Boom, item *Item) error {
	key, err := bm.Put(item)
	if err != nil {
		return errors.Wrap(err, "datastore.Put")
	}
	item.ID = key.ID()
	return nil
}

func (store *ItemStore) Get(bm *boom.Boom, item *Item) error {
	err := bm.Get(item)
	if err != nil {
		return errors.Wrap(err, "datastore.Get")
	}
	return nil
}

func (store *ItemStore) Update(bm *boom.Boom, item *Item) error {
	_, err := bm.RunInTransaction(func(tx *boom.Transaction) error {
		st := Item{
			Kind: item.Kind,
			ID:   item.ID,
		}
		if err := tx.Get(&st); err != nil {
			return errors.Wrap(err, "tx.Get")
		}
		st.Contents = item.Contents
		st.ContentsOrg = item.ContentsOrg
		st.CryptKey = item.CryptKey
		st.EncryptedContents = item.EncryptedContents
		_, err := tx.Put(&st)
		if err != nil {
			return errors.Wrap(err, "tx.Put")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "datastore.RunInTransaction")
	}

	return nil
}
