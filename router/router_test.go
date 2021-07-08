package router

import (
	"github.com/golang/mock/gomock"
	"github.com/gorilla/sessions"
	"github.com/traPtitech/Jomon/model/mock_model"
	"github.com/traPtitech/Jomon/service/mock_service"
	"go.uber.org/zap"
)

type Repository struct {
	*mock_model.MockCommentRepository
	*mock_model.MockFileRepository
	*mock_model.MockGroupBudgetRepository
	*mock_model.MockGroupRepository
	*mock_model.MockRequestRepository
	*mock_model.MockRequestFileRepository
	*mock_model.MockRequestTagRepository
	*mock_model.MockRequestTargetRepository
	*mock_model.MockTransactionDetailRepository
	*mock_model.MockTransactionTagRepository
	*mock_model.MockTransactionRepository
	*mock_model.MockUserRepository
}

type Service struct {
	*mock_service.MockService
}

type TestHandlers struct {
	Repository   *Repository
	Logger       *zap.Logger
	Service      *Service
	SessionName  string
	SessionStore sessions.Store
}

func NewMockEntRepository(ctrl *gomock.Controller) *Repository {
	return &Repository{
		MockCommentRepository:           mock_model.NewMockCommentRepository(ctrl),
		MockFileRepository:              mock_model.NewMockFileRepository(ctrl),
		MockGroupBudgetRepository:       mock_model.NewMockGroupBudgetRepository(ctrl),
		MockGroupRepository:             mock_model.NewMockGroupRepository(ctrl),
		MockRequestRepository:           mock_model.NewMockRequestRepository(ctrl),
		MockRequestFileRepository:       mock_model.NewMockRequestFileRepository(ctrl),
		MockRequestTagRepository:        mock_model.NewMockRequestTagRepository(ctrl),
		MockRequestTargetRepository:     mock_model.NewMockRequestTargetRepository(ctrl),
		MockTransactionRepository:       mock_model.NewMockTransactionRepository(ctrl),
		MockTransactionDetailRepository: mock_model.NewMockTransactionDetailRepository(ctrl),
		MockTransactionTagRepository:    mock_model.NewMockTransactionTagRepository(ctrl),
		MockUserRepository:              mock_model.NewMockUserRepository(ctrl),
	}
}
