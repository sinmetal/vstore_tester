package order

import (
	"context"
	"testing"

	"github.com/sinmetal/vstore_tester/config"
	_ "go.mercari.io/datastore/clouddatastore"
)

func TestOrderP1API_Post(t *testing.T) {
	ctx := context.Background()
	ctx = config.SetProjectID(ctx, "souzoh-demo-gcp-001")

	api := OrderP1API{}
	_, err := api.Post(ctx, &OrderPostRequest{
		Email: "hoge@example.com",
		Details: []*DetailForm{
			&DetailForm{
				ItemID: "1",
				Count:  1,
			},
			&DetailForm{
				ItemID: "2",
				Count:  2,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}
