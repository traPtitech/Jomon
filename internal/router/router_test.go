package router

import (
	"encoding/gob"
	"testing"
	"time"

	"github.com/traPtitech/Jomon/internal/service"
	"github.com/traPtitech/Jomon/internal/service/mock_service"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/internal/testutil/random"
	"go.uber.org/mock/gomock"
)

type MockRepository struct {
	*mock_service.MockAccountManagerRepository
	*mock_service.MockCommentRepository
	*mock_service.MockFileRepository
	*mock_service.MockApplicationRepository
	*mock_service.MockApplicationStatusRepository
	*mock_service.MockApplicationFileRepository
	*mock_service.MockApplicationTagRepository
	*mock_service.MockApplicationTargetRepository
	*mock_service.MockTagRepository
	*mock_service.MockUserRepository
}

type MockStorage struct {
	*mock_service.MockStorage
}

func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	return &MockRepository{
		MockAccountManagerRepository:    mock_service.NewMockAccountManagerRepository(ctrl),
		MockCommentRepository:           mock_service.NewMockCommentRepository(ctrl),
		MockFileRepository:              mock_service.NewMockFileRepository(ctrl),
		MockApplicationRepository:       mock_service.NewMockApplicationRepository(ctrl),
		MockApplicationStatusRepository: mock_service.NewMockApplicationStatusRepository(ctrl),
		MockApplicationFileRepository:   mock_service.NewMockApplicationFileRepository(ctrl),
		MockApplicationTagRepository:    mock_service.NewMockApplicationTagRepository(ctrl),
		MockApplicationTargetRepository: mock_service.NewMockApplicationTargetRepository(ctrl),
		MockTagRepository:               mock_service.NewMockTagRepository(ctrl),
		MockUserRepository:              mock_service.NewMockUserRepository(ctrl),
	}
}

func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	return &MockStorage{
		MockStorage: mock_service.NewMockStorage(ctrl),
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
	service := service.New(repository, storage)
	sessionName := "session"

	return &TestHandlers{
		Handlers{
			Repository:  repository,
			Storage:     storage,
			Service:     service,
			SessionName: sessionName,
		},
		repository,
		storage,
	}, nil
}

func makeUser(t *testing.T, accountManager bool) *service.User {
	t.Helper()
	date := time.Now()

	return &service.User{
		ID:             uuid.New(),
		Name:           random.AlphaNumeric(t, 20),
		DisplayName:    random.AlphaNumeric(t, 20),
		AccountManager: accountManager,
		CreatedAt:      date,
		UpdatedAt:      date,
	}
}
