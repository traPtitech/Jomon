package service

import (
	"github.com/golang/mock/gomock"
	"github.com/traPtitech/Jomon/storage/mock_storage"
)

type Storage struct {
	*mock_storage.MockStorage
}

func NewMockStorage(ctrl *gomock.Controller) *Storage {
	return &Storage{
		MockStorage: mock_storage.NewMockStorage(ctrl),
	}
}
