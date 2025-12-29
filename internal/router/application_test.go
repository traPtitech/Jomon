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
	"github.com/traPtitech/Jomon/internal/ent"
	"github.com/traPtitech/Jomon/internal/model"
	"github.com/traPtitech/Jomon/internal/nulltime"
	"github.com/traPtitech/Jomon/internal/testutil"
	"github.com/traPtitech/Jomon/internal/testutil/random"
	"go.uber.org/mock/gomock"
)

func modelTagToTagOverview(t *model.Tag) *TagResponse {
	return &TagResponse{
		ID:        t.ID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func modelApplicationTargetDetailToTargetOverview(
	t *model.ApplicationTargetDetail,
) *TargetOverview {
	return &TargetOverview{
		ID:        t.ID,
		Target:    t.Target,
		Amount:    t.Amount,
		CreatedAt: t.CreatedAt,
	}
}

func modelApplicationStatusToStatusResponseOverview(
	s *model.ApplicationStatus,
) *StatusResponseOverview {
	return &StatusResponseOverview{
		CreatedBy: s.CreatedBy,
		Status:    s.Status,
		CreatedAt: s.CreatedAt,
	}
}

func modelCommentToCommentDetail(c *model.Comment) *CommentDetail {
	return &CommentDetail{
		ID:        c.ID,
		User:      c.User,
		Comment:   c.Comment,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// FIXME: この処理はapplication.goにも書かれてある
func modelApplicationDetailToApplicationResponse(
	r *model.ApplicationDetail,
) *ApplicationDetailResponse {
	return &ApplicationDetailResponse{
		ApplicationResponse: ApplicationResponse{
			ID:        r.ID,
			Status:    r.Status,
			CreatedBy: r.CreatedBy,
			Title:     r.Title,
			Content:   r.Content,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
			Targets: lo.Map(
				r.Targets,
				func(m *model.ApplicationTargetDetail, _ int) *TargetOverview {
					return modelApplicationTargetDetailToTargetOverview(m)
				},
			),
			Tags: lo.Map(r.Tags, func(m *model.Tag, _ int) *TagResponse {
				return modelTagToTagOverview(m)
			}),
		},
		Statuses: lo.Map(
			r.Statuses,
			func(m *model.ApplicationStatus, _ int) *StatusResponseOverview {
				return modelApplicationStatusToStatusResponseOverview(m)
			},
		),
		Comments: lo.Map(r.Comments, func(m *model.Comment, _ int) *CommentDetail {
			return modelCommentToCommentDetail(m)
		}),
		Files: r.Files,
	}
}

func TestHandlers_GetApplications(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date1 := time.Now()
		date2 := date1.Add(time.Hour)
		application1 := &model.ApplicationResponse{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date1,
			UpdatedAt: date1,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses:  []*model.ApplicationStatus{},
		}
		application2 := &model.ApplicationResponse{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date2,
			UpdatedAt: date2,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses:  []*model.ApplicationStatus{},
		}
		applications := []*model.ApplicationResponse{application2, application1}

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/applications", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplications(c.Request().Context(), model.ApplicationQuery{
				Limit:  100,
				Offset: 0,
			}).
			Return(applications, nil)

		require.NoError(t, h.Handlers.GetApplications(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*ApplicationResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := []*ApplicationResponse{
			{
				ID:        application2.ID,
				Status:    application2.Status,
				CreatedAt: application2.CreatedAt,
				UpdatedAt: application2.UpdatedAt,
				CreatedBy: application2.CreatedBy,
				Title:     application2.Title,
				Content:   application2.Content,
				Targets:   []*TargetOverview{},
				Tags:      []*TagResponse{},
			},
			{
				ID:        application1.ID,
				Status:    application1.Status,
				CreatedAt: application1.CreatedAt,
				UpdatedAt: application1.UpdatedAt,
				CreatedBy: application1.CreatedBy,
				Title:     application1.Title,
				Content:   application1.Content,
				Targets:   []*TargetOverview{},
			},
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		applications := []*model.ApplicationResponse{}

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/applications", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplications(c.Request().Context(), model.ApplicationQuery{
				Limit:  100,
				Offset: 0,
			}).
			Return(applications, nil)

		require.NoError(t, h.Handlers.GetApplications(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*ApplicationResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		require.Empty(t, got)
	})

	t.Run("Success3", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date1 := time.Now()
		application1 := &model.ApplicationResponse{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date1,
			UpdatedAt: date1,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses:  []*model.ApplicationStatus{},
		}
		applications := []*model.ApplicationResponse{application1}

		e := echo.New()
		status := "submitted"
		path := fmt.Sprintf("/api/applications?status=%s", status)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplications(c.Request().Context(), model.ApplicationQuery{
				Status: &status,
				Limit:  100,
				Offset: 0,
			}).
			Return(applications, nil)

		require.NoError(t, h.Handlers.GetApplications(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*ApplicationResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := []*ApplicationResponse{
			{
				ID:        application1.ID,
				Status:    application1.Status,
				CreatedAt: application1.CreatedAt,
				UpdatedAt: application1.UpdatedAt,
				CreatedBy: application1.CreatedBy,
				Title:     application1.Title,
				Content:   application1.Content,
				Tags:      []*TagResponse{},
				Targets:   []*TargetOverview{},
			},
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("Success4", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date1 := time.Now()
		date2str := date1.Add(time.Hour).Format("2006-01-02")
		date2, err := nulltime.ParseDate(date2str)
		require.NoError(t, err)
		application1 := &model.ApplicationResponse{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date1,
			UpdatedAt: date1,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses:  []*model.ApplicationStatus{},
		}
		applications := []*model.ApplicationResponse{application1}

		e := echo.New()
		path := fmt.Sprintf("/api/applications?until=%s", date2str)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplications(c.Request().Context(), model.ApplicationQuery{
				Until:  date2,
				Limit:  100,
				Offset: 0,
			}).
			Return(applications, nil)

		require.NoError(t, h.Handlers.GetApplications(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*ApplicationResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := []*ApplicationResponse{
			{
				ID:        application1.ID,
				Status:    application1.Status,
				CreatedAt: application1.CreatedAt,
				UpdatedAt: application1.UpdatedAt,
				CreatedBy: application1.CreatedBy,
				Title:     application1.Title,
				Content:   application1.Content,
				Tags:      []*TagResponse{},
				Targets:   []*TargetOverview{},
			},
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("Success5", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date1 := time.Now()
		date2str := date1.Add(-time.Hour).Format("2006-01-02")
		date2, err := nulltime.ParseDate(date2str)
		require.NoError(t, err)
		application1 := &model.ApplicationResponse{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date1,
			UpdatedAt: date1,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses:  []*model.ApplicationStatus{},
		}
		applications := []*model.ApplicationResponse{application1}

		e := echo.New()
		path := fmt.Sprintf("/api/applications?since=%s", date2str)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplications(c.Request().Context(), model.ApplicationQuery{
				Since:  date2,
				Limit:  100,
				Offset: 0,
			}).
			Return(applications, nil)
		require.NoError(t, h.Handlers.GetApplications(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*ApplicationResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := []*ApplicationResponse{
			{
				ID:        application1.ID,
				Status:    application1.Status,
				CreatedAt: application1.CreatedAt,
				UpdatedAt: application1.UpdatedAt,
				CreatedBy: application1.CreatedBy,
				Title:     application1.Title,
				Content:   application1.Content,
				Targets:   []*TargetOverview{},
				Tags:      []*TagResponse{},
			},
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("Success6", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date1 := time.Now()
		tag1 := model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 10),
			CreatedAt: date1,
			UpdatedAt: date1,
		}
		tag1ov := TagResponse{
			ID:        tag1.ID,
			Name:      tag1.Name,
			CreatedAt: tag1.CreatedAt,
			UpdatedAt: tag1.UpdatedAt,
		}
		application1 := &model.ApplicationResponse{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date1,
			UpdatedAt: date1,
			Tags:      []*model.Tag{&tag1},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses:  []*model.ApplicationStatus{},
		}
		applications := []*model.ApplicationResponse{application1}

		e := echo.New()
		path := fmt.Sprintf("/api/applications?tag=%s", tag1.Name)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplications(c.Request().Context(), model.ApplicationQuery{
				Tag:    &tag1.Name,
				Limit:  100,
				Offset: 0,
			}).
			Return(applications, nil)
		require.NoError(t, h.Handlers.GetApplications(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*ApplicationResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := []*ApplicationResponse{
			{
				ID:        application1.ID,
				Status:    application1.Status,
				CreatedAt: application1.CreatedAt,
				UpdatedAt: application1.UpdatedAt,
				CreatedBy: application1.CreatedBy,
				Title:     application1.Title,
				Content:   application1.Content,
				Tags:      []*TagResponse{&tag1ov},
				Targets:   []*TargetOverview{},
			},
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("Success7", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		application := &model.ApplicationResponse{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses:  []*model.ApplicationStatus{},
		}
		modelApplications := []*model.ApplicationResponse{application}

		e := echo.New()
		path := fmt.Sprintf("/api/applications?created_by=%s", application.CreatedBy.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplications(c.Request().Context(), model.ApplicationQuery{
				Limit:     100,
				Offset:    0,
				CreatedBy: application.CreatedBy},
			).
			Return(modelApplications, nil)

		err = h.Handlers.GetApplications(c)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*ApplicationResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := []*ApplicationResponse{
			{
				ID:        application.ID,
				Status:    application.Status,
				CreatedAt: date,
				UpdatedAt: date,
				CreatedBy: application.CreatedBy,
				Title:     application.Title,
				Content:   application.Content,
				Tags:      []*TagResponse{},
				Targets:   []*TargetOverview{},
			},
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("InvaildStatus", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		e := echo.New()
		path := "/api/applications?status=invalid-status"
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.GetApplications(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, "invalid status"), err)
	})

	t.Run("FailedToGetApplications", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/applications", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		resErr := errors.New("Failed to get applications.")
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplications(c.Request().Context(), model.ApplicationQuery{
				Limit:  100,
				Offset: 0,
			}).
			Return(nil, resErr)

		err = h.Handlers.GetApplications(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})
}

func TestHandlers_PostApplication(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				CreatedBy: uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		reqApplication := Application{
			CreatedBy: application.CreatedBy,
			Title:     application.Title,
			Content:   application.Content,
		}
		reqBody, err := json.Marshal(reqApplication)
		require.NoError(t, err)
		tags := []*model.Tag{}
		targets := []*model.ApplicationTarget{}

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/applications", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockApplicationRepository.
			EXPECT().
			CreateApplication(
				c.Request().Context(),
				reqApplication.Title, reqApplication.Content,
				tags, targets,
				reqApplication.CreatedBy).
			Return(application, nil)

		require.NoError(t, h.Handlers.PostApplication(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *ApplicationDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelApplicationDetailToApplicationResponse(application)
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithTags", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}
		tags := []*model.Tag{tag}
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
			Tags:      tags,
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				CreatedBy: uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		reqApplication := Application{
			CreatedBy: application.CreatedBy,
			Title:     application.Title,
			Content:   application.Content,
			Tags:      []uuid.UUID{tag.ID},
		}
		reqBody, err := json.Marshal(reqApplication)
		require.NoError(t, err)
		targets := []*model.ApplicationTarget{}

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/applications", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockTagRepository.
			EXPECT().
			GetTag(c.Request().Context(), tag.ID).
			Return(tag, nil)
		h.Repository.MockApplicationRepository.
			EXPECT().
			CreateApplication(
				c.Request().Context(),
				reqApplication.Title, reqApplication.Content,
				tags, targets,
				reqApplication.CreatedBy).
			Return(application, nil)

		require.NoError(t, h.Handlers.PostApplication(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *ApplicationDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelApplicationDetailToApplicationResponse(application)
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithTarget", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		target := &model.ApplicationTarget{
			Target: uuid.New(),
			Amount: random.Numeric(t, 1000000),
		}
		tgd := &model.ApplicationTargetDetail{
			ID:        uuid.New(),
			Target:    target.Target,
			Amount:    target.Amount,
			CreatedAt: date,
		}
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{tgd},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				CreatedBy: uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		tg := &Target{
			Target: target.Target,
			Amount: target.Amount,
		}
		reqApplication := Application{
			CreatedBy: application.CreatedBy,
			Title:     application.Title,
			Content:   application.Content,
			Targets:   []*Target{tg},
		}
		reqBody, err := json.Marshal(reqApplication)
		require.NoError(t, err)
		tags := []*model.Tag{}

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/applications", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockApplicationRepository.
			EXPECT().
			CreateApplication(
				c.Request().Context(),
				reqApplication.Title, reqApplication.Content,
				tags, []*model.ApplicationTarget{target},
				reqApplication.CreatedBy).
			Return(application, nil)

		require.NoError(t, h.Handlers.PostApplication(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *ApplicationDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelApplicationDetailToApplicationResponse(application)
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("UnknownTagID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
		}
		unknownTagID := uuid.New()
		reqApplication := Application{
			CreatedBy: application.CreatedBy,
			Title:     application.Title,
			Content:   application.Content,
			Tags:      []uuid.UUID{unknownTagID},
		}
		reqBody, err := json.Marshal(reqApplication)
		require.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/applications", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown id"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockTagRepository.
			EXPECT().
			GetTag(c.Request().Context(), unknownTagID).
			Return(nil, resErr)

		err = h.Handlers.PostApplication(c)
		require.Error(t, err)
		// FIXME: http.StatusNotFoundだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
	})

	t.Run("UnknownUserID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
		}
		reqApplication := Application{
			CreatedBy: application.CreatedBy,
			Title:     application.Title,
			Content:   application.Content,
		}
		reqBody, err := json.Marshal(reqApplication)
		require.NoError(t, err)
		tags := []*model.Tag{}
		targets := []*model.ApplicationTarget{}

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/applications", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown id"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockApplicationRepository.
			EXPECT().
			CreateApplication(
				c.Request().Context(),
				reqApplication.Title, reqApplication.Content,
				tags, targets,
				reqApplication.CreatedBy).
			Return(nil, resErr)

		err = h.Handlers.PostApplication(c)
		require.Error(t, err)
		// FIXME: http.StatusNotFoundだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
	})
}

func TestHandlers_GetApplication(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/applications/:applicationID")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(c.Request().Context(), application.ID).
			Return(application, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			GetComments(c.Request().Context(), application.ID).
			Return(nil, nil)

		require.NoError(t, h.Handlers.GetApplication(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *ApplicationDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelApplicationDetailToApplicationResponse(application)
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithComments", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)
		date := time.Now()

		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		comment1 := &model.Comment{
			ID:        uuid.New(),
			User:      application.CreatedBy,
			Comment:   random.AlphaNumeric(t, 100),
			CreatedAt: date,
			UpdatedAt: date,
		}
		comment2 := &model.Comment{
			ID:        uuid.New(),
			User:      application.CreatedBy,
			Comment:   random.AlphaNumeric(t, 100),
			CreatedAt: date,
			UpdatedAt: date,
		}
		comments := []*model.Comment{comment1, comment2}

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/applications/:applicationID")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(c.Request().Context(), application.ID).
			Return(application, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			GetComments(c.Request().Context(), application.ID).
			Return(comments, nil)

		require.NoError(t, h.Handlers.GetApplication(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *ApplicationDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelApplicationDetailToApplicationResponse(application)
		exp.Comments = lo.Map(comments, func(c *model.Comment, _ int) *CommentDetail {
			return modelCommentToCommentDetail(c)
		})
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithTarget", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		target := &model.ApplicationTargetDetail{
			ID:        uuid.New(),
			Target:    uuid.New(),
			Amount:    random.Numeric(t, 1000000),
			PaidAt:    time.Time{},
			CreatedAt: date,
		}
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{target},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/applications/:applicationID")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(c.Request().Context(), application.ID).
			Return(application, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			GetComments(c.Request().Context(), application.ID).
			Return(nil, nil)
		require.NoError(t, h.Handlers.GetApplication(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *ApplicationDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelApplicationDetailToApplicationResponse(application)
		exp.Targets = []*TargetOverview{
			modelApplicationTargetDetailToTargetOverview(target),
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		invalidUUID := "invalid-uuid"
		_, resErr := uuid.Parse(invalidUUID)

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s", invalidUUID)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/applications/:applicationID")
		c.SetParamNames("applicationID")
		c.SetParamValues(invalidUUID)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.GetApplication(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s", uuid.Nil)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/applications/:applicationID")
		c.SetParamNames("applicationID")
		c.SetParamValues(uuid.Nil.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		resErr := errors.New("invalid UUID")

		err = h.Handlers.GetApplication(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("UnknownID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		unknownID := uuid.New()

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s", unknownID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/applications/:applicationID")
		c.SetParamNames("applicationID")
		c.SetParamValues(unknownID.String())

		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown id"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(c.Request().Context(), unknownID).
			Return(nil, resErr)

		err = h.Handlers.GetApplication(c)
		require.Error(t, err)
		// FIXME: http.StatusNotFoundだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
	})
}

func TestHandlers_PutApplication(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		reqApplication := PutApplication{
			Title:   random.AlphaNumeric(t, 30),
			Content: random.AlphaNumeric(t, 50),
			Tags:    []uuid.UUID{},
			Targets: []*Target{},
		}
		reqBody, err := json.Marshal(reqApplication)
		require.NoError(t, err)
		tags := []*model.Tag{}
		updateApplication := &model.ApplicationDetail{
			ID:        application.ID,
			Status:    application.Status,
			CreatedBy: application.CreatedBy,
			CreatedAt: application.CreatedAt,
			UpdatedAt: time.Now(),
			Title:     reqApplication.Title,
			Content:   reqApplication.Content,
			Tags:      tags,
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses:  application.Statuses,
			Comments:  application.Comments,
			Files:     application.Files,
		}

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s", application.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/applications/:applicationID")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())
		c.Set(loginUserKey, user)

		targets := lo.Map(
			updateApplication.Targets,
			func(t *model.ApplicationTargetDetail, _ int) *model.ApplicationTarget {
				return &model.ApplicationTarget{
					Target: t.Target,
					Amount: t.Amount,
				}
			},
		)
		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(c.Request().Context(), application.ID).
			Return(application, nil)
		h.Repository.MockApplicationRepository.
			EXPECT().
			UpdateApplication(
				c.Request().Context(),
				application.ID,
				reqApplication.Title, reqApplication.Content,
				tags, targets).
			Return(updateApplication, nil)

		require.NoError(t, h.Handlers.PutApplication(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *ApplicationDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelApplicationDetailToApplicationResponse(updateApplication)
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithTag", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		tag1 := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}
		tag2 := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}
		tags := []*model.Tag{tag1, tag2}
		reqApplication := PutApplication{
			Title:   random.AlphaNumeric(t, 30),
			Content: random.AlphaNumeric(t, 50),
			Tags:    []uuid.UUID{tag1.ID, tag2.ID},
			Targets: []*Target{},
		}
		reqBody, err := json.Marshal(reqApplication)
		require.NoError(t, err)
		updateApplication := &model.ApplicationDetail{
			ID:        application.ID,
			Status:    application.Status,
			CreatedBy: application.CreatedBy,
			Title:     application.Title,
			Content:   reqApplication.Content,
			CreatedAt: application.CreatedAt,
			UpdatedAt: time.Now(),
			Tags:      tags,
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses:  application.Statuses,
			Comments:  application.Comments,
			Files:     application.Files,
		}

		targets := lo.Map(
			updateApplication.Targets,
			func(t *model.ApplicationTargetDetail, _ int) *model.ApplicationTarget {
				return &model.ApplicationTarget{
					Target: t.Target,
					Amount: t.Amount,
				}
			},
		)
		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/applications/:applicationID")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(c.Request().Context(), application.ID).
			Return(application, nil)
		h.Repository.MockTagRepository.
			EXPECT().
			GetTag(c.Request().Context(), tag1.ID).
			Return(tag1, nil)
		h.Repository.MockTagRepository.
			EXPECT().
			GetTag(c.Request().Context(), tag2.ID).
			Return(tag2, nil)
		h.Repository.MockApplicationRepository.
			EXPECT().
			UpdateApplication(
				c.Request().Context(),
				application.ID,
				reqApplication.Title, reqApplication.Content,
				tags, targets).
			Return(updateApplication, nil)

		require.NoError(t, h.Handlers.PutApplication(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *ApplicationDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelApplicationDetailToApplicationResponse(updateApplication)
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithTarget", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		target1 := &model.ApplicationTargetDetail{
			ID:        uuid.New(),
			Target:    uuid.New(),
			Amount:    random.Numeric(t, 100000),
			PaidAt:    time.Time{},
			CreatedAt: date,
		}
		target2 := &model.ApplicationTargetDetail{
			ID:        uuid.New(),
			Target:    uuid.New(),
			Amount:    random.Numeric(t, 100000),
			PaidAt:    time.Time{},
			CreatedAt: date,
		}
		targetDetails := []*model.ApplicationTargetDetail{target1, target2}
		reqApplication := PutApplication{
			Title:   random.AlphaNumeric(t, 30),
			Content: random.AlphaNumeric(t, 50),
			Tags:    []uuid.UUID{},
			Targets: []*Target{
				{Target: target1.Target, Amount: target1.Amount},
				{Target: target2.Target, Amount: target2.Amount},
			},
		}
		reqBody, err := json.Marshal(reqApplication)
		require.NoError(t, err)
		tags := []*model.Tag{}
		updateApplication := &model.ApplicationDetail{
			ID:        application.ID,
			Status:    application.Status,
			CreatedBy: application.CreatedBy,
			Title:     application.Title,
			Content:   reqApplication.Content,
			CreatedAt: application.CreatedAt,
			UpdatedAt: time.Now(),
			Tags:      tags,
			Targets:   targetDetails,
			Statuses:  application.Statuses,
			Comments:  application.Comments,
			Files:     application.Files,
		}

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/applications/:applicationID")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())
		c.Set(loginUserKey, user)

		targets := lo.Map(
			updateApplication.Targets,
			func(t *model.ApplicationTargetDetail, _ int) *model.ApplicationTarget {
				return &model.ApplicationTarget{
					Target: t.Target,
					Amount: t.Amount,
				}
			},
		)
		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(c.Request().Context(), application.ID).
			Return(application, nil)
		h.Repository.MockApplicationRepository.
			EXPECT().
			UpdateApplication(
				c.Request().Context(),
				application.ID,
				reqApplication.Title, reqApplication.Content,
				tags, targets).
			Return(updateApplication, nil)

		require.NoError(t, h.Handlers.PutApplication(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *ApplicationDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelApplicationDetailToApplicationResponse(updateApplication)
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithComment", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		reqApplication := PutApplication{
			Title:   random.AlphaNumeric(t, 30),
			Content: random.AlphaNumeric(t, 50),
			Tags:    []uuid.UUID{},
			Targets: []*Target{},
		}
		reqBody, err := json.Marshal(reqApplication)
		require.NoError(t, err)
		tags := []*model.Tag{}
		comment1 := &model.Comment{
			ID:        uuid.New(),
			User:      application.CreatedBy,
			Comment:   random.AlphaNumeric(t, 100),
			CreatedAt: date,
			UpdatedAt: date,
		}
		comment2 := &model.Comment{
			ID:        uuid.New(),
			User:      application.CreatedBy,
			Comment:   random.AlphaNumeric(t, 100),
			CreatedAt: date,
			UpdatedAt: date,
		}
		comments := []*model.Comment{comment1, comment2}
		updateApplication := &model.ApplicationDetail{
			ID:        application.ID,
			Status:    application.Status,
			CreatedBy: application.CreatedBy,
			Title:     reqApplication.Title,
			Content:   reqApplication.Content,
			CreatedAt: application.CreatedAt,
			UpdatedAt: time.Now(),
			Tags:      tags,
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses:  application.Statuses,
			Comments:  comments,
			Files:     application.Files,
		}

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s", application.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/applications/:applicationID")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())
		c.Set(loginUserKey, user)

		targets := lo.Map(
			updateApplication.Targets,
			func(t *model.ApplicationTargetDetail, _ int) *model.ApplicationTarget {
				return &model.ApplicationTarget{
					Target: t.Target,
					Amount: t.Amount,
				}
			},
		)
		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(c.Request().Context(), application.ID).
			Return(application, nil)
		h.Repository.MockApplicationRepository.
			EXPECT().
			UpdateApplication(
				c.Request().Context(),
				application.ID,
				reqApplication.Title, reqApplication.Content,
				tags, targets).
			Return(updateApplication, nil)

		require.NoError(t, h.Handlers.PutApplication(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *ApplicationDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelApplicationDetailToApplicationResponse(updateApplication)
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		invalidUUID := "invalid-uuid"
		_, resErr := uuid.Parse(invalidUUID)

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s", invalidUUID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/applications/:applicationID")
		c.SetParamNames("applicationID")
		c.SetParamValues(invalidUUID)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.PutApplication(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s", uuid.Nil)
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/applications/:applicationID")
		c.SetParamNames("applicationID")
		c.SetParamValues(uuid.Nil.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		resErr := errors.New("invalid UUID")

		err = h.Handlers.PutApplication(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("UnknownID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		unknownID := uuid.New()
		reqApplication := PutApplication{
			Title:   random.AlphaNumeric(t, 30),
			Content: random.AlphaNumeric(t, 50),
			Tags:    []uuid.UUID{},
			Targets: []*Target{},
		}
		reqBody, err := json.Marshal(reqApplication)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s", unknownID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/applications/:applicationID")
		c.SetParamNames("applicationID")
		c.SetParamValues(unknownID.String())
		c.Set(loginUserKey, user)

		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown id"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(c.Request().Context(), unknownID).
			Return(nil, resErr)

		err = h.Handlers.PutApplication(c)
		require.Error(t, err)
		// FIXME: http.StatusNotFoundだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
	})

	t.Run("UnknownTagID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}
		reqApplication := PutApplication{
			Title:   random.AlphaNumeric(t, 30),
			Content: random.AlphaNumeric(t, 50),
			Tags:    []uuid.UUID{tag.ID},
			Targets: []*Target{},
		}
		reqBody, err := json.Marshal(reqApplication)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/applications/:applicationID")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())
		c.Set(loginUserKey, user)

		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown id"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(c.Request().Context(), application.ID).
			Return(application, nil)
		h.Repository.MockTagRepository.
			EXPECT().
			GetTag(c.Request().Context(), tag.ID).
			Return(nil, resErr)

		err = h.Handlers.PutApplication(c)
		require.Error(t, err)
		// FIXME: http.StatusNotFoundだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
	})
}

func TestHandlers_PutStatus(t *testing.T) {
	t.Parallel()

	t.Run("SuccessByCreator", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.FixRequired,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.FixRequired,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status:  model.Submitted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)
		comment := &model.Comment{
			ID:        uuid.New(),
			User:      user.ID,
			Comment:   reqStatus.Comment,
			CreatedAt: date,
			UpdatedAt: date,
		}
		status := &model.ApplicationStatus{
			ID:        uuid.New(),
			CreatedBy: user.ID,
			Status:    reqStatus.Status,
			CreatedAt: date,
		}

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s/status", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/applications/:applicationID/status")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(ctx, application.ID).
			Return(application, nil)
		h.Repository.MockApplicationStatusRepository.
			EXPECT().
			CreateStatus(ctx, application.ID, user.ID, reqStatus.Status).
			Return(status, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			CreateComment(ctx, reqStatus.Comment, application.ID, user.ID).
			Return(comment, nil)

		require.NoError(t, h.Handlers.PutStatus(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *StatusResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &StatusResponse{
			CreatedBy: user.ID,
			Status:    status.Status,
			Comment: CommentDetail{
				ID:        comment.ID,
				User:      comment.User,
				Comment:   comment.Comment,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
			},
			CreatedAt: status.CreatedAt,
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessByAccountManagerFromSubmittedToFixRequired", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status:  model.FixRequired,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)
		comment := &model.Comment{
			ID:        uuid.New(),
			User:      user.ID,
			Comment:   reqStatus.Comment,
			CreatedAt: date,
			UpdatedAt: date,
		}
		status := &model.ApplicationStatus{
			ID:        uuid.New(),
			CreatedBy: user.ID,
			Status:    reqStatus.Status,
			CreatedAt: date,
		}

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s/status", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/applications/:applicationID/status")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(ctx, application.ID).
			Return(application, nil)
		h.Repository.MockApplicationStatusRepository.
			EXPECT().
			CreateStatus(ctx, application.ID, user.ID, reqStatus.Status).
			Return(status, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			CreateComment(ctx, reqStatus.Comment, application.ID, user.ID).
			Return(comment, nil)

		require.NoError(t, h.Handlers.PutStatus(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *StatusResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &StatusResponse{
			CreatedBy: user.ID,
			Status:    status.Status,
			Comment: CommentDetail{
				ID:        comment.ID,
				User:      comment.User,
				Comment:   comment.Comment,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
			},
			CreatedAt: status.CreatedAt,
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessByAccountManagerFromSubmittedToAccepted", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status:  model.Accepted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)
		comment := &model.Comment{
			ID:        uuid.New(),
			User:      user.ID,
			Comment:   reqStatus.Comment,
			CreatedAt: date,
			UpdatedAt: date,
		}
		status := &model.ApplicationStatus{
			ID:        uuid.New(),
			CreatedBy: user.ID,
			Status:    reqStatus.Status,
			CreatedAt: date,
		}

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s/status", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/applications/:applicationID/status")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(ctx, application.ID).
			Return(application, nil)
		h.Repository.MockApplicationStatusRepository.
			EXPECT().
			CreateStatus(ctx, application.ID, user.ID, reqStatus.Status).
			Return(status, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			CreateComment(ctx, reqStatus.Comment, application.ID, user.ID).
			Return(comment, nil)

		require.NoError(t, h.Handlers.PutStatus(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *StatusResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &StatusResponse{
			CreatedBy: user.ID,
			Status:    status.Status,
			Comment: CommentDetail{
				ID:        comment.ID,
				User:      comment.User,
				Comment:   comment.Comment,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
			},
			CreatedAt: status.CreatedAt,
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessByAccountManagerFromSubmittedToFixRequired", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status:  model.FixRequired,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		comment := &model.Comment{
			ID:        uuid.New(),
			User:      user.ID,
			Comment:   reqStatus.Comment,
			CreatedAt: date,
			UpdatedAt: date,
		}
		status := &model.ApplicationStatus{
			ID:        uuid.New(),
			CreatedBy: user.ID,
			Status:    reqStatus.Status,
			CreatedAt: date,
		}

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s/status", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/applications/:applicationID/status")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(ctx, application.ID).
			Return(application, nil)
		h.Repository.MockApplicationStatusRepository.
			EXPECT().
			CreateStatus(ctx, application.ID, user.ID, reqStatus.Status).
			Return(status, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			CreateComment(ctx, reqStatus.Comment, application.ID, user.ID).
			Return(comment, nil)

		require.NoError(t, h.Handlers.PutStatus(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *StatusResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &StatusResponse{
			CreatedBy: user.ID,
			Status:    status.Status,
			Comment: CommentDetail{
				ID:        comment.ID,
				User:      comment.User,
				Comment:   comment.Comment,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
			},
			CreatedAt: status.CreatedAt,
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessByAccountManagerFromFixRequiredToSubmitted", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.FixRequired,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.FixRequired,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status:  model.Submitted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)
		comment := &model.Comment{
			ID:        uuid.New(),
			User:      user.ID,
			Comment:   reqStatus.Comment,
			CreatedAt: date,
			UpdatedAt: date,
		}
		status := &model.ApplicationStatus{
			ID:        uuid.New(),
			CreatedBy: user.ID,
			Status:    reqStatus.Status,
			CreatedAt: date,
		}

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s/status", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/applications/:applicationID/status")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(ctx, application.ID).
			Return(application, nil)
		h.Repository.MockApplicationStatusRepository.
			EXPECT().
			CreateStatus(ctx, application.ID, user.ID, reqStatus.Status).
			Return(status, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			CreateComment(ctx, reqStatus.Comment, application.ID, user.ID).
			Return(comment, nil)

		require.NoError(t, h.Handlers.PutStatus(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *StatusResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &StatusResponse{
			CreatedBy: user.ID,
			Status:    status.Status,
			Comment: CommentDetail{
				ID:        comment.ID,
				User:      comment.User,
				Comment:   comment.Comment,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
			},
			CreatedAt: status.CreatedAt,
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessByAccountManagerFromAcceptedToSubmitted", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Accepted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.Accepted,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		target := &model.ApplicationTargetDetail{
			ID:        uuid.New(),
			Target:    uuid.New(),
			PaidAt:    time.Time{},
			CreatedAt: date,
		}
		targets := []*model.ApplicationTargetDetail{target}
		reqStatus := PutStatus{
			Status:  model.Submitted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)
		comment := &model.Comment{
			ID:        uuid.New(),
			User:      user.ID,
			Comment:   reqStatus.Comment,
			CreatedAt: date,
			UpdatedAt: date,
		}
		status := &model.ApplicationStatus{
			ID:        uuid.New(),
			CreatedBy: user.ID,
			Status:    reqStatus.Status,
			CreatedAt: date,
		}

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s/status", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/applications/:applicationID/status")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(ctx, application.ID).
			Return(application, nil)
		h.Repository.MockApplicationTargetRepository.
			EXPECT().
			GetApplicationTargets(ctx, application.ID).
			Return(targets, nil)
		h.Repository.MockApplicationStatusRepository.
			EXPECT().
			CreateStatus(ctx, application.ID, user.ID, reqStatus.Status).
			Return(status, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			CreateComment(ctx, reqStatus.Comment, application.ID, user.ID).
			Return(comment, nil)

		require.NoError(t, h.Handlers.PutStatus(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *StatusResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &StatusResponse{
			CreatedBy: user.ID,
			Status:    status.Status,
			Comment: CommentDetail{
				ID:        comment.ID,
				User:      comment.User,
				Comment:   comment.Comment,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
			},
			CreatedAt: status.CreatedAt,
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("InvalidStatus", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.FixRequired,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.FixRequired,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		invalidStatus := random.AlphaNumeric(t, 20)
		reqBody, err := json.Marshal(&struct {
			Status  string `json:"status"`
			Comment string `json:"comment"`
		}{
			Status:  invalidStatus,
			Comment: random.AlphaNumeric(t, 20),
		})
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s/status", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/applications/:applicationID/status")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		resErr := echo.NewHTTPError(http.StatusBadRequest)
		resErrMessage := echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("invalid Status %s", invalidStatus))
		resErrMessage.Internal = fmt.Errorf("invalid Status %s", invalidStatus)
		resErr.Message = resErrMessage

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい
		require.Equal(t, resErr, err)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		invalidUUID := "invalid-uuid"
		_, resErr := uuid.Parse(invalidUUID)
		reqStatus := PutStatus{
			Status:  model.Submitted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s/status", invalidUUID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/applications/:applicationID/status")
		c.SetParamNames("applicationID")
		c.SetParamValues(invalidUUID)
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("NillUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		reqStatus := PutStatus{
			Status:  model.Submitted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s/status", uuid.Nil)
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/applications/:applicationID/status")
		c.SetParamNames("applicationID")
		c.SetParamValues(uuid.Nil.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		_, resErr := uuid.Parse(c.Param("applicationID"))

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("SameStatusError", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Status(random.Numeric(t, 5) + 1),
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.Status(random.Numeric(t, 5) + 1),
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status:  application.Status,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s/status", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/applications/:applicationID/status")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(ctx, application.ID).
			Return(application, nil)

		resErr := errors.New("invalid application: same status")

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("CommentRequiredErrorFromSubmittedToFixRequired", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)
		date := time.Now()

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}

		reqStatus := PutStatus{
			Status: model.FixRequired,
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s/status", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/applications/:applicationID/status")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(ctx, application.ID).
			Return(application, nil)

		resErr := fmt.Errorf(
			"unable to change %v to %v without comment",
			application.Status.String(),
			reqStatus.Status.String())

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("CommentRequiredErrorFromSubmittedToRejected", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status: model.Rejected,
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s/status", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/applications/:applicationID/status")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(ctx, application.ID).
			Return(application, nil)

		resErr := fmt.Errorf(
			"unable to change %v to %v without comment",
			application.Status.String(),
			reqStatus.Status.String())

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("CommentRequiredErrorFromAcceptedToSubmitted", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Accepted,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.Accepted,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status: model.Submitted,
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s/status", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/applications/:applicationID/status")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(ctx, application.ID).
			Return(application, nil)

		resErr := fmt.Errorf(
			"unable to change %v to %v without comment",
			application.Status.String(),
			reqStatus.Status.String())

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("AccountManagerNoPrivilege", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.FixRequired,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.FixRequired,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status:  model.Accepted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s/status", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/applications/:applicationID/status")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(ctx, application.ID).
			Return(application, nil)

		resErr := fmt.Errorf(
			"accountManager unable to change %v to %v",
			application.Status.String(),
			reqStatus.Status.String())

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusForbiddenだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("AlreadyPaid", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Accepted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.Accepted,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		target := &model.ApplicationTargetDetail{
			ID:        uuid.New(),
			Target:    uuid.New(),
			PaidAt:    date,
			CreatedAt: date,
		}
		targets := []*model.ApplicationTargetDetail{target}
		reqStatus := PutStatus{
			Status:  model.Submitted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s/status", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/applications/:applicationID/status")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		ctx = c.Request().Context()
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(ctx, application.ID).
			Return(application, nil)
		h.Repository.MockApplicationTargetRepository.
			EXPECT().
			GetApplicationTargets(ctx, application.ID).
			Return(targets, nil)

		resErr := errors.New("someone already paid")

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("CreatorNoPrivilege", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status:  model.Accepted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s/status", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/applications/:applicationID/status")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(ctx, application.ID).
			Return(application, nil)

		resErr := fmt.Errorf(
			"creator unable to change %v to %v",
			application.Status.String(), reqStatus.Status.String())

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusForbiddenだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("NoPrivilege", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		application := &model.ApplicationDetail{
			ID:        uuid.New(),
			Status:    model.FixRequired,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.ApplicationTargetDetail{},
			Statuses: []*model.ApplicationStatus{{
				ID:        uuid.New(),
				Status:    model.FixRequired,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Comments: []*model.Comment{},
			Files:    []uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status:  model.Submitted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/applications/%s/status", application.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/applications/:applicationID/status")
		c.SetParamNames("applicationID")
		c.SetParamValues(application.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockApplicationRepository.
			EXPECT().
			GetApplication(ctx, application.ID).
			Return(application, nil)

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusForbiddenだけ判定したい
		require.Equal(
			t,
			echo.NewHTTPError(http.StatusForbidden, "you are not application creator"),
			err)
	})
}

// TODO: TestHandlers_PostComment
