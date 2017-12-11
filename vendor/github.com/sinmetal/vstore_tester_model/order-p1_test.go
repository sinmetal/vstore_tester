package vstore_tester_model

import (
	"context"
	"fmt"
	"testing"

	"go.mercari.io/datastore"
	"go.mercari.io/datastore/boom"
	"go.mercari.io/datastore/clouddatastore"
)

func TestOrderP1Store_Put(t *testing.T) {
	ctx := context.Background()

	o := datastore.WithProjectID("souzoh-demo-gcp-001")
	client, err := clouddatastore.FromContext(ctx, o)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	bm := boom.FromClient(ctx, client)

	store := OrderP1Store{}
	err = store.Put(bm,
		&OrderP1{
			Email: "hoge@example.com",
		},
		[]*OrderP1Detail{
			&OrderP1Detail{
				ID: "d1",
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOrderP1Store_GetWithDetail(t *testing.T) {
	ctx := context.Background()

	o := datastore.WithProjectID("souzoh-demo-gcp-001")
	client, err := clouddatastore.FromContext(ctx, o)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	bm := boom.FromClient(ctx, client)

	orderID := "hoge"
	orderDetailID := fmt.Sprintf("%s-1", orderID)
	store := OrderP1Store{}
	err = store.Put(bm,
		&OrderP1{
			ID:         orderID,
			Email:      "hoge@example.com",
			DetailKeys: []datastore.Key{client.NameKey("OrderP1Detail", orderDetailID, nil)},
		},
		[]*OrderP1Detail{
			&OrderP1Detail{
				ID: orderDetailID,
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	{
		o, ds, err := store.GetWithDetail(client, bm, orderID)
		if err != nil {
			t.Fatal(err)
		}
		if e, g := orderID, o.ID; e != g {
			t.Fatalf("expected OrderID %s; got %s", e, g)
		}
		if e, g := 1, len(ds); e != g {
			t.Fatalf("expected Detail.Len %d; got %d", e, g)
		}
	}

}
