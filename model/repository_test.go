package model

import (
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
	os.Mkdir(testutil.GetEnvOrDefault("UPLOAD_DIR", "./uploads"), 0777)
	storage, err := storage.NewLocalStorage(testutil.GetEnvOrDefault("UPLOAD_DIR", "./uploads"))
	if err != nil {
		return nil, nil, err
	}
	return client, storage, nil
}
