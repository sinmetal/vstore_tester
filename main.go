package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/sinmetal/vstore_tester/client"
	"github.com/sinmetal/vstore_tester/config"
	"github.com/sinmetal/vstore_tester/item"
	"github.com/sinmetal/vstore_tester/order"

	"github.com/favclip/ucon"
)

func main() {
	defer client.CloseDatastoreClient()

	ctx := context.Background()
	if err := item.SetUpKMS(ctx); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	ucon.Middleware(UseContext)
	ucon.Orthodox()

	itemAPI := item.ItemAPI{}
	ucon.HandleFunc(http.MethodPost, "/item", itemAPI.Post)
	ucon.HandleFunc(http.MethodPost, "/item/onlyoneclient", itemAPI.PostForOnlyOneClient)
	ucon.HandleFunc(http.MethodPut, "/item/onlyoneclient", itemAPI.UpdateForOnlyOneClient)
	ucon.HandleFunc(http.MethodGet, "/item/onlyoneclient", itemAPI.GetForOnlyOneClient)
	ucon.HandleFunc(http.MethodPost, "/item/onlyoneclientotherProject", itemAPI.PostForOnlyOneClientOtherProject)

	ucon.HandleFunc(http.MethodPost, "/item/createclienteverytimeretry", itemAPI.PostForCreateClientEveryTimeRetry)
	ucon.HandleFunc(http.MethodPost, "/item/allocatedid", itemAPI.AllocatedID)

	orderP1API := order.OrderP1API{}
	ucon.HandleFunc(http.MethodPost, "/orderP1", orderP1API.Post)

	ucon.DefaultMux.Prepare()
	http.Handle("/", ucon.DefaultMux)

	fmt.Println("Start Listen port 8080")
	ucon.ListenAndServe(":8080")

	fmt.Println("Exit")
}

// UseContext is Set Context
func UseContext(b *ucon.Bubble) error {
	if b.Context == nil {
		ctx := context.Background()
		b.Context = config.SetProjectID(ctx, config.ProjectID)
	} else {
		b.Context = config.SetProjectID(b.Context, config.ProjectID)
	}

	return b.Next()
}
