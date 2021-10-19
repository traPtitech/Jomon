package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/testutil/random"
)

// TODO: 直す
func TestHandlers_GetGroups(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		budget1 := random.Numeric(t, 1000000)
		budget2 := random.Numeric(t, 1000000)

		group1 := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget1,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		group2 := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget2,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		groups := []*model.Group{group1, group2}

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/api/groups", nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroups(c.Request().Context()).
			Return(groups, nil)

		resOverview := []*GroupOverview{}
		for _, group := range groups {
			resOverview = append(resOverview, &GroupOverview{
				ID:          group.ID,
				Name:        group.Name,
				Description: group.Description,
				Budget:      group.Budget,
				CreatedAt:   group.CreatedAt,
				UpdatedAt:   group.UpdatedAt,
			})
		}
		resBody, err := json.Marshal(resOverview)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetGroups(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		groups := []*model.Group{}

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/api/groups", nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroups(c.Request().Context()).
			Return(groups, nil)

		resOverview := []*GroupOverview{}
		resBody, err := json.Marshal(resOverview)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetGroups(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("FailedToGetGroups", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/api/groups", nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		mocErr := errors.New("failed to get groups")
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroups(c.Request().Context()).
			Return(nil, mocErr)

		err = h.Handlers.GetGroups(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
		}
	})
}

func TestHandlers_GetMembers(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {

		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		budget := random.Numeric(t, 1000000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		user1 := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		user2 := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		members := []*model.User{user1, user2}
		memberIDs := []uuid.UUID{user1.ID, user2.ID}

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/groups/%s/members", group.ID.String()), nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetMembers(c.Request().Context(), group.ID).
			Return(members, nil)
		res := &MemberResponse{
			ID: memberIDs,
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetMembers(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("Success2", func(t *testing.T) {

		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		budget := random.Numeric(t, 1000000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		members := []*model.User{}

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/groups/%s/members", group.ID.String()), nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetMembers(c.Request().Context(), group.ID).
			Return(members, nil)
		res := &MemberResponse{
			ID: nil,
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetMembers(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/api/groups/hoge/members", nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues("hoge")

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		_, resErr := uuid.Parse(c.Param("groupID"))

		err = h.Handlers.GetMembers(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/groups/%s/members", uuid.Nil.String()), nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(uuid.Nil.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		resErr := errors.New("invalid UUID")

		err = h.Handlers.GetMembers(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("UnknownID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		unknownID := uuid.New()

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/groups/%s/members", unknownID.String()), nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(unknownID.String())

		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown id"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetMembers(c.Request().Context(), unknownID).
			Return(nil, resErr)

		err = h.Handlers.GetMembers(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
		}
	})
}

/*
func TestHandlers_PostMember(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		budget := random.Numeric(t, 1000000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		ctx := context.Background()
		th.Repository.MockGroupRepository.
			EXPECT().
			CreateMember(ctx, group.ID, user.ID).
			Return(&model.Member{
				ID: user.ID,
			}, nil)

		req := Member{
			ID: user.ID,
		}

		var resBody Member
		path := fmt.Sprintf("/api/groups/%s/members", group.ID.String())
		statusCode, _ := th.doRequest(t, echo.POST, path, &req, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, user.ID, resBody.ID)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		req := Member{
			ID: user.ID,
		}

		path := "/api/groups/hoge/members" // Invalid UUID
		statusCode, _ := th.doRequest(t, echo.POST, path, &req, nil)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		req := Member{
			ID: user.ID,
		}

		path := fmt.Sprintf("/api/groups/%s/members", uuid.Nil.String())
		statusCode, _ := th.doRequest(t, echo.POST, path, &req, nil)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("UnknownGroupID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		unknownGroupID := uuid.New()
		ctx := context.Background()
		var e *ent.ConstraintError
		errors.As(errors.New("unknown group id"), &e)

		th.Repository.MockGroupRepository.
			EXPECT().
			CreateMember(ctx, unknownGroupID, user.ID).
			Return(nil, e)

		req := Member{
			ID: user.ID,
		}

		path := fmt.Sprintf("/api/groups/%s/members", unknownGroupID.String())
		statusCode, _ := th.doRequest(t, echo.POST, path, &req, nil)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("UnknownUserID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		budget := random.Numeric(t, 1000000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		unknownUserID := uuid.New()
		ctx := context.Background()
		var e *ent.ConstraintError
		errors.As(errors.New("unknown user id"), &e)

		th.Repository.MockGroupRepository.
			EXPECT().
			CreateMember(ctx, group.ID, unknownUserID).
			Return(nil, e)

		req := Member{
			ID: unknownUserID,
		}

		path := fmt.Sprintf("/api/groups/%s/members", group.ID.String())
		statusCode, _ := th.doRequest(t, echo.POST, path, &req, nil)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})
}

func TestHandlers_DeleteMember(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		budget := random.Numeric(t, 1000000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		ctx := context.Background()
		th.Repository.MockGroupRepository.
			EXPECT().
			DeleteMember(ctx, group.ID, user.ID).
			Return(nil)

		req := Member{
			ID: user.ID,
		}

		path := fmt.Sprintf("/api/groups/%s/members", group.ID.String())
		statusCode, _ := th.doRequest(t, echo.DELETE, path, &req, nil)
		assert.Equal(t, http.StatusOK, statusCode)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		req := Member{
			ID: user.ID,
		}

		path := "/api/groups/hoge/members" // Invalid UUID
		statusCode, _ := th.doRequest(t, echo.DELETE, path, &req, nil)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		req := Member{
			ID: user.ID,
		}

		path := fmt.Sprintf("/api/groups/%s/members", uuid.Nil.String())
		statusCode, _ := th.doRequest(t, echo.DELETE, path, &req, nil)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("UnknownGroupID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		req := Member{
			ID: user.ID,
		}

		ctx := context.Background()
		unknownGroupID := uuid.New()
		var e *ent.NotFoundError
		errors.As(errors.New("unknown group id"), &e)

		th.Repository.MockGroupRepository.
			EXPECT().
			DeleteMember(ctx, unknownGroupID, user.ID).
			Return(e)

		path := fmt.Sprintf("/api/groups/%s/members", unknownGroupID.String())
		statusCode, _ := th.doRequest(t, echo.DELETE, path, &req, nil)
		assert.Equal(t, http.StatusNotFound, statusCode)
	})

	t.Run("UnknownMemberID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		budget := random.Numeric(t, 1000000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		unknownMemberID := uuid.New()
		ctx := context.Background()
		var e *ent.NotFoundError
		errors.As(errors.New("unknown member id"), &e)

		th.Repository.MockGroupRepository.
			EXPECT().
			DeleteMember(ctx, group.ID, unknownMemberID).
			Return(e)

		req := Member{
			ID: unknownMemberID,
		}

		path := fmt.Sprintf("/api/groups/%s/members", group.ID.String())
		statusCode, _ := th.doRequest(t, echo.DELETE, path, &req, nil)
		assert.Equal(t, http.StatusNotFound, statusCode)
	})

}
*/
