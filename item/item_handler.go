package item

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"go.mercari.io/datastore"
	"go.mercari.io/datastore/boom"
	"go.mercari.io/datastore/clouddatastore"

	"github.com/pkg/errors"
	"github.com/sinmetal/slog"
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
	log := slog.Start(time.Now())
	defer log.Flush()

	log.Info("Item.Post")
	{
		body, err := json.Marshal(form)
		if err != nil {
			fmt.Println(err.Error())
		}
		log.Info(string(body))
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
		log.Errorf("failed FromContext. err = %s", err.Error())
		return nil, errors.Wrap(err, "FromContext")
	}
	defer client.Close()

	bm := boom.FromClient(ctx, client)
	err = Retry(func(attempt int) (retry bool, err error) {
		log.Info(formatDatastoreRetryCountLog(attempt))
		err = store.Put(bm, item)
		if err != nil {
			log.Infof("store.Put. err = %s", err.Error())
			return true, errors.Wrap(err, "store.Put")
		}
		return false, nil
	}, 8)
	if err != nil {
		log.Errorf("failed store.Put. err = %s", err.Error())
		return nil, err
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

func (api *ItemAPI) PostForCreateClientEveryTimeRetry(ctx context.Context, form *ItemAPIPostRequest) (*ItemAPIPostResponse, error) {
	log := slog.Start(time.Now())
	defer log.Flush()

	log.Info("Item.PostForCreateClientEveryTimeRetry")
	{
		body, err := json.Marshal(form)
		if err != nil {
			fmt.Println(err.Error())
		}
		log.Info(string(body))
	}

	item := &vtm.Item{
		Kind:        "ItemV1PostForCreateClientEveryTimeRetry",
		Lot:         form.Lot,
		Index:       form.Index,
		Contents:    form.Contents,
		ContentsOrg: form.Contents,
	}

	projectID, err := config.GetProjectID(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "config.GetProjectID")
	}

	store := vtm.ItemStore{}

	var bm *boom.Boom
	err = Retry(func(attempt int) (retry bool, err error) {
		log.Info(formatDatastoreRetryCountLog(attempt))
		client, err := FromContext(ctx, projectID)
		if err != nil {
			log.Infof("FromContext. err = %s", err.Error())
			return true, errors.Wrap(err, "FromContext")
		}
		defer client.Close()

		bm = boom.FromClient(ctx, client)
		err = store.Put(bm, item)
		if err != nil {
			log.Infof("store.Put. err = %s", err.Error())
			return true, errors.Wrap(err, "store.Put")
		}
		return false, nil
	}, 8)
	if err != nil {
		log.Errorf("failed store.Put. err = %s", err.Error())
		return nil, err
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
	log := slog.Start(time.Now())
	defer log.Flush()

	log.Info("Item.PostForOnlyOneClient")
	{
		body, err := json.Marshal(form)
		if err != nil {
			fmt.Println(err.Error())
		}
		log.Info(string(body))
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
		log.Errorf("failed client.GetDatastoreClient. err = %s", err.Error())
		return nil, errors.Wrap(err, "client.GetDatastoreClient")
	}

	bm := boom.FromClient(ctx, client)
	err = Retry(func(attempt int) (retry bool, err error) {
		log.Info(formatDatastoreRetryCountLog(attempt))
		err = store.Put(bm, item)
		if err != nil {
			log.Infof("store.Put. err = %s", err.Error())
			return true, errors.Wrap(err, "store.Put")
		}
		return false, nil
	}, 8)
	if err != nil {
		log.Errorf("failed store.Put. err = %s", err.Error())
		return nil, err
	}

	log.Infof("stored item id = %d", bm.Key(item).ID())

	return &ItemAPIPostResponse{
		Key:       bm.Key(item).Encode(),
		Lot:       item.Lot,
		Index:     item.Index,
		Contents:  item.Contents,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}, nil
}

func (api *ItemAPI) PostForOnlyOneClientOtherProject(ctx context.Context, form *ItemAPIPostRequest) (*ItemAPIPostResponse, error) {
	log := slog.Start(time.Now())
	defer log.Flush()

	log.Info("Item.PostForOnlyOneClientOtherProject")
	{
		body, err := json.Marshal(form)
		if err != nil {
			fmt.Println(err.Error())
		}
		log.Info(string(body))
	}

	item := &vtm.Item{
		Kind:        "ItemV1OnlyOneClient",
		Lot:         form.Lot,
		Index:       form.Index,
		Contents:    form.Contents,
		ContentsOrg: form.Contents,
	}

	store := vtm.ItemStore{}

	projectID := "identification-service-qa"
	client, err := client.GetDatastoreClientForOtherProject(ctx, projectID)
	if err != nil {
		log.Errorf("failed client.GetDatastoreClientForOtherProject. err = %s", err.Error())
		return nil, errors.Wrap(err, "client.GetDatastoreClient")
	}

	bm := boom.FromClient(ctx, client)
	err = Retry(func(attempt int) (retry bool, err error) {
		log.Info(formatDatastoreRetryCountLog(attempt))
		err = store.Put(bm, item)
		if err != nil {
			log.Infof("store.Put. err = %s", err.Error())
			return true, errors.Wrap(err, "store.Put")
		}
		return false, nil
	}, 8)
	if err != nil {
		log.Errorf("failed store.Put. err = %s", err.Error())
		return nil, err
	}

	log.Infof("stored item id = %d", bm.Key(item).ID())

	return &ItemAPIPostResponse{
		Key:       bm.Key(item).Encode(),
		Lot:       item.Lot,
		Index:     item.Index,
		Contents:  item.Contents,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}, nil
}

type ItemAPIPutRequest struct {
	Key string `json:"key"`
}

type ItemAPIPutResponse struct {
	Key       string    `json:"key"`
	Lot       string    `json:"lot"`
	Index     int       `json:"index"`
	Contents  []string  `json:"contents"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (api *ItemAPI) UpdateForOnlyOneClient(ctx context.Context, form *ItemAPIPutRequest) (*ItemAPIPutResponse, error) {
	log := slog.Start(time.Now())
	defer log.Flush()

	log.Info("Item.UpdateForOnlyOneClient")
	{
		body, err := json.Marshal(form)
		if err != nil {
			fmt.Println(err.Error())
		}
		log.Info(string(body))
	}

	projectID, err := config.GetProjectID(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "config.GetProjectID")
	}
	client, err := client.GetDatastoreClient(ctx, projectID)
	if err != nil {
		log.Errorf("failed client.GetDatastoreClient. err = %s", err.Error())
		return nil, errors.Wrap(err, "client.GetDatastoreClient")
	}
	store := vtm.ItemStore{}

	bm := boom.FromClient(ctx, client)
	key, err := client.DecodeKey(form.Key)
	if err != nil {
		log.Errorf("Invalid Key. err = %s", err.Error())
		return nil, errors.Wrap(err, "client.DecodeKey")
	}
	log.Infof("parameter key id = %d", key.ID())

	var si vtm.Item
	si.Kind = "ItemV1OnlyOneClient"
	si.ID = key.ID()
	if err := store.Get(bm, &si); err != nil {
		log.Errorf("failed datastore.Get. err = %s", err.Error())
		return nil, errors.Wrap(err, "datastore.Get")
	}

	// 適当に暗号化して更新する
	err = Retry(func(attempt int) (retry bool, err error) {
		log.Info(formatKMSRetryCountLog(attempt))
		cipterText, cryptKey, err := encrypt(ctx, si.Contents[0])
		if err != nil {
			log.Errorf("failed kms.encrypt. err = %s", err.Error())
			return true, errors.Wrap(err, "kms.encrypt")
		}
		si.EncryptedContents = cipterText
		si.CryptKey = cryptKey
		return false, nil
	}, 4)
	if err != nil {
		log.Errorf("failed KMS.encrypt. err = %s", err.Error())
		return nil, err
	}

	err = Retry(func(attempt int) (retry bool, err error) {
		log.Info(formatDatastoreRetryCountLog(attempt))
		err = store.Update(bm, &si)
		if err == datastore.ErrNoSuchEntity {
			return false, err
		}
		if err != nil {
			log.Infof("store.Update. err = %s", err.Error())
			return true, errors.Wrap(err, "store.Update")
		}
		return false, nil
	}, 8)
	if err != nil {
		log.Errorf("failed store.Update. err = %s", err.Error())
		return nil, err
	}

	return &ItemAPIPutResponse{
		Key:       bm.Key(si).Encode(),
		Lot:       si.Lot,
		Index:     si.Index,
		Contents:  si.Contents,
		CreatedAt: si.CreatedAt,
		UpdatedAt: si.UpdatedAt,
	}, nil
}

type ItemAPIGetRequest struct {
	Key string `json:"key" swagger:",in=query"`
}

type ItemAPIGetResponse struct {
	Key       string    `json:"key"`
	Lot       string    `json:"lot"`
	Index     int       `json:"index"`
	Contents  []string  `json:"contents"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (api *ItemAPI) GetForOnlyOneClient(ctx context.Context, form *ItemAPIGetRequest) (*ItemAPIGetResponse, error) {
	log := slog.Start(time.Now())
	defer log.Flush()

	log.Info("Item.GetForOnlyOneClient")
	{
		body, err := json.Marshal(form)
		if err != nil {
			fmt.Println(err.Error())
		}
		log.Info(string(body))
	}

	projectID, err := config.GetProjectID(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "config.GetProjectID")
	}
	client, err := client.GetDatastoreClient(ctx, projectID)
	if err != nil {
		log.Errorf("failed client.GetDatastoreClient. err = %s", err.Error())
		return nil, errors.Wrap(err, "client.GetDatastoreClient")
	}
	store := vtm.ItemStore{}

	bm := boom.FromClient(ctx, client)
	key, err := client.DecodeKey(form.Key)
	if err != nil {
		log.Errorf("Invalid Key. err = %s", err.Error())
		return nil, errors.Wrap(err, "client.DecodeKey")
	}

	var si vtm.Item
	si.Kind = "ItemV1OnlyOneClient"
	si.ID = key.ID()
	err = Retry(func(attempt int) (retry bool, err error) {
		if err := store.Get(bm, &si); err != nil {
			log.Errorf("failed datastore.Get. err = %s", err.Error())
			return true, errors.Wrap(err, "datastore.Get")
		}

		return false, nil
	}, 8)
	if err != nil {
		log.Errorf("failed datastore.Get. err = %s", err.Error())
		return nil, errors.Wrap(err, "datastore.Get")
	}

	// 復号化する
	err = Retry(func(attempt int) (retry bool, err error) {
		log.Info(formatKMSRetryCountLog(attempt))
		plainText, err := decrypt(ctx, si.EncryptedContents)
		if err != nil {
			log.Errorf("failed kms.encrypt. err = %s", err.Error())
			return true, errors.Wrap(err, "kms.encrypt")
		}
		log.Infof("plainText = %s", plainText)
		return false, nil
	}, 4)
	if err != nil {
		log.Errorf("failed KMS.encrypt. err = %s", err.Error())
		return nil, err
	}

	return &ItemAPIGetResponse{
		Key:       bm.Key(si).Encode(),
		Lot:       si.Lot,
		Index:     si.Index,
		Contents:  si.Contents,
		CreatedAt: si.CreatedAt,
		UpdatedAt: si.UpdatedAt,
	}, nil
}

type Func func(attempt int) (retry bool, err error)

func Retry(fn Func, maxRetries int) error {
	var err error

	attempt := 1
	for {
		b, err := fn(attempt)
		if b == false || err == nil {
			break
		}
		attempt++
		if attempt > maxRetries {
			return err
		}
		time.Sleep(time.Second*time.Duration(attempt*attempt) + time.Millisecond*time.Duration(rand.Intn(100*attempt)))
	}

	return err
}

func formatDatastoreRetryCountLog(attempt int) string {
	return fmt.Sprintf("Datastore Retry Count = %d", attempt)
}

func formatKMSRetryCountLog(attempt int) string {
	return fmt.Sprintf("KMS Retry Count = %d", attempt)
}
