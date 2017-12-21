package client

import (
	"context"
	"sync"

	"go.mercari.io/datastore"
	"go.mercari.io/datastore/clouddatastore"
)

var mu sync.RWMutex
var ds datastore.Client

var muOther sync.RWMutex
var dsOther datastore.Client

func GetDatastoreClient(ctx context.Context, projectID string) (datastore.Client, error) {
	mu.RLock()
	if ds == nil {
		mu.RUnlock()
		createWithSetDatastoreClient(ctx, projectID)
	} else {
		mu.RUnlock()
	}
	return ds, nil
}

func GetDatastoreClientForOtherProject(ctx context.Context, projectID string) (datastore.Client, error) {
	muOther.RLock()
	if dsOther == nil {
		muOther.RUnlock()
		createWithSetDatastoreClientForOtherProject(ctx, projectID)
	} else {
		muOther.RUnlock()
	}
	return dsOther, nil
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

func createWithSetDatastoreClientForOtherProject(ctx context.Context, projectID string) error {
	muOther.Lock()
	defer muOther.Unlock()
	client, err := fromContext(ctx, projectID)
	if err != nil {
		return err
	}
	dsOther = client
	return nil
}

func fromContext(ctx context.Context, projectID string) (datastore.Client, error) {
	o := datastore.WithProjectID(projectID)
	return clouddatastore.FromContext(ctx, o)
}
