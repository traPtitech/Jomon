//go:generate go tool mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/internal/logging"
	"github.com/traPtitech/Jomon/internal/nulltime"
	"go.uber.org/zap"
)

type ApplicationRepository interface {
	GetApplications(ctx context.Context, query ApplicationQuery) ([]*ApplicationResponse, error)
	CreateApplication(
		ctx context.Context, title string, content string,
		tags []*Tag, targets []*ApplicationTarget, userID uuid.UUID,
	) (*ApplicationDetail, error)
	GetApplication(ctx context.Context, applicationID uuid.UUID) (*ApplicationDetail, error)
	UpdateApplication(
		ctx context.Context, applicationID uuid.UUID, title string, content string,
		tags []*Tag, targets []*ApplicationTarget,
	) (*ApplicationDetail, error)
}

type Application struct {
	ID        uuid.UUID
	Amount    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ApplicationResponse struct {
	ID        uuid.UUID
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uuid.UUID
	Title     string
	Content   string
	Tags      []*Tag
	Targets   []*ApplicationTargetDetail
	Statuses  []*ApplicationStatus
}

type ApplicationDetail struct {
	ID        uuid.UUID
	Status    Status
	Title     string
	Content   string
	Comments  []*Comment
	Files     []uuid.UUID
	Statuses  []*ApplicationStatus
	Tags      []*Tag
	Targets   []*ApplicationTargetDetail
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uuid.UUID
}

type ApplicationQuery struct {
	Sort      *string
	Target    uuid.UUID
	Status    *string
	Since     nulltime.NullTime
	Until     nulltime.NullTime
	Limit     int
	Offset    int
	Tag       *string
	CreatedBy uuid.UUID
}

func (s *Service) GetApplications(
	ctx context.Context, query ApplicationQuery,
) ([]*ApplicationResponse, error) {
	logger := logging.GetLogger(ctx)
	applications, err := s.repository.GetApplications(ctx, query)
	if err != nil {
		logger.Error("failed to get applications from repository", zap.Error(err))
		return nil, err
	}
	return applications, nil
}

type CreateApplicationInputs struct {
	Title   string
	Content string
	Tags    []*Tag
	Targets []*ApplicationTarget
}

func (s *Service) CreateApplication(
	ctx context.Context, user *User, inputs CreateApplicationInputs,
) (*ApplicationDetail, error) {
	logger := logging.GetLogger(ctx)
	application, err := s.repository.CreateApplication(
		ctx, inputs.Title, inputs.Content, inputs.Tags, inputs.Targets, user.ID,
	)
	if err != nil {
		logger.Error("failed to create application in repository", zap.Error(err))
		return nil, err
	}
	return application, nil
}

func (s *Service) GetApplication(
	ctx context.Context, applicationID uuid.UUID,
) (*ApplicationDetail, error) {
	logger := logging.GetLogger(ctx)
	application, err := s.repository.GetApplication(ctx, applicationID)
	if err != nil {
		logger.Error("failed to get application from repository", zap.Error(err))
		return nil, err
	}
	comments, err := s.repository.GetComments(ctx, applicationID)
	if err != nil {
		logger.Error("failed to get comments from repository", zap.Error(err))
		return nil, err
	}
	application.Comments = comments
	return application, nil
}

type UpdateApplicationInputs struct {
	Title   string
	Content string
	TagIDs  []uuid.UUID
	Targets []*ApplicationTarget
}

func (s *Service) UpdateApplication(
	ctx context.Context,
	applicationID uuid.UUID, user *User, inputs UpdateApplicationInputs,
) (*ApplicationDetail, error) {
	logger := logging.GetLogger(ctx).With(
		zap.String("applicationID", applicationID.String()), zap.String("userID", user.ID.String()))
	currentApplication, err := s.repository.GetApplication(ctx, applicationID)
	if err != nil {
		logger.Error("failed to get application from repository", zap.Error(err))
		return nil, err
	}
	if !s.isApplicationCreator(currentApplication, user.ID) {
		logger.Info("user is not application creator")
		return nil, NewForbiddenError("user is not application creator")
	}
	tags := []*Tag{}
	for _, tagID := range inputs.TagIDs {
		// FIXME: N+1
		tag, err := s.repository.GetTag(ctx, tagID)
		if err != nil {
			logger.Error("failed to get tag from repository", zap.Error(err))
			return nil, err
		}
		tags = append(tags, tag)
	}
	application, err := s.repository.UpdateApplication(
		ctx, applicationID, inputs.Title, inputs.Content, tags, inputs.Targets)
	if err != nil {
		logger.Error("failed to update application in repository", zap.Error(err))
		return nil, err
	}
	return application, nil
}

func (s *Service) isApplicationCreator(
	application *ApplicationDetail, userID uuid.UUID,
) bool {
	return application.CreatedBy == userID
}

func (s *Service) CreateCommentToApplication(
	ctx context.Context, applicationID uuid.UUID, user *User, content string,
) (*Comment, error) {
	logger := logging.GetLogger(ctx).With(
		zap.String("applicationID", applicationID.String()), zap.String("userID", user.ID.String()))
	comment, err := s.repository.CreateComment(ctx, content, applicationID, user.ID)
	if err != nil {
		logger.Error("failed to create comment in repository", zap.Error(err))
		return nil, err
	}
	return comment, nil
}

type UpdateStatusInputs struct {
	Status  Status
	Comment string
}

func (s *Service) UpdateApplicationStatus(
	ctx context.Context, applicationID uuid.UUID, user *User, inputs UpdateStatusInputs,
) (*ApplicationStatus, *Comment, error) {
	logger := logging.GetLogger(ctx).With(
		zap.String("applicationID", applicationID.String()), zap.String("userID", user.ID.String()))
	application, err := s.repository.GetApplication(ctx, applicationID)
	if err != nil {
		logger.Error("failed to get application from repository", zap.Error(err))
		return nil, nil, err
	}
	if !s.isApplicationCreator(application, user.ID) && !s.isAccountManager(user) {
		logger.Info("user is not application creator nor account manager")
		return nil, nil, NewForbiddenError("user is not application creator nor account manager")
	}
	if inputs.Status == application.Status {
		return nil, nil, NewBadInputError("status is not changed")
	}
	if inputs.Comment == "" && !s.isAbleNoCommentUpdateStatus(application.Status, inputs.Status) {
		message := fmt.Sprintf(
			"unable to change %v to %v without comment",
			application.Status.String(), inputs.Status.String())
		return nil, nil, NewBadInputError(message)
	}
	if s.isAccountManager(user) &&
		!s.isAbleAccountManagerUpdateStatus(application.Status, inputs.Status) {
		message := fmt.Sprintf(
			"account manager is unable to change %v to %v",
			application.Status.String(), inputs.Status.String())
		return nil, nil, NewBadInputError(message)
	}
	if s.isAccountManagerRevertingAcceptedApplication(user, application.Status, inputs.Status) {
		targets, err := s.repository.GetApplicationTargets(ctx, applicationID)
		if err != nil {
			logger.Error("failed to get application targets from repository", zap.Error(err))
			return nil, nil, err
		}
		paid := lo.Reduce(targets, func(p bool, target *ApplicationTargetDetail, _ int) bool {
			return p || target.PaidAt.Valid
		}, false)
		if paid {
			return nil, nil, NewBadInputError("someone already paid")
		}
	}
	if s.isApplicationCreator(application, user.ID) &&
		!s.isAbleCreatorChangeStatus(application.Status, inputs.Status) {
		message := fmt.Sprintf(
			"application creator is unable to change %v to %v",
			application.Status.String(), inputs.Status.String())
		return nil, nil, NewForbiddenError(message)
	}
	newStatus, err := s.repository.CreateStatus(ctx, applicationID, user.ID, inputs.Status)
	if err != nil {
		logger.Error("failed to create status in repository", zap.Error(err))
		return nil, nil, err
	}
	var comment *Comment
	if inputs.Comment != "" {
		comment, err = s.repository.CreateComment(ctx, inputs.Comment, applicationID, user.ID)
		if err != nil {
			logger.Error("failed to create comment in repository", zap.Error(err))
			return nil, nil, err
		}
	}
	return newStatus, comment, nil
}

func (s *Service) isAccountManager(user *User) bool {
	return user.AccountManager
}

func (s *Service) isAbleNoCommentUpdateStatus(currentStatus, newStatus Status) bool {
	switch currentStatus {
	case Submitted:
		return newStatus != FixRequired && newStatus != Rejected
	case Accepted:
		return newStatus != Submitted
	case FixRequired, Completed, Rejected:
		return true
	}
	// the switch above performs exhaustive check
	panic("unreachable")
}

func (s *Service) isAbleAccountManagerUpdateStatus(currentStatus, newStatus Status) bool {
	return newStatus == Rejected && currentStatus == Submitted ||
		newStatus == Submitted && currentStatus == FixRequired ||
		newStatus == Accepted && currentStatus == Submitted ||
		newStatus == Submitted && currentStatus == Accepted ||
		newStatus == FixRequired && currentStatus == Submitted
}

func (s *Service) isAbleCreatorChangeStatus(currentStatus, newStatus Status) bool {
	return currentStatus == FixRequired && newStatus == Submitted
}

func (s *Service) isAccountManagerRevertingAcceptedApplication(
	user *User, currentStatus, newStatus Status,
) bool {
	return user.AccountManager && currentStatus == Accepted && newStatus == Submitted
}
