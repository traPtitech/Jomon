package router

import (
	"testing"
	"time"

	"github.com/traPtitech/Jomon/storage/mock_storage"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/model/mock_model"
	"github.com/traPtitech/Jomon/testutil/random"
	"go.uber.org/zap"
)

type MockRepository struct {
	*mock_model.MockCommentRepository
	*mock_model.MockFileRepository
	*mock_model.MockGroupBudgetRepository
	*mock_model.MockGroupRepository
	*mock_model.MockRequestRepository
	*mock_model.MockRequestStatusRepository
	*mock_model.MockRequestFileRepository
	*mock_model.MockRequestTagRepository
	*mock_model.MockRequestTargetRepository
	*mock_model.MockTagRepository
	*mock_model.MockTransactionDetailRepository
	*mock_model.MockTransactionTagRepository
	*mock_model.MockTransactionRepository
	*mock_model.MockUserRepository
}

func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	return &MockRepository{
		MockCommentRepository:           mock_model.NewMockCommentRepository(ctrl),
		MockFileRepository:              mock_model.NewMockFileRepository(ctrl),
		MockGroupBudgetRepository:       mock_model.NewMockGroupBudgetRepository(ctrl),
		MockGroupRepository:             mock_model.NewMockGroupRepository(ctrl),
		MockRequestRepository:           mock_model.NewMockRequestRepository(ctrl),
		MockRequestStatusRepository:     mock_model.NewMockRequestStatusRepository(ctrl),
		MockRequestFileRepository:       mock_model.NewMockRequestFileRepository(ctrl),
		MockRequestTagRepository:        mock_model.NewMockRequestTagRepository(ctrl),
		MockRequestTargetRepository:     mock_model.NewMockRequestTargetRepository(ctrl),
		MockTagRepository:               mock_model.NewMockTagRepository(ctrl),
		MockTransactionRepository:       mock_model.NewMockTransactionRepository(ctrl),
		MockTransactionDetailRepository: mock_model.NewMockTransactionDetailRepository(ctrl),
		MockTransactionTagRepository:    mock_model.NewMockTransactionTagRepository(ctrl),
		MockUserRepository:              mock_model.NewMockUserRepository(ctrl),
	}
}

type TestHandlers struct {
	Handlers   *Handlers
	Repository *MockRepository
}

func NewTestHandlers(_ *testing.T, ctrl *gomock.Controller) (*TestHandlers, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	repository := NewMockRepository(ctrl)
	sessionStore := sessions.NewCookieStore([]byte("session"))
	sessionName := "session"

	return &TestHandlers{&Handlers{
		Repository:   repository,
		Storage:      mock_storage.NewMockStorage(ctrl),
		Logger:       logger,
		SessionName:  sessionName,
		SessionStore: sessionStore,
	}, repository}, nil
}

func makeUser(t *testing.T, admin bool) *model.User {
	t.Helper()
	date := time.Now()

	return &model.User{
		ID:          uuid.New(),
		Name:        random.AlphaNumeric(t, 20),
		DisplayName: random.AlphaNumeric(t, 20),
		Admin:       admin,
		CreatedAt:   date,
		UpdatedAt:   date,
	}
}
