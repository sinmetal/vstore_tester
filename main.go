package vstore_tester

import (
	"net/http"

	"github.com/favclip/ucon"

	"google.golang.org/appengine"
)

func init() {
	ucon.Middleware(UseAppengineContext)
	ucon.Orthodox()

	itemAPI := ItemAPI{}
	ucon.HandleFunc(http.MethodPost, "/item", itemAPI.Post)
	ucon.HandleFunc(http.MethodPost, "/item/allocatedid", itemAPI.AllocatedID)
	ucon.HandleFunc(http.MethodPost, "/item/allocatedidorg", itemAPI.AllocatedIDOrg)

	ucon.DefaultMux.Prepare()
	http.Handle("/", ucon.DefaultMux)
}

func UseAppengineContext(b *ucon.Bubble) error {
	if b.Context == nil {
		b.Context = appengine.NewContext(b.R)
	} else {
		b.Context = appengine.WithContext(b.Context, b.R)
	}

	return b.Next()
}
