package model

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/storage"
	"github.com/traPtitech/Jomon/storage/mock_storage"
	"github.com/traPtitech/Jomon/testutil"
)

func setup(t *testing.T) (*ent.Client, *mock_storage.MockStorage, error) {
	client, err := SetupTestEntClient(t)
	if err != nil {
		return nil, nil, err
	}
	/*
		os.Mkdir(os.Getenv("UPLOAD_DIR"), 0777)
	*/
	ctrl := gomock.NewController(t)
	storage := mock_storage.NewMockStorage(ctrl)
	if err != nil {
		return nil, nil, err
	}
	return client, storage, nil
}
