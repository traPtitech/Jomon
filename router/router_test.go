package router

import (
	"encoding/gob"
	"testing"
	"time"

	"github.com/traPtitech/Jomon/storage/mock_storage"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/model/mock_model"
	"github.com/traPtitech/Jomon/testutil/random"
	"go.uber.org/mock/gomock"
)

type MockRepository struct {
	*mock_model.MockAccountManagerRepository
	*mock_model.MockCommentRepository
	*mock_model.MockFileRepository
	*mock_model.MockRequestRepository
	*mock_model.MockRequestStatusRepository
	*mock_model.MockRequestFileRepository
	*mock_model.MockRequestTagRepository
	*mock_model.MockRequestTargetRepository
	*mock_model.MockTagRepository
	*mock_model.MockUserRepository
}

type MockStorage struct {
	*mock_storage.MockStorage
}

func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	return &MockRepository{
		MockAccountManagerRepository: mock_model.NewMockAccountManagerRepository(ctrl),
		MockCommentRepository:        mock_model.NewMockCommentRepository(ctrl),
		MockFileRepository:           mock_model.NewMockFileRepository(ctrl),
		MockRequestRepository:        mock_model.NewMockRequestRepository(ctrl),
		MockRequestStatusRepository:  mock_model.NewMockRequestStatusRepository(ctrl),
		MockRequestFileRepository:    mock_model.NewMockRequestFileRepository(ctrl),
		MockRequestTagRepository:     mock_model.NewMockRequestTagRepository(ctrl),
		MockRequestTargetRepository:  mock_model.NewMockRequestTargetRepository(ctrl),
		MockTagRepository:            mock_model.NewMockTagRepository(ctrl),
		MockUserRepository:           mock_model.NewMockUserRepository(ctrl),
	}
}

func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	return &MockStorage{
		MockStorage: mock_storage.NewMockStorage(ctrl),
	}
}

type TestHandlers struct {
	Handlers   Handlers
	Repository *MockRepository
	Storage    *MockStorage
}

func NewTestHandlers(_ *testing.T, ctrl *gomock.Controller) (*TestHandlers, error) {
	gob.Register(User{})
	repository := NewMockRepository(ctrl)
	storage := NewMockStorage(ctrl)
	sessionName := "session"

	return &TestHandlers{
		Handlers{
			Repository:  repository,
			Storage:     storage,
			SessionName: sessionName,
		},
		repository,
		storage,
	}, nil
}

func makeUser(t *testing.T, accountManager bool) *model.User {
	t.Helper()
	date := time.Now()

	return &model.User{
		ID:             uuid.New(),
		Name:           random.AlphaNumeric(t, 20),
		DisplayName:    random.AlphaNumeric(t, 20),
		AccountManager: accountManager,
		CreatedAt:      date,
		UpdatedAt:      date,
	}
}
