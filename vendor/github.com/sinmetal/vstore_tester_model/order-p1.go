package vstore_tester_model

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"go.mercari.io/datastore"
	"go.mercari.io/datastore/boom"
)

type OrderP1 struct {
	ID            string          `boom:"id" datastore:"-"`
	Email         string          `json:"email"`
	Price         int             `json:"price"`
	DetailKeys    []datastore.Key `json:"detailKeys"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
	SchemaVersion int             `json:"-"`
}

type OrderP1Detail struct {
	ID            string        `boom:"id" datastore:"-"`
	ItemKey       datastore.Key `json:"itemKey"`
	Price         int           `json:"price"`
	Count         int           `json:"count"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
	SchemaVersion int           `json:"-"`
}

type OrderP1Store struct{}

var _ datastore.PropertyLoadSaver = &OrderP1{}
var _ datastore.PropertyLoadSaver = &OrderP1Detail{}

func (model *OrderP1) Load(ctx context.Context, ps []datastore.Property) error {
	err := datastore.LoadStruct(ctx, model, ps)
	if err != nil {
		return err
	}

	return nil
}

func (model *OrderP1) Save(ctx context.Context) ([]datastore.Property, error) {
	model.SchemaVersion = 1
	if model.CreatedAt.IsZero() {
		model.CreatedAt = time.Now()
	}
	model.UpdatedAt = time.Now()

	return datastore.SaveStruct(ctx, model)
}

func (model *OrderP1Detail) Load(ctx context.Context, ps []datastore.Property) error {
	err := datastore.LoadStruct(ctx, model, ps)
	if err != nil {
		return err
	}

	return nil
}

func (model *OrderP1Detail) Save(ctx context.Context) ([]datastore.Property, error) {
	model.SchemaVersion = 1
	if model.CreatedAt.IsZero() {
		model.CreatedAt = time.Now()
	}
	model.UpdatedAt = time.Now()

	return datastore.SaveStruct(ctx, model)
}

func (store *OrderP1Store) Put(bm *boom.Boom, order *OrderP1, details []*OrderP1Detail) error {
	if len(details) < 20 {
		return fmt.Errorf("detail must be 20 pieces or less. details.len = %d", len(details))
	}

	tx, err := bm.NewTransaction()
	if err != nil {
		return errors.Wrap(err, "datastore.NewTransaction")
	}
	_, err = tx.Put(order)
	if err != nil {
		return errors.Wrap(err, "order.Put")
	}
	_, err = tx.PutMulti(details)
	if err != nil {
		return errors.Wrap(err, "orderDetail.PutMulti")
	}
	_, err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "datastore.Commit")
	}

	return nil
}

func (store *OrderP1Store) GetWithDetail(client datastore.Client, bm *boom.Boom, id string) (*OrderP1, []*OrderP1Detail, error) {
	o := OrderP1{
		ID: id,
	}

	err := bm.Get(&o)
	if err != nil {
		return nil, nil, errors.Wrap(err, "order.Get")
	}

	ds := make([]*OrderP1Detail, len(o.DetailKeys))
	for i, v := range o.DetailKeys {
		ds[i] = &OrderP1Detail{
			ID: v.Name(),
		}
	}
	err = bm.GetMulti(ds)
	if err != nil {
		return nil, nil, errors.Wrap(err, "orderDetail.GetMulti")
	}

	return &o, ds, nil
}
