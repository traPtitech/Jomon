package router

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/model/mock_model"
	"github.com/traPtitech/Jomon/service/mock_service"
	"github.com/traPtitech/Jomon/testutil/random"
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
	*mock_model.MockTagRepository
	*mock_model.MockTransactionDetailRepository
	*mock_model.MockTransactionTagRepository
	*mock_model.MockTransactionRepository
	*mock_model.MockUserRepository
}

type Service struct {
	*mock_service.MockService
}

type TestHandlers struct {
	Handler      *Handlers
	Repository   *Repository
	Logger       *zap.Logger
	Service      *Service
	SessionName  string
	SessionStore sessions.Store
	Echo         *echo.Echo
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
		MockTagRepository:               mock_model.NewMockTagRepository(ctrl),
		MockTransactionRepository:       mock_model.NewMockTransactionRepository(ctrl),
		MockTransactionDetailRepository: mock_model.NewMockTransactionDetailRepository(ctrl),
		MockTransactionTagRepository:    mock_model.NewMockTransactionTagRepository(ctrl),
		MockUserRepository:              mock_model.NewMockUserRepository(ctrl),
	}
}

func NewMockService(ctrl *gomock.Controller) *Service {
	return &Service{
		MockService: mock_service.NewMockService(ctrl),
	}
}

func SetupTestHandlers(t *testing.T, ctrl *gomock.Controller) (*TestHandlers, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	repository := NewMockEntRepository(ctrl)
	service := NewMockService(ctrl)
	sessionStore := sessions.NewCookieStore([]byte("session"))
	sessionName := "session"
	h := Handlers{
		Repository:   repository,
		Logger:       logger,
		Service:      service,
		SessionName:  sessionName,
		SessionStore: sessionStore,
	}

	e := echo.New()
	e.Use(session.Middleware(h.SessionStore))
	SetRouting(e, h)

	return &TestHandlers{
		Handler:      &h,
		Repository:   repository,
		Logger:       logger,
		Service:      service,
		SessionName:  sessionName,
		SessionStore: sessionStore,
		Echo:         e,
	}, nil
}

func (th *TestHandlers) doRequest(t *testing.T, method string, path string, reqBody interface{}, resBody interface{}) (statusCode int, rec *httptest.ResponseRecorder) {
	t.Helper()
	req := httptest.NewRequest(method, path, requestEncode(t, reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	th.Echo.ServeHTTP(rec, req)

	if resBody != nil {
		responseDecode(t, rec, resBody)
	}

	return rec.Code, rec
}

func (th *TestHandlers) doRequestWithLogin(t *testing.T, accessUser *model.User, method string, path string, reqBody interface{}, resBody interface{}) (statusCode int, rec *httptest.ResponseRecorder) {
	t.Helper()

	ctx := context.Background()
	th.Repository.MockUserRepository.
		EXPECT().
		GetUserByID(ctx, accessUser.ID).
		Return(accessUser, nil).
		AnyTimes()

	req := httptest.NewRequest(method, path, requestEncode(t, reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()

	sess, err := th.Handler.SessionStore.Get(req, th.Handler.SessionName)
	require.NoError(t, err)
	sess.Values["userID"] = accessUser.ID.String()
	err = sess.Save(req, rec)
	require.NoError(t, err)

	th.Echo.ServeHTTP(rec, req)

	if resBody != nil {
		responseDecode(t, rec, resBody)
	}

	return rec.Code, rec
}

func requestEncode(t *testing.T, body interface{}) *strings.Reader {
	t.Helper()

	b, err := json.Marshal(body)
	require.NoError(t, err)

	return strings.NewReader(string(b))
}

func responseDecode(t *testing.T, rec *httptest.ResponseRecorder, i interface{}) {
	t.Helper()

	err := json.Unmarshal(rec.Body.Bytes(), i)
	require.NoError(t, err)
}

func mustMakeUser(t *testing.T, admin bool) *model.User {
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
