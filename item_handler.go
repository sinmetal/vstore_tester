package vstore_tester

import (
	"context"
	"time"

	"go.mercari.io/datastore"
	"go.mercari.io/datastore/aedatastore"
	"go.mercari.io/datastore/boom"

	aestore "google.golang.org/appengine/datastore"
)

func FromContext(ctx context.Context) (datastore.Client, error) {
	return aedatastore.FromContext(ctx)
}

type ItemAPI struct{}

type ItemAPIAllocatedIDResponse struct {
	ID int64
}

func (api *ItemAPI) AllocatedID(ctx context.Context) (*ItemAPIAllocatedIDResponse, error) {
	store := ItemStore{}

	client, err := FromContext(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	bm := boom.FromClient(ctx, client)
	item := &Item{
		Kind:     "ItemV1",
		Contents: []string{"From AllocatedID"},
	}
	k, err := store.AllocatedID(bm, item)
	if err != nil {
		return nil, err
	}
	item.ID = k.ID()
	err = store.Put(bm, item)
	if err != nil {
		return nil, err
	}

	return &ItemAPIAllocatedIDResponse{
		ID: k.ID(),
	}, nil
}

type ItemAPIAllocatedIDOrgResponse struct {
	Low  int64
	High int64
}

func (api *ItemAPI) AllocatedIDOrg(ctx context.Context) (*ItemAPIAllocatedIDOrgResponse, error) {
	low, high, err := aestore.AllocateIDs(ctx, "ItemV1", nil, 1)
	if err != nil {
		return nil, err
	}

	return &ItemAPIAllocatedIDOrgResponse{
		Low:  low,
		High: high,
	}, nil
}

type ItemAPIPostRequest struct {
	Contents []string
}

type ItemAPIPostResponse struct {
	Key       string
	Contents  []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (api *ItemAPI) Post(ctx context.Context, form *ItemAPIPostRequest) (*ItemAPIPostResponse, error) {
	store := ItemStore{}

	client, err := FromContext(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	bm := boom.FromClient(ctx, client)
	item := &Item{
		Kind:     "ItemV1",
		Contents: form.Contents,
	}
	err = store.Put(bm, item)

	return &ItemAPIPostResponse{
		Key:       bm.Key(item).Encode(),
		Contents:  item.Contents,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}, nil
}
