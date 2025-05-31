package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/testutil"
	"github.com/traPtitech/Jomon/testutil/random"
	"go.uber.org/mock/gomock"
)

func TestHandlers_GetGroups(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
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
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/groups", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroups(c.Request().Context()).
			Return(groups, nil)

		require.NoError(t, h.Handlers.GetGroups(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*GroupOverview
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := lo.Map(groups, func(group *model.Group, _ int) *GroupOverview {
			return &GroupOverview{
				ID:          group.ID,
				Name:        group.Name,
				Description: group.Description,
				Budget:      group.Budget,
				CreatedAt:   group.CreatedAt,
				UpdatedAt:   group.UpdatedAt,
			}
		})
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		var groups []*model.Group

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/groups", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroups(c.Request().Context()).
			Return(groups, nil)

		require.NoError(t, h.Handlers.GetGroups(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*GroupOverview
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := []*GroupOverview{}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("FailedToGetGroups", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/groups", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		resErr := errors.New("failed to get groups")
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroups(c.Request().Context()).
			Return(nil, resErr)

		err = h.Handlers.GetGroups(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})
}

func TestHandlers_PostGroup(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
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
		// FIXME: #833
		reqBody, err := json.Marshal(&struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Budget      int    `json:"budget"`
		}{
			Name:        group.Name,
			Description: group.Description,
			Budget:      *group.Budget,
		})
		require.NoError(t, err)
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/groups", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			CreateGroup(c.Request().Context(), group.Name, group.Description, group.Budget).
			Return(group, nil)

		require.NoError(t, h.Handlers.PostGroup(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got GroupOverview
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &GroupOverview{
			ID:          group.ID,
			Name:        group.Name,
			Description: group.Description,
			Budget:      group.Budget,
			CreatedAt:   group.CreatedAt,
			UpdatedAt:   group.UpdatedAt,
		}
		testutil.RequireEqual(t, exp, &got, opts...)
	})

	t.Run("FailedWithCreateGroup", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		budget := random.Numeric(t, 1000000)

		e := echo.New()
		// FIXME: #833
		reqBody, err := json.Marshal(&struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Budget      int    `json:"budget"`
		}{
			Name:        "test",
			Description: "test",
			Budget:      budget,
		})
		require.NoError(t, err)
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/groups", bytes.NewReader(reqBody))
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
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})
}

func TestHandlers_GetGroupDetail(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
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
		path := fmt.Sprintf("/api/groups/%s", group.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		// FIXME: #822
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
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

		err = h.Handlers.GetGroupDetail(c)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		var got GroupDetail
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &GroupDetail{
			ID:          group.ID,
			Name:        group.Name,
			Description: group.Description,
			Budget:      group.Budget,
			Owners:      ownerIDs,
			Members:     memberIDs,
			CreatedAt:   group.CreatedAt,
			UpdatedAt:   group.UpdatedAt,
		}
		testutil.RequireEqual(t, exp, &got, opts...)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
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
		path := fmt.Sprintf("/api/groups/%s", group.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		// FIXME: #822
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
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

		err = h.Handlers.GetGroupDetail(c)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		var got GroupDetail
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &GroupDetail{
			ID:          group.ID,
			Name:        group.Name,
			Description: group.Description,
			Budget:      group.Budget,
			Owners:      []*uuid.UUID{},
			Members:     []*uuid.UUID{},
			CreatedAt:   group.CreatedAt,
			UpdatedAt:   group.UpdatedAt,
		}
		testutil.RequireEqual(t, exp, &got, opts...)
	})

	t.Run("FailedWithUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		invalidUUID := "invalid-uuid"
		_, resErr := uuid.Parse(invalidUUID)

		e := echo.New()
		path := fmt.Sprintf("/api/groups/%s", invalidUUID)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(invalidUUID)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		err = h.Handlers.GetGroupDetail(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("NilGroupID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		e := echo.New()
		path := fmt.Sprintf("/api/groups/%s", uuid.Nil.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(uuid.Nil.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		resErr := errors.New("invalid UUID")

		err = h.Handlers.GetGroupDetail(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("UnknownGroupID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)
		unknownGroupID := uuid.New()
		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown group id"), &resErr)

		e := echo.New()
		path := fmt.Sprintf("/api/groups/%s", unknownGroupID)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(unknownGroupID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroup(c.Request().Context(), unknownGroupID).
			Return(nil, resErr)

		err = h.Handlers.GetGroupDetail(c)
		require.Error(t, err)
		// FIXME: http.StatusNotFoundだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
	})

	t.Run("FailedToGetGroup", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
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
		path := fmt.Sprintf("/api/groups/%s", group.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroup(c.Request().Context(), group.ID).
			Return(nil, resErr)

		err = h.Handlers.GetGroupDetail(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})

	t.Run("FailedToGetOwners", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
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
		path := fmt.Sprintf("/api/groups/%s", group.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroup(c.Request().Context(), group.ID).
			Return(group, nil)

		h.Repository.MockGroupRepository.
			EXPECT().
			GetOwners(c.Request().Context(), group.ID).
			Return(nil, resErr)

		err = h.Handlers.GetGroupDetail(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})

	t.Run("FailedToGetMembers", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
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
		path := fmt.Sprintf("/api/groups/%s", group.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
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
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})
}

func TestHandlers_PutGroup(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
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
		// FIXME: #833
		reqBody, err := json.Marshal(&struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Budget      int    `json:"budget"`
		}{
			Name:        updated.Name,
			Description: updated.Description,
			Budget:      *updated.Budget,
		})
		require.NoError(t, err)
		path := fmt.Sprintf("/api/groups/%s", group.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetOwners(c.Request().Context(), group.ID).
			Return([]*model.Owner{{ID: user.ID}}, nil)
		h.Repository.MockGroupRepository.
			EXPECT().
			UpdateGroup(
				c.Request().Context(),
				group.ID, updated.Name,
				updated.Description, updated.Budget).
			Return(updated, nil)
		require.NoError(t, h.Handlers.PutGroup(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got GroupOverview
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &GroupOverview{
			ID:          updated.ID,
			Name:        updated.Name,
			Description: updated.Description,
			Budget:      updated.Budget,
			CreatedAt:   updated.CreatedAt,
			UpdatedAt:   updated.UpdatedAt,
		}
		testutil.RequireEqual(t, exp, &got, opts...)
	})

	t.Run("FailedWithUpdateGroup", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
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

		// FIXME: #833
		reqBody, err := json.Marshal(&struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Budget      int    `json:"budget"`
		}{
			Name:        updated.Name,
			Description: updated.Description,
			Budget:      *updated.Budget,
		})
		require.NoError(t, err)
		path := fmt.Sprintf("/api/groups/%s", group.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetPath("api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		resErr := errors.New("Failed to get requests.")
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetOwners(c.Request().Context(), group.ID).
			Return([]*model.Owner{{ID: user.ID}}, nil)
		h.Repository.MockGroupRepository.
			EXPECT().
			UpdateGroup(
				c.Request().Context(),
				group.ID, updated.Name, updated.Description, updated.Budget).
			Return(nil, resErr)

		err = h.Handlers.PutGroup(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})

	t.Run("FailedWithUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		invalidUUID := "invalid-uuid"
		_, resErr := uuid.Parse(invalidUUID)

		e := echo.New()
		// FIXME: #833
		reqBody, err := json.Marshal(&struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Budget      int    `json:"budget"`
		}{
			Name:        "test",
			Description: "test",
			Budget:      1000000,
		})
		require.NoError(t, err)
		path := fmt.Sprintf("/api/groups/%s", invalidUUID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(invalidUUID)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		err = h.Handlers.PutGroup(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})
}

func TestHandlers_DeleteGroup(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
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
		path := fmt.Sprintf("/api/groups/%s", group.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetOwners(c.Request().Context(), group.ID).
			Return([]*model.Owner{{ID: user.ID}}, nil)
		h.Repository.MockGroupRepository.
			EXPECT().
			DeleteGroup(c.Request().Context(), group.ID).
			Return(nil)

		require.NoError(t, h.Handlers.DeleteGroup(c))
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("FailedWithDeleteGroup", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
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
		path := fmt.Sprintf("/api/groups/%s", group.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		resErr := errors.New("Failed to get requests.")
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetOwners(c.Request().Context(), group.ID).
			Return([]*model.Owner{{ID: user.ID}}, nil)
		h.Repository.MockGroupRepository.
			EXPECT().
			DeleteGroup(c.Request().Context(), group.ID).
			Return(resErr)

		err = h.Handlers.DeleteGroup(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})

	t.Run("FailedWithUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		invalidUUID := "invalid-uuid"
		_, resErr := uuid.Parse(invalidUUID)

		e := echo.New()
		path := fmt.Sprintf("/api/groups/%s", invalidUUID)
		req := httptest.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/groups/:groupID")
		c.SetParamNames("groupID")
		c.SetParamValues(invalidUUID)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		err = h.Handlers.DeleteGroup(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})
}

func TestHandlers_PostMember(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
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
		member := []uuid.UUID{user.ID}
		reqBody, err := json.Marshal(member)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/groups/%s/members", group.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetOwners(c.Request().Context(), group.ID).
			Return([]*model.Owner{{ID: user.ID}}, nil)
		modelMember := &model.Member{
			ID: user.ID,
		}
		h.Repository.MockGroupRepository.
			EXPECT().
			AddMembers(c.Request().Context(), group.ID, []uuid.UUID{user.ID}).
			Return([]*model.Member{
				modelMember,
			}, nil)

		require.NoError(t, err)
		require.NoError(t, h.Handlers.PostMember(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []uuid.UUID
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		exp := []uuid.UUID{user.ID}
		testutil.RequireEqual(t, exp, got)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		invalidUUID := "invalid-uuid"
		_, resErr := uuid.Parse(invalidUUID)
		member := []uuid.UUID{uuid.New()}
		reqBody, err := json.Marshal(member)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/groups/%s/members", invalidUUID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(invalidUUID)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.PostMember(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		member := []uuid.UUID{uuid.Nil}
		reqBody, err := json.Marshal(member)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/groups/%s/members", uuid.Nil.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(reqBody))
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
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("UnknownGroupID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		unknownGroupID := uuid.New()
		var resErr *ent.ConstraintError
		errors.As(errors.New("unknown group id"), &resErr)

		member := []uuid.UUID{user.ID}
		reqBody, err := json.Marshal(member)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/groups/%s/members", unknownGroupID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(unknownGroupID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetOwners(c.Request().Context(), unknownGroupID).
			Return([]*model.Owner{{ID: user.ID}}, nil)
		h.Repository.MockGroupRepository.
			EXPECT().
			AddMembers(c.Request().Context(), unknownGroupID, []uuid.UUID{user.ID}).
			Return(nil, resErr)

		err = h.Handlers.PostMember(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})

	t.Run("UnknownUserID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
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
		path := fmt.Sprintf("/api/groups/%s/members", group.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetOwners(c.Request().Context(), group.ID).
			Return([]*model.Owner{{ID: user.ID}}, nil)
		h.Repository.MockGroupRepository.
			EXPECT().
			AddMembers(c.Request().Context(), group.ID, []uuid.UUID{unknownUserID}).
			Return(nil, resErr)

		err = h.Handlers.PostMember(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})
}

func TestHandlers_DeleteMember(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
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
		member := []uuid.UUID{user.ID}
		reqBody, err := json.Marshal(member)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/groups/%s/members", group.ID.String())
		req := httptest.NewRequestWithContext(
			ctx, http.MethodDelete, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			DeleteMembers(c.Request().Context(), group.ID, []uuid.UUID{user.ID}).
			Return(nil)

		require.NoError(t, h.Handlers.DeleteMember(c))
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("NilGroupUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
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
		path := fmt.Sprintf("/api/groups/%s/members", uuid.Nil.String())
		req := httptest.NewRequestWithContext(
			ctx, http.MethodDelete, path, bytes.NewReader(reqBody))
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
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("UnknownGroupID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
		unknownGroupID := uuid.New()
		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown group id"), &resErr)
		member := []uuid.UUID{user.ID}
		reqBody, err := json.Marshal(member)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/groups/%s/members", unknownGroupID.String())
		req := httptest.NewRequestWithContext(
			ctx, http.MethodDelete, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(unknownGroupID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			DeleteMembers(c.Request().Context(), unknownGroupID, []uuid.UUID{user.ID}).
			Return(resErr)

		err = h.Handlers.DeleteMember(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})

	t.Run("UnknownMemberID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
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
		path := fmt.Sprintf("/api/groups/%s/members", group.ID.String())
		req := httptest.NewRequestWithContext(
			ctx, http.MethodDelete, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetOwners(c.Request().Context(), group.ID).
			Return([]*model.Owner{{ID: user.ID}}, nil)
		h.Repository.MockGroupRepository.
			EXPECT().
			DeleteMembers(c.Request().Context(), group.ID, []uuid.UUID{unknownUserID}).
			Return(resErr)

		err = h.Handlers.DeleteMember(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})

	t.Run("InvalidGroupUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		invalidUUID := "invalid-uuid"
		_, resErr := uuid.Parse(invalidUUID)
		member := []uuid.UUID{uuid.New()}
		reqBody, err := json.Marshal(member)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/groups/%s/members", invalidUUID)
		req := httptest.NewRequestWithContext(
			ctx, http.MethodDelete, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/members")
		c.SetParamNames("groupID")
		c.SetParamValues(invalidUUID)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.DeleteMember(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})
}

func TestHandlers_PostOwner(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
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
		owner := []uuid.UUID{user.ID}
		reqBody, err := json.Marshal(owner)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/groups/%s/owners", group.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/owners")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())
		c.Set(loginUserKey, user)

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

		require.NoError(t, h.Handlers.PostOwner(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []uuid.UUID
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		exp := []uuid.UUID{user.ID}
		testutil.RequireEqual(t, exp, got)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		invalidUUID := "invalid-uuid"
		_, resErr := uuid.Parse(invalidUUID)
		owner := []string{invalidUUID}
		reqBody, err := json.Marshal(owner)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/groups/%s/owners", invalidUUID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/owners")
		c.SetParamNames("groupID")
		c.SetParamValues(invalidUUID)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.PostOwner(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		owner := []uuid.UUID{uuid.New()}
		reqBody, err := json.Marshal(owner)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/groups/%s/owners", uuid.Nil.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(reqBody))
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
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("UnknownGroupID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
		unknownGroupID := uuid.New()
		var resErr *ent.ConstraintError
		errors.As(errors.New("unknown group id"), &resErr)
		owner := []uuid.UUID{user.ID}
		reqBody, err := json.Marshal(owner)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/groups/%s/owners", unknownGroupID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/owner")
		c.SetParamNames("groupID")
		c.SetParamValues(unknownGroupID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			AddOwners(c.Request().Context(), unknownGroupID, []uuid.UUID{user.ID}).
			Return(nil, resErr)

		err = h.Handlers.PostOwner(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("UnknownUserID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
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
		path := fmt.Sprintf("/api/groups/%s/owners", group.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/owners")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			AddOwners(c.Request().Context(), group.ID, []uuid.UUID{unknownUserID}).
			Return(nil, resErr)

		err = h.Handlers.PostOwner(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})
}

func TestHandlers_DeleteOwner(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
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
		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
		owner := []uuid.UUID{user.ID}
		reqBody, err := json.Marshal(owner)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/groups/%s/owners", group.ID.String())
		req := httptest.NewRequestWithContext(
			ctx, http.MethodDelete, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/owners")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			DeleteOwners(c.Request().Context(), group.ID, []uuid.UUID{user.ID}).
			Return(nil)

		require.NoError(t, h.Handlers.DeleteOwner(c))
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("NilGroupUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		owner := []uuid.UUID{uuid.Nil}
		reqBody, err := json.Marshal(owner)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/groups/%s/owners", uuid.Nil.String())
		req := httptest.NewRequestWithContext(
			ctx, http.MethodDelete, path, bytes.NewReader(reqBody))
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
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("UnknownGroupID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
		unknownGroupID := uuid.New()
		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown group id"), &resErr)
		owner := []uuid.UUID{user.ID}
		reqBody, err := json.Marshal(owner)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/groups/%s/owners", unknownGroupID)
		req := httptest.NewRequestWithContext(
			ctx, http.MethodDelete, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/owners")
		c.SetParamNames("groupID")
		c.SetParamValues(unknownGroupID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			DeleteOwners(c.Request().Context(), unknownGroupID, []uuid.UUID{user.ID}).
			Return(resErr)

		err = h.Handlers.DeleteOwner(c)
		require.Error(t, err)
		require.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
	})

	t.Run("UnknownOwnerID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
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
		path := fmt.Sprintf("/api/groups/%s/owners", group.ID.String())
		req := httptest.NewRequestWithContext(
			ctx, http.MethodDelete, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/owners")
		c.SetParamNames("groupID")
		c.SetParamValues(group.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetOwners(c.Request().Context(), group.ID).
			Return([]*model.Owner{{ID: user.ID}}, nil)
		h.Repository.MockGroupRepository.
			EXPECT().
			DeleteOwners(c.Request().Context(), group.ID, []uuid.UUID{unknownUserID}).
			Return(resErr)

		err = h.Handlers.DeleteOwner(c)
		require.Error(t, err)
		// FIXME: http.StatusNotFoundだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
	})

	t.Run("InvalidGroupUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		invalidUUID := "invalid-uuid"
		_, resErr := uuid.Parse(invalidUUID)
		owner := []uuid.UUID{uuid.New()}
		reqBody, err := json.Marshal(owner)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/groups/%s/owners", invalidUUID)
		req := httptest.NewRequestWithContext(
			ctx, http.MethodDelete, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/groups/:groupID/owners")
		c.SetParamNames("groupID")
		c.SetParamValues(invalidUUID)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.DeleteOwner(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})
}
