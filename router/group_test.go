package router

import (
	"bytes"
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
		resErr := errors.New("failed to get groups")
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroups(c.Request().Context()).
			Return(nil, resErr)

		err = h.Handlers.GetGroups(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
		}
	})
}

func TestHandlers_PostGroup(t *testing.T) {
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

		e := echo.New()
		reqBody := fmt.Sprintf(
			`{"name":"%s","description":"%s","budget":%d}`,
			group.Name, group.Description, *group.Budget)
		req, err := http.NewRequest(
			http.MethodPost,
			"/api/groups",
			strings.NewReader(reqBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			CreateGroup(c.Request().Context(), group.Name, group.Description, group.Budget).
			Return(group, nil)

		res := &GroupOverview{
			ID:          group.ID,
			Name:        group.Name,
			Description: group.Description,
			Budget:      group.Budget,
			CreatedAt:   group.CreatedAt,
			UpdatedAt:   group.UpdatedAt,
		}

		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PostGroup(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("FailedWithCreateGroup", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		budget := random.Numeric(t, 1000000)

		e := echo.New()
		reqBody := fmt.Sprintf(`{"name":"test","description":"test","budget":%d}`, budget)
		req, err := http.NewRequest(
			http.MethodPost,
			"/api/groups",
			strings.NewReader(reqBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		resErr := errors.New("failed to create group")
		h.Repository.MockGroupRepository.
			EXPECT().
			CreateGroup(c.Request().Context(), "test", "test", &budget).
			Return(nil, resErr)

		err = h.Handlers.PostGroup(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
		}
	})
}

func TestHandlers_GetGroupDetail(t *testing.T) {
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
		user3 := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       false,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		user4 := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       false,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		owner1 := model.Owner{ID: user1.ID}
		owner2 := model.Owner{ID: user2.ID}
		owners := []*model.Owner{&owner1, &owner2}
		ownerIDs := []*uuid.UUID{&user1.ID, &user2.ID}

		member1 := model.Member{ID: user3.ID}
		member2 := model.Member{ID: user4.ID}
		members := []*model.Member{&member1, &member2}
		memberIDs := []*uuid.UUID{&member1.ID, &member2.ID}

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/groups/%s", group.ID), nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroup(c.Request().Context(), group.ID).
			Return(group, nil)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetOwners(c.Request().Context(), group.ID).
			Return(owners, nil)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetMembers(c.Request().Context(), group.ID).
			Return(members, nil)

		res := &GroupDetail{
			ID:          group.ID,
			Name:        group.Name,
			Description: group.Description,
			Budget:      group.Budget,
			Owners:      ownerIDs,
			Members:     memberIDs,
			CreatedAt:   group.CreatedAt,
			UpdatedAt:   group.UpdatedAt,
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		err = h.Handlers.GetGroupDetail(c)
		if assert.NoError(t, err) {
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

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/groups/%s", group.ID), nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroup(c.Request().Context(), group.ID).
			Return(group, nil)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetOwners(c.Request().Context(), group.ID).
			Return([]*model.Owner{}, nil)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetMembers(c.Request().Context(), group.ID).
			Return([]*model.Member{}, nil)

		res := &GroupDetail{
			ID:          group.ID,
			Name:        group.Name,
			Description: group.Description,
			Budget:      group.Budget,
			Owners:      []*uuid.UUID{},
			Members:     []*uuid.UUID{},
			CreatedAt:   group.CreatedAt,
			UpdatedAt:   group.UpdatedAt,
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		err = h.Handlers.GetGroupDetail(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("FailedWithUUID", func(t *testing.T) {
		t.Parallel()

		invalidUUID := "invalid-uuid"
		_, resErr := uuid.Parse(invalidUUID)

		ctrl := gomock.NewController(t)
		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/api/groups/invalid-uuid", nil)
		require.NoError(t, err)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues("invalid-uuid")

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		err = h.Handlers.GetGroupDetail(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("NilGroupID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/groups/%s", uuid.Nil.String()),
			nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(uuid.Nil.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		resErr := errors.New("invalid UUID")

		err = h.Handlers.GetGroupDetail(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("UnknownGroupID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		unknownGroupID := uuid.New()
		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown group id"), &resErr)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/groups/%s", unknownGroupID),
			nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(unknownGroupID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroup(c.Request().Context(), unknownGroupID).
			Return(nil, resErr)

		err = h.Handlers.GetGroupDetail(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
		}
	})

	t.Run("FailedToGetGroup", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		resErr := errors.New("failed to get groups")

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

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/groups/%s", group.ID), nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroup(c.Request().Context(), group.ID).
			Return(nil, resErr)

		err = h.Handlers.GetGroupDetail(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
		}
	})

	t.Run("FailedToGetOwners", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		resErr := errors.New("failed to get owners")

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

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/groups/%s", group.ID), nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroup(c.Request().Context(), group.ID).
			Return(group, nil)

		h.Repository.MockGroupRepository.
			EXPECT().
			GetOwners(c.Request().Context(), group.ID).
			Return(nil, resErr)

		err = h.Handlers.GetGroupDetail(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
		}
	})

	t.Run("FailedToGetMembers", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		resErr := errors.New("failed to get members")

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
		owner := model.Owner{ID: user.ID}
		owners := []*model.Owner{&owner}

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/groups/%s", group.ID), nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroup(c.Request().Context(), group.ID).
			Return(group, nil)

		h.Repository.MockGroupRepository.
			EXPECT().
			GetOwners(c.Request().Context(), group.ID).
			Return(owners, nil)

		h.Repository.MockGroupRepository.
			EXPECT().
			GetMembers(c.Request().Context(), group.ID).
			Return(nil, resErr)

		err = h.Handlers.GetGroupDetail(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
		}
	})
}

func TestHandlers_PutGroup(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()
		date2 := time.Now().Add(time.Hour)

		budget := random.Numeric(t, 1000000)
		budget2 := random.Numeric(t, 1000000)

		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		updated := &model.Group{
			ID:          group.ID,
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget2,
			CreatedAt:   date2,
			UpdatedAt:   date2,
		}

		e := echo.New()
		reqBody := fmt.Sprintf(
			`{"name":"%s","description":"%s","budget":%d}`,
			updated.Name, updated.Description, *updated.Budget)
		req, err := http.NewRequest(
			http.MethodPut,
			fmt.Sprintf("/api/groups/%s", group.ID.String()),
			strings.NewReader(reqBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			UpdateGroup(
				c.Request().Context(),
				group.ID, updated.Name,
				updated.Description, updated.Budget).
			Return(updated, nil)

		res := &GroupOverview{
			ID:          updated.ID,
			Name:        updated.Name,
			Description: updated.Description,
			Budget:      updated.Budget,
			CreatedAt:   updated.CreatedAt,
			UpdatedAt:   updated.UpdatedAt,
		}

		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PutGroup(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("FailedWithUpdateGroup", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()
		date2 := time.Now().Add(time.Hour)

		budget := random.Numeric(t, 1000000)
		budget2 := random.Numeric(t, 1000000)

		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		updated := &model.Group{
			ID:          group.ID,
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget2,
			CreatedAt:   date2,
			UpdatedAt:   date2,
		}

		e := echo.New()
		reqBody := fmt.Sprintf(
			`{"name":"%s","description":"%s","budget":%d}`,
			updated.Name, updated.Description, *updated.Budget)
		req, err := http.NewRequest(
			http.MethodPut,
			fmt.Sprintf("/api/groups/%s", group.ID.String()),
			strings.NewReader(reqBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		resErr := errors.New("Failed to get requests.")
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			UpdateGroup(
				c.Request().Context(),
				group.ID, updated.Name, updated.Description, updated.Budget).
			Return(nil, resErr)

		err = h.Handlers.PutGroup(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
		}
	})

	t.Run("FailedWithUUID", func(t *testing.T) {
		t.Parallel()

		invalidUUID := "invalid-uuid"
		_, resErr := uuid.Parse(invalidUUID)

		ctrl := gomock.NewController(t)
		e := echo.New()
		req, err := http.NewRequest(
			http.MethodPut,
			"/api/groups/invalid-uuid",
			strings.NewReader(`{"name":"test","description":"test","budget":1000000}`))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues("invalid-uuid")

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		err = h.Handlers.PutGroup(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})
}

func TestHandlers_DeleteGroup(t *testing.T) {
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

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/groups/%s", group.ID.String()),
			nil)
		require.NoError(t, err)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			DeleteGroup(c.Request().Context(), group.ID).
			Return(nil)

		if assert.NoError(t, h.Handlers.DeleteGroup(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedWithDeleteGroup", func(t *testing.T) {
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

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/groups/%s", group.ID.String()),
			nil)
		require.NoError(t, err)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		resErr := errors.New("Failed to get requests.")
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			DeleteGroup(c.Request().Context(), group.ID).
			Return(resErr)

		err = h.Handlers.DeleteGroup(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
		}
	})

	t.Run("FailedWithUUID", func(t *testing.T) {
		t.Parallel()

		invalidUUID := "invalid-uuid"
		_, resErr := uuid.Parse(invalidUUID)

		ctrl := gomock.NewController(t)
		e := echo.New()
		req, err := http.NewRequest(http.MethodDelete, "/api/groups/invalid-uuid", nil)
		require.NoError(t, err)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues("invalid-uuid")

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		err = h.Handlers.DeleteGroup(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})
}

func TestHandlers_PostMember(t *testing.T) {
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

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		member := []uuid.UUID{user.ID}
		reqBody, err := json.Marshal(member)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("/api/groups/%s/members", group.ID.String()),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		modelMember := &model.Member{
			ID: user.ID,
		}
		h.Repository.MockGroupRepository.
			EXPECT().
			AddMembers(c.Request().Context(), group.ID, []uuid.UUID{user.ID}).
			Return([]*model.Member{
				modelMember,
			}, nil)

		res := []uuid.UUID{modelMember.ID}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PostMember(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		member := []uuid.UUID{uuid.New()}
		reqBody, err := json.Marshal(member)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodPost,
			"/api/groups/hoge/members",
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues("hoge")

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		_, resErr := uuid.Parse(c.Param("groupID"))

		err = h.Handlers.PostMember(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		member := []uuid.UUID{uuid.Nil}
		reqBody, err := json.Marshal(member)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("/api/groups/%s/members", uuid.Nil.String()),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(uuid.Nil.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		resErr := errors.New("invalid UUID")

		err = h.Handlers.PostMember(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("UnknownGroupID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
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
		var resErr *ent.ConstraintError
		errors.As(errors.New("unknown group id"), &resErr)

		member := []uuid.UUID{user.ID}
		reqBody, err := json.Marshal(member)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("/api/groups/%s/members", unknownGroupID),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(unknownGroupID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			AddMembers(c.Request().Context(), unknownGroupID, []uuid.UUID{user.ID}).
			Return(nil, resErr)

		err = h.Handlers.PostMember(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
		}
	})

	t.Run("UnknownUserID", func(t *testing.T) {
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

		unknownUserID := uuid.New()
		var resErr *ent.ConstraintError
		errors.As(errors.New("unknown user id"), &resErr)

		member := []uuid.UUID{unknownUserID}
		reqBody, err := json.Marshal(member)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("/api/groups/%s/members", group.ID.String()),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			AddMembers(c.Request().Context(), group.ID, []uuid.UUID{unknownUserID}).
			Return(nil, resErr)

		err = h.Handlers.PostMember(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
		}
	})
}

func TestHandlers_DeleteMember(t *testing.T) {
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

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		member := []uuid.UUID{user.ID}
		reqBody, err := json.Marshal(member)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/groups/%s/members", group.ID.String()),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			DeleteMembers(c.Request().Context(), group.ID, []uuid.UUID{user.ID}).
			Return(nil)

		if assert.NoError(t, h.Handlers.DeleteMember(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("NilGroupUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		member := []uuid.UUID{user.ID}
		reqBody, err := json.Marshal(member)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/groups/%s/members", uuid.Nil.String()),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(uuid.Nil.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		resErr := errors.New("invalid UUID")

		err = h.Handlers.DeleteMember(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("UnknownGroupID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
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
		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown group id"), &resErr)

		member := []uuid.UUID{user.ID}
		reqBody, err := json.Marshal(member)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/groups/%s/members", unknownGroupID.String()),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(unknownGroupID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			DeleteMembers(c.Request().Context(), unknownGroupID, []uuid.UUID{user.ID}).
			Return(resErr)

		err = h.Handlers.DeleteMember(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
		}
	})

	t.Run("UnknownMemberID", func(t *testing.T) {
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

		unknownUserID := uuid.New()
		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown member id"), &resErr)

		member := []uuid.UUID{unknownUserID}
		reqBody, err := json.Marshal(member)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/groups/%s/members", group.ID.String()),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			DeleteMembers(c.Request().Context(), group.ID, []uuid.UUID{unknownUserID}).
			Return(resErr)

		err = h.Handlers.DeleteMember(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
		}
	})

	t.Run("InvalidGroupUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		invID := "po"

		_, resErr := uuid.Parse(invID)

		member := []uuid.UUID{uuid.New()}
		reqBody, err := json.Marshal(member)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/groups/%s/members", invID),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(invID)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.DeleteMember(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})
}

func TestHandlers_PostOwner(t *testing.T) {
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

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		owner := []uuid.UUID{user.ID}
		reqBody, err := json.Marshal(owner)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("/api/groups/%s/owners", group.ID.String()),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/owners")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		modelOwner := &model.Owner{
			ID: user.ID,
		}
		h.Repository.MockGroupRepository.
			EXPECT().
			AddOwners(c.Request().Context(), group.ID, owner).
			Return([]*model.Owner{
				modelOwner,
			}, nil)

		res := owner
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PostOwner(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		owner := []string{"hoge"}
		reqBody, err := json.Marshal(owner)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodPost,
			"/api/groups/hoge/owners",
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/owners")
		c.SetParamNames("groupID")
		c.SetParamValues("hoge")

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		_, resErr := uuid.Parse(c.Param("groupID"))

		err = h.Handlers.PostOwner(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		owner := []uuid.UUID{uuid.New()}
		reqBody, err := json.Marshal(owner)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("/api/groups/%s/owners", uuid.Nil.String()),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/owners")
		c.SetParamNames("groupID")
		c.SetParamValues(uuid.Nil.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		resErr := errors.New("invalid UUID")

		err = h.Handlers.PostOwner(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("UnknownGroupID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
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
		var resErr *ent.ConstraintError
		errors.As(errors.New("unknown group id"), &resErr)

		owner := []uuid.UUID{user.ID}
		reqBody, err := json.Marshal(owner)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("/api/groups/%s/owners", unknownGroupID),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/owner")
		c.SetParamNames("groupID")
		c.SetParamValues(unknownGroupID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			AddOwners(c.Request().Context(), unknownGroupID, []uuid.UUID{user.ID}).
			Return(nil, resErr)

		err = h.Handlers.PostOwner(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("UnknownUserID", func(t *testing.T) {
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

		unknownUserID := uuid.New()
		var resErr *ent.ConstraintError
		errors.As(errors.New("unknown user id"), &resErr)

		owner := []uuid.UUID{unknownUserID}
		reqBody, err := json.Marshal(owner)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("/api/groups/%s/owners", group.ID.String()),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/owners")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			AddOwners(c.Request().Context(), group.ID, []uuid.UUID{unknownUserID}).
			Return(nil, resErr)

		err = h.Handlers.PostOwner(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})
}

func TestHandlers_DeleteOwner(t *testing.T) {
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

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		owner := []uuid.UUID{user.ID}
		reqBody, err := json.Marshal(owner)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/groups/%s/owners", group.ID.String()),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/owners")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			DeleteOwners(c.Request().Context(), group.ID, []uuid.UUID{user.ID}).
			Return(nil)

		if assert.NoError(t, h.Handlers.DeleteOwner(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("NilGroupUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		owner := []uuid.UUID{uuid.Nil}
		reqBody, err := json.Marshal(owner)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/groups/%s/owners", uuid.Nil.String()),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/owners")
		c.SetParamNames("groupID")
		c.SetParamValues(uuid.Nil.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		resErr := errors.New("invalid UUID")

		err = h.Handlers.DeleteOwner(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("UnknownGroupID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
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
		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown group id"), &resErr)

		owner := []uuid.UUID{user.ID}
		reqBody, err := json.Marshal(owner)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/groups/%s/owners", unknownGroupID),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/owners")
		c.SetParamNames("groupID")
		c.SetParamValues(unknownGroupID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			DeleteOwners(c.Request().Context(), unknownGroupID, []uuid.UUID{user.ID}).
			Return(resErr)

		err = h.Handlers.DeleteOwner(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
		}
	})

	t.Run("UnknownOwnerID", func(t *testing.T) {
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

		unknownUserID := uuid.New()
		reqBody, err := json.Marshal([]uuid.UUID{unknownUserID})
		require.NoError(t, err)
		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown owner id"), &resErr)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/groups/%s/owners", group.ID.String()),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/owners")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			DeleteOwners(c.Request().Context(), group.ID, []uuid.UUID{unknownUserID}).
			Return(resErr)

		err = h.Handlers.DeleteOwner(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
		}
	})

	t.Run("InvalidGroupUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		invID := "po"
		owner := []uuid.UUID{uuid.New()}
		reqBody, err := json.Marshal(owner)
		require.NoError(t, err)

		_, resErr := uuid.Parse(invID)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/groups/%s/owners", invID),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/owners")
		c.SetParamNames("groupID")
		c.SetParamValues(invID)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.DeleteOwner(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})
}
