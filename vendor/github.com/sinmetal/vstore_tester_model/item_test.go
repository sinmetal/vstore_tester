package vstore_tester_model

import (
	"context"
	"fmt"
	"testing"

	"go.mercari.io/datastore"
	"go.mercari.io/datastore/boom"
	"go.mercari.io/datastore/clouddatastore"
)

func TestItemStore_Update(t *testing.T) {
	ctx := context.Background()

	o := datastore.WithProjectID("souzoh-demo-gcp-001")
	client, err := clouddatastore.FromContext(ctx, o)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	bm := boom.FromClient(ctx, client)

	store := ItemStore{}

	item := Item{}
	if err := store.Put(bm, &item); err != nil {
		t.Fatal(err)
	}
	item.CryptKey = "hogeKey"

	if err := store.Update(bm, &item); err != nil {
		t.Fatal(err)
	}

	si := Item{}
	si.ID = item.ID
	if err := store.Get(bm, &si); err != nil {
		t.Fatal(err)
	}
	if e, g := item.CryptKey, si.CryptKey; e != g {
		t.Fatalf("expected item.CryptKey is %d; got %s", e, g)
	}
}

func TestKeyEncodeDecord(t *testing.T) {
	t.Skip()
	ctx := context.Background()

	o := datastore.WithProjectID("souzoh-demo-gcp-001")
	client, err := clouddatastore.FromContext(ctx, o)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	bm := boom.FromClient(ctx, client)

	item := Item{}
	item.Kind = "HogeKind"
	item.ID = 100
	beforeKey := bm.Key(&item)
	encodedKey := beforeKey.Encode()
	fmt.Println(encodedKey)
	afterKey, err := client.DecodeKey(encodedKey)
	if err != nil {
		t.Fatal(err)
	}
	if e, g := beforeKey, afterKey; e != g {
		t.Fatalf("expected Key %v; got %v", e, g)
	}
}
