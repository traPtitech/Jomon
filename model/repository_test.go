package model

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/storage/mock_storage"
)

func setup(t *testing.T, ctx context.Context) (*ent.Client, *mock_storage.MockStorage, error) {
	t.Helper()
	client, err := SetupTestEntClient(t)
	if err != nil {
		return nil, nil, err
	}
	err = dropAll(t, ctx, client)
	if err != nil {
		return nil, nil, err
	}
	ctrl := gomock.NewController(t)
	storage := mock_storage.NewMockStorage(ctrl)
	if err != nil {
		return nil, nil, err
	}
	return client, storage, nil
}

func dropAll(t *testing.T, ctx context.Context, client *ent.Client) error {
	t.Helper()
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
