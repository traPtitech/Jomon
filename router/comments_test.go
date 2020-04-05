package router

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/traPtitech/Jomon/model"
)

type commentRepositoryMock struct {
	mock.Mock
	asr *assert.Assertions
}

func NewCommentRepositoryMock(t *testing.T) *commentRepositoryMock {
	m := new(commentRepositoryMock)
	m.asr = assert.New(t)
	return m
}

func (m *commentRepositoryMock) GetComment(applicationId uuid.UUID, commentId int) (model.Comment, error) {
	ret := m.Called(applicationId, commentId)
	return ret.Get(0).(model.Comment), ret.Error(1)
}

func (m *commentRepositoryMock) CreateComment(applicationId uuid.UUID, commentText string, userId string) (model.Comment, error) {
	ret := m.Called(applicationId, commentText, userId)

	m.asr.NotEqual("", commentText)

	return ret.Get(0).(model.Comment), ret.Error(1)
}

func (m *commentRepositoryMock) PutComment(applicationId uuid.UUID, commentId int, commentText string) (model.Comment, error) {
	ret := m.Called(applicationId, commentId, commentText)

	m.asr.NotEqual("", commentText)

	return ret.Get(0).(model.Comment), ret.Error(1)
}

func (m *commentRepositoryMock) DeleteComment(applicationId uuid.UUID, commentId int) error {
	ret := m.Called(applicationId, commentId)
	return ret.Error(0)
}

func GenerateComment(
	appId uuid.UUID,
	commentId int,
	createUserTrapID string,
	commentText string,
) model.Comment {
	return model.Comment{
		ID:            commentId,
		ApplicationID: appId,
		UserTrapID: model.User{
			TrapId: createUserTrapID,
		},
		Comment: commentText,
	}
}

