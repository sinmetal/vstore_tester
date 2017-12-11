package item

import (
	"context"
	"fmt"
	"time"

	"go.mercari.io/datastore"
	"go.mercari.io/datastore/boom"
	"go.mercari.io/datastore/clouddatastore"

	"github.com/pkg/errors"
	"github.com/sinmetal/vstore_tester/client"
	"github.com/sinmetal/vstore_tester/config"
	vtm "github.com/sinmetal/vstore_tester_model"
)

func FromContext(ctx context.Context, projectID string) (datastore.Client, error) {
	o := datastore.WithProjectID(projectID)
	return clouddatastore.FromContext(ctx, o)
}

type ItemAPI struct{}

type ItemAPIAllocatedIDResponse struct {
	ID int64
}

func (api *ItemAPI) AllocatedID(ctx context.Context) (*ItemAPIAllocatedIDResponse, error) {
	store := vtm.ItemStore{}

	projectID, err := config.GetProjectID(ctx)
	if err != nil {
		return nil, err
	}
	client, err := FromContext(ctx, projectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	bm := boom.FromClient(ctx, client)
	item := &vtm.Item{
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

type ItemAPIPostRequest struct {
	Contents []string
}

type ItemAPIPostResponse struct {
	ID        int64
	Contents  []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (api *ItemAPI) Post(ctx context.Context, form *ItemAPIPostRequest) (*ItemAPIPostResponse, error) {
	fmt.Println("Item Post !!!")

	item := &vtm.Item{
		Kind:     "ItemV1",
		Contents: form.Contents,
	}

	if err := api.insertItem(ctx, item); err != nil {
		fmt.Println(err.Error())
	}
	if err := api.insertItemOnlyOneClient(ctx, &vtm.Item{
		Kind:     "ItemV1OnlyOneClient",
		Contents: form.Contents,
	}); err != nil {
		fmt.Println(err.Error())
	}

	return &ItemAPIPostResponse{
		ID:        item.ID,
		Contents:  item.Contents,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}, nil
}

func (api *ItemAPI) insertItem(ctx context.Context, item *vtm.Item) error {
	store := vtm.ItemStore{}

	projectID, err := config.GetProjectID(ctx)
	if err != nil {
		return errors.Wrap(err, "config.GetProjectID")
	}
	client, err := FromContext(ctx, projectID)
	if err != nil {
		return errors.Wrap(err, "FromContext")
	}
	defer client.Close()

	bm := boom.FromClient(ctx, client)
	err = store.Put(bm, item)
	if err != nil {
		errors.Wrap(err, "store.Put")
	}
	return nil
}

func (api *ItemAPI) insertItemOnlyOneClient(ctx context.Context, item *vtm.Item) error {
	store := vtm.ItemStore{}

	projectID, err := config.GetProjectID(ctx)
	if err != nil {
		return errors.Wrap(err, "config.GetProjectID")
	}
	client, err := client.GetDatastoreClient(ctx, projectID)
	if err != nil {
		return errors.Wrap(err, "client.GetDatastoreClient")
	}

	bm := boom.FromClient(ctx, client)
	err = store.Put(bm, item)
	if err != nil {
		errors.Wrap(err, "store.Put")
	}
	return nil
}
