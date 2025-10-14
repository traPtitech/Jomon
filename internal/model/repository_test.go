package model

import (
	"context"
	"testing"

	"github.com/traPtitech/Jomon/ent"
)

const (
	dbPrefix = "jomon_test_repo_"
)

//nolint:ireturn
func setup(t *testing.T, ctx context.Context, dbName string) (*ent.Client, error) {
	t.Helper()
	client, err := SetupTestEntClient(t, ctx, dbPrefix+dbName)
	if err != nil {
		return nil, err
	}
	err = dropAll(t, ctx, client)
	if err != nil {
		return nil, err
	}
	return client, nil
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
	_, err = client.Application.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = client.ApplicationStatus.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = client.ApplicationTarget.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = client.Tag.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = client.User.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