func TestPostComment(t *testing.T) {
	appRepMock := NewApplicationRepositoryMock(t)

	title := "夏コミの交通費をお願いします。"
	remarks := "〇〇駅から〇〇駅への移動"
	amount := 1000
	paidAt := time.Now().Round(time.Second)

	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	appRepMock.On("GetApplication", id, mock.Anything).Return(GenerateApplication(id, "User2", model.ApplicationType{Type: model.Contest}, title, remarks, amount, paidAt), nil)
	appRepMock.On("GetApplication", mock.Anything, mock.Anything).Return(model.Application{}, gorm.ErrRecordNotFound)

	adminRepMock := NewAdministratorRepositoryMock("AdminUserId")

	commentRepoMock := NewCommentRepositoryMock(t)

	commentId := int(randSrc.Int63())
	commentText := "This is comment."
	commentRepoMock.On("CreateComment", id, commentText, "UserId").Return(GenerateComment(id, commentId, "UserId", commentText), nil)

	userRepMock := NewUserRepositoryMock(t, "UserId", "AdminUserId")

	service := Service{
		Administrators: adminRepMock,
		Applications:   appRepMock,
		Comments:       commentRepoMock,
		Users:          userRepMock,
	}

	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		body := fmt.Sprintf(`
		{
			"comment":"%s"
		}
		`, commentText)

		req := httptest.NewRequest(http.MethodPost, "/api/applications/"+id.String()+"/comments", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, "application/json")
		req.Header.Set("Authorization", userRepMock.token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/applications/:applicationId/comments")
		c.SetParamNames("applicationId")
		c.SetParamValues(id.String())

		route, pathParam, err := router.FindRoute(req.Method, req.URL)
		if err != nil {
			panic(err)
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		c, err = service.SetMyUser(c)
		if err != nil {
			panic(err)
		}

		err = service.PostComments(c)
		asr.NoError(err)
		asr.Equal(http.StatusCreated, rec.Code)

		err = validateResponse(&ctx, requestValidationInput, rec)
		asr.NoError(err)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		anotherId, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}

		body := fmt.Sprintf(`
		{
			"comment":"%s"
		}
		`, commentText)

		req := httptest.NewRequest(http.MethodPost, "/api/applications/"+anotherId.String()+"/comments", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, "application/json")
		req.Header.Set("Authorization", userRepMock.token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/applications/:applicationId/comments")
		c.SetParamNames("applicationId")
		c.SetParamValues(anotherId.String())

		route, pathParam, err := router.FindRoute(req.Method, req.URL)
		if err != nil {
			panic(err)
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		c, err = service.SetMyUser(c)
		if err != nil {
			panic(err)
		}

		err = service.PostComments(c)
		asr.NoError(err)

		asr.Equal(http.StatusNotFound, rec.Code)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		body := fmt.Sprintf(`
		{
			"comment":""
		}
		`)

		req := httptest.NewRequest(http.MethodPost, "/api/applications/"+id.String()+"/comments", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, "application/json")
		req.Header.Set("Authorization", userRepMock.token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/applications/:applicationId/comments")
		c.SetParamNames("applicationId")
		c.SetParamValues(id.String())

		route, pathParam, err := router.FindRoute(req.Method, req.URL)
		if err != nil {
			panic(err)
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		c, err = service.SetMyUser(c)
		if err != nil {
			panic(err)
		}

		err = service.PostComments(c)
		asr.NoError(err)
		asr.Equal(http.StatusBadRequest, rec.Code)
	})
}
func TestPutComment(t *testing.T) {
	appRepMock := NewApplicationRepositoryMock(t)

	title := "夏コミの交通費をお願いします。"
	remarks := "〇〇駅から〇〇駅への移動"
	amount := 1000
	paidAt := time.Now().Round(time.Second)

	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	appRepMock.On("GetApplication", id, mock.Anything).Return(GenerateApplication(id, "User2", model.ApplicationType{Type: model.Contest}, title, remarks, amount, paidAt), nil)
	appRepMock.On("GetApplication", mock.Anything, mock.Anything).Return(model.Application{}, gorm.ErrRecordNotFound)

	adminRepMock := NewAdministratorRepositoryMock("AdminUserId")

	commentRepoMock := NewCommentRepositoryMock(t)

	commentId := int(randSrc.Int63())
	commentText := "This is comment."
	commentRepoMock.On("GetComment", id, commentId).Return(GenerateComment(id, commentId, "UserId", commentText), nil)
	commentRepoMock.On("GetComment", mock.Anything, mock.Anything).Return(model.Comment{}, gorm.ErrRecordNotFound)

	newCommentText := "This is new comment."
	commentRepoMock.On("PutComment", id, commentId, newCommentText).Return(GenerateComment(id, commentId, "UserId", newCommentText), nil)
	commentRepoMock.On("PutComment", mock.Anything, mock.Anything, mock.Anything).Return(model.Comment{}, gorm.ErrRecordNotFound)

	anotherToken := "AnotherToken"

	userRepMock := NewUserRepositoryMock(t, "UserId", "AdminUserId")
	userRepMock.On("GetMyUser", anotherToken).Return(model.User{TrapId: "AnotherId"}, nil)

	service := Service{
		Administrators: adminRepMock,
		Applications:   appRepMock,
		Comments:       commentRepoMock,
		Users:          userRepMock,
	}

	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		body := fmt.Sprintf(`
		{
			"comment":"%s"
		}
		`, newCommentText)

		req := httptest.NewRequest(http.MethodPut, "/api/applications/"+id.String()+"/comments/"+strconv.Itoa(commentId), strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, "application/json")
		req.Header.Set("Authorization", userRepMock.token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/applications/:applicationId/comments/:commentId")
		c.SetParamNames("applicationId", "commentId")
		c.SetParamValues(id.String(), strconv.Itoa(commentId))

		route, pathParam, err := router.FindRoute(req.Method, req.URL)
		if err != nil {
			panic(err)
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		c, err = service.SetMyUser(c)
		if err != nil {
			panic(err)
		}

		err = service.PutComments(c)
		asr.NoError(err)
		asr.Equal(http.StatusOK, rec.Code)

		err = validateResponse(&ctx, requestValidationInput, rec)
		asr.NoError(err)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		anotherCommentId := int(randSrc.Int63())

		body := fmt.Sprintf(`
		{
			"comment":"%s"
		}
		`, commentText)

		req := httptest.NewRequest(http.MethodPut, "/api/applications/"+id.String()+"/comments/"+strconv.Itoa(anotherCommentId), strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, "application/json")
		req.Header.Set("Authorization", userRepMock.token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/applications/:applicationId/comments/:commentId")
		c.SetParamNames("applicationId", "commentId")
		c.SetParamValues(id.String(), strconv.Itoa(anotherCommentId))

		route, pathParam, err := router.FindRoute(req.Method, req.URL)
		if err != nil {
			panic(err)
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		c, err = service.SetMyUser(c)
		if err != nil {
			panic(err)
		}

		err = service.PutComments(c)
		asr.NoError(err)
		asr.Equal(http.StatusNotFound, rec.Code)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		body := fmt.Sprintf(`
		{
			"comment":"%s"
		}
		`, commentText)

		req := httptest.NewRequest(http.MethodPut, "/api/applications/"+id.String()+"/comments/"+strconv.Itoa(commentId), strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, "application/json")
		req.Header.Set("Authorization", anotherToken)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/applications/:applicationId/comments/:commentId")
		c.SetParamNames("applicationId", "commentId")
		c.SetParamValues(id.String(), strconv.Itoa(commentId))

		route, pathParam, err := router.FindRoute(req.Method, req.URL)
		if err != nil {
			panic(err)
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		c, err = service.SetMyUser(c)
		if err != nil {
			panic(err)
		}

		err = service.PutComments(c)
		asr.NoError(err)
		asr.Equal(http.StatusForbidden, rec.Code)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		body := fmt.Sprintf(`
		{
			"comment":""
		}
		`)

		req := httptest.NewRequest(http.MethodPut, "/api/applications/"+id.String()+"/comments/"+strconv.Itoa(commentId), strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, "application/json")
		req.Header.Set("Authorization", userRepMock.token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/applications/:applicationId/comments/:commentId")
		c.SetParamNames("applicationId", "commentId")
		c.SetParamValues(id.String(), strconv.Itoa(commentId))

		route, pathParam, err := router.FindRoute(req.Method, req.URL)
		if err != nil {
			panic(err)
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		c, err = service.SetMyUser(c)
		if err != nil {
			panic(err)
		}

		err = service.PutComments(c)
		asr.NoError(err)
		asr.Equal(http.StatusBadRequest, rec.Code)
	})
}

func TestDeleteComment(t *testing.T) {
	appRepMock := NewApplicationRepositoryMock(t)

	title := "夏コミの交通費をお願いします。"
	remarks := "〇〇駅から〇〇駅への移動"
	amount := 1000
	paidAt := time.Now().Round(time.Second)

	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	appRepMock.On("GetApplication", id, mock.Anything).Return(GenerateApplication(id, "User2", model.ApplicationType{Type: model.Contest}, title, remarks, amount, paidAt), nil)
	appRepMock.On("GetApplication", mock.Anything, mock.Anything).Return(model.Application{}, gorm.ErrRecordNotFound)

	adminRepMock := NewAdministratorRepositoryMock("AdminUserId")

	commentRepoMock := NewCommentRepositoryMock(t)

	commentId := int(randSrc.Int63())
	commentText := "This is comment."
	commentRepoMock.On("GetComment", id, commentId).Return(GenerateComment(id, commentId, "UserId", commentText), nil)
	commentRepoMock.On("GetComment", mock.Anything, mock.Anything).Return(model.Comment{}, gorm.ErrRecordNotFound)

	commentRepoMock.On("DeleteComment", id, commentId).Return(nil)
	commentRepoMock.On("DeleteComment", mock.Anything, mock.Anything).Return(gorm.ErrRecordNotFound)

	anotherToken := "AnotherToken"

	userRepMock := NewUserRepositoryMock(t, "UserId", "AdminUserId")
	userRepMock.On("GetMyUser", anotherToken).Return(model.User{TrapId: "AnotherId"}, nil)

	service := Service{
		Administrators: adminRepMock,
		Applications:   appRepMock,
		Comments:       commentRepoMock,
		Users:          userRepMock,
	}

	t.Parallel()

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		req := httptest.NewRequest(http.MethodDelete, "/api/applications/"+id.String()+"/comments/"+strconv.Itoa(commentId), nil)
		req.Header.Set("Authorization", userRepMock.token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/applications/:applicationId/comments/:commentId")
		c.SetParamNames("applicationId", "commentId")
		c.SetParamValues(id.String(), strconv.Itoa(commentId))

		route, pathParam, err := router.FindRoute(req.Method, req.URL)
		if err != nil {
			panic(err)
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		c, err = service.SetMyUser(c)
		if err != nil {
			panic(err)
		}

		err = service.DeleteComments(c)
		asr.NoError(err)
		asr.Equal(http.StatusNoContent, rec.Code)

		err = validateResponse(&ctx, requestValidationInput, rec)
		asr.NoError(err)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		anotherCommentId := int(randSrc.Int63())

		req := httptest.NewRequest(http.MethodDelete, "/api/applications/"+id.String()+"/comments/"+strconv.Itoa(anotherCommentId), nil)
		req.Header.Set("Authorization", userRepMock.token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/applications/:applicationId/comments/:commentId")
		c.SetParamNames("applicationId", "commentId")
		c.SetParamValues(id.String(), strconv.Itoa(anotherCommentId))

		route, pathParam, err := router.FindRoute(req.Method, req.URL)
		if err != nil {
			panic(err)
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		c, err = service.SetMyUser(c)
		if err != nil {
			panic(err)
		}

		err = service.DeleteComments(c)
		asr.NoError(err)
		asr.Equal(http.StatusNotFound, rec.Code)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)
		e := echo.New()
		ctx := context.TODO()

		req := httptest.NewRequest(http.MethodDelete, "/api/applications/"+id.String()+"/comments/"+strconv.Itoa(commentId), nil)
		req.Header.Set("Authorization", anotherToken)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/applications/:applicationId/comments/:commentId")
		c.SetParamNames("applicationId", "commentId")
		c.SetParamValues(id.String(), strconv.Itoa(commentId))

		route, pathParam, err := router.FindRoute(req.Method, req.URL)
		if err != nil {
			panic(err)
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParam,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			panic(err)
		}

		c, err = service.SetMyUser(c)
		if err != nil {
			panic(err)
		}

		err = service.DeleteComments(c)
		asr.NoError(err)
		asr.Equal(http.StatusForbidden, rec.Code)
	})
}
