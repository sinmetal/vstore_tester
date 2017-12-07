package order

import (
	"context"
	"fmt"

	"go.mercari.io/datastore"
	"go.mercari.io/datastore/boom"
	"go.mercari.io/datastore/clouddatastore"

	"github.com/google/uuid"

	"github.com/sinmetal/vstore_tester/config"
	vtm "github.com/sinmetal/vstore_tester_model"
)

func FromContext(ctx context.Context, projectID string) (datastore.Client, error) {
	o := datastore.WithProjectID(projectID)
	return clouddatastore.FromContext(ctx, o)
}

type OrderP1API struct{}

type OrderPostRequest struct {
	Email   string
	Details []*DetailForm
}

type DetailForm struct {
	ItemID string `json:"itemID"`
	Count  int    `json:"count"`
}

type OrderPostResponse struct {
	Order   *vtm.OrderP1
	Details []*vtm.OrderP1Detail
}

func (api *OrderP1API) Post(ctx context.Context, form *OrderPostRequest) (*OrderPostResponse, error) {
	store := vtm.OrderP1Store{}

	projectID, err := config.GetProjectID(ctx)
	if err != nil {
		return nil, err
	}
	client, err := FromContext(ctx, projectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	id := uuid.New().String()
	bm := boom.FromClient(ctx, client)
	sp := 0
	ds := make([]*vtm.OrderP1Detail, len(form.Details))
	for i, v := range form.Details {
		ds[i] = &vtm.OrderP1Detail{
			ID:      fmt.Sprintf("%s-_-%d", id, i),
			ItemKey: client.NameKey("Item", v.ItemID, nil),
			Price:   1000,
			Count:   v.Count,
		}
		sp += ds[i].Price * ds[i].Count
	}
	order := &vtm.OrderP1{
		ID:    uuid.New().String(),
		Price: sp,
	}
	err = store.Put(bm, order, ds)
	if err != nil {
		return nil, err
	}
	return &OrderPostResponse{
		Order:   order,
		Details: ds,
	}, nil
}
