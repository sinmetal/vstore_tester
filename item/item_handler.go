package item

import (
	"context"
	"encoding/json"
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
	Lot      string   `json:"lot"`
	Index    int      `json:"index"`
	Contents []string `json:"contents"`
}

type ItemAPIPostResponse struct {
	Key       string    `json:"key"`
	Lot       string    `json:"lot"`
	Index     int       `json:"index"`
	Contents  []string  `json:"contents"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (api *ItemAPI) Post(ctx context.Context, form *ItemAPIPostRequest) (*ItemAPIPostResponse, error) {
	fmt.Println("Item Post !!!")
	{
		body, err := json.Marshal(form)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(string(body))
	}

	item := &vtm.Item{
		Kind:        "ItemV1",
		Lot:         form.Lot,
		Index:       form.Index,
		Contents:    form.Contents,
		ContentsOrg: form.Contents,
	}

	store := vtm.ItemStore{}

	projectID, err := config.GetProjectID(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "config.GetProjectID")
	}
	client, err := FromContext(ctx, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "FromContext")
	}
	defer client.Close()

	bm := boom.FromClient(ctx, client)
	err = store.Put(bm, item)
	if err != nil {
		return nil, errors.Wrap(err, "store.Put")
	}

	return &ItemAPIPostResponse{
		Key:       bm.Key(item).Encode(),
		Lot:       item.Lot,
		Index:     item.Index,
		Contents:  item.Contents,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}, nil
}

func (api *ItemAPI) PostForOnlyOneClient(ctx context.Context, form *ItemAPIPostRequest) (*ItemAPIPostResponse, error) {
	fmt.Println("Item PostForOnlyOneClient !!!")
	{
		body, err := json.Marshal(form)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(string(body))
	}

	item := &vtm.Item{
		Kind:        "ItemV1OnlyOneClient",
		Lot:         form.Lot,
		Index:       form.Index,
		Contents:    form.Contents,
		ContentsOrg: form.Contents,
	}

	store := vtm.ItemStore{}

	projectID, err := config.GetProjectID(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "config.GetProjectID")
	}
	client, err := client.GetDatastoreClient(ctx, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "client.GetDatastoreClient")
	}

	bm := boom.FromClient(ctx, client)
	err = store.Put(bm, item)
	if err != nil {
		return nil, errors.Wrap(err, "store.Put")
	}

	return &ItemAPIPostResponse{
		Key:       bm.Key(item).Encode(),
		Lot:       item.Lot,
		Index:     item.Index,
		Contents:  item.Contents,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}, nil
}
