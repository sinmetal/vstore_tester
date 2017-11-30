package main

import (
	"context"
	"net/http"

	"github.com/sinmetal/vstore_tester/config"
	"github.com/sinmetal/vstore_tester/item"

	"github.com/favclip/ucon"
)

func main() {
	ucon.Middleware(UseContext)
	ucon.Orthodox()

	itemAPI := item.ItemAPI{}
	ucon.HandleFunc(http.MethodPost, "/item", itemAPI.Post)
	ucon.HandleFunc(http.MethodPost, "/item/allocatedid", itemAPI.AllocatedID)

	ucon.DefaultMux.Prepare()
	http.Handle("/", ucon.DefaultMux)

	ucon.ListenAndServe(":8080")
}

func UseContext(b *ucon.Bubble) error {
	if b.Context == nil {
		ctx := context.Background()
		b.Context = config.SetProjectID(ctx, config.ProjectID)
	} else {
		b.Context = config.SetProjectID(b.Context, config.ProjectID)
	}

	return b.Next()
}
