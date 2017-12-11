package client

import (
	"context"
	"sync"

	"go.mercari.io/datastore"
	"go.mercari.io/datastore/clouddatastore"
)

var mu sync.RWMutex
var ds datastore.Client

func GetDatastoreClient(ctx context.Context, projectID string) (datastore.Client, error) {
	mu.RLock()
	defer mu.RUnlock()
	if ds == nil {
		createWithSetDatastoreClient(ctx, projectID)
	}
	return ds, nil
}

func CloseDatastoreClient() error {
	if ds == nil {
		return nil
	}
	return ds.Close()
}

func createWithSetDatastoreClient(ctx context.Context, projectID string) error {
	mu.Lock()
	defer mu.Unlock()
	client, err := fromContext(ctx, projectID)
	if err != nil {
		return err
	}
	ds = client
	return nil
}

func fromContext(ctx context.Context, projectID string) (datastore.Client, error) {
	o := datastore.WithProjectID(projectID)
	return clouddatastore.FromContext(ctx, o)
}
