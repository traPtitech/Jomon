package model

import (
	"context"
	"os"
	"testing"

	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/storage"
	"github.com/traPtitech/Jomon/testutil"
)

func setup(t *testing.T) (*ent.Client, storage.Storage, error) {
	client, err := SetupTestEntClient(t)
	if err != nil {
		return nil, nil, err
	}
	ctx := context.Background()
	err = dropAll(ctx, client)
	if err != nil {
		return nil, nil, err
	}
	os.Mkdir(os.Getenv("UPLOAD_DIR"), 0777)
	storage, err := storage.NewLocalStorage(testutil.GetEnvOrDefault("UPLOAD_DIR", "./uploads"))
	if err != nil {
		return nil, nil, err
	}
	return client, storage, nil
}

func dropAll(ctx context.Context, client *ent.Client) error {
	_, err := client.Comment.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = client.File.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = client.Group.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = client.GroupBudget.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = client.Request.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = client.RequestStatus.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = client.RequestTarget.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = client.Tag.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = client.Transaction.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = client.TransactionDetail.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = client.User.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
