package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/internal/ent/usersubject"
	"github.com/traPtitech/Jomon/internal/service"
)

var userSubjectErrorConverter = &entErrorConverter{
	msgBadInput: "failed to process user subject due to invalid input",
	msgNotFound: "user subject not found",
}

func (repo *EntRepository) GetUserBySubject(
	ctx context.Context, subject string,
) (*service.User, error) {
	us, err := repo.client.UserSubject.
		Query().
		Where(usersubject.IDEQ(subject)).
		WithUser().
		Only(ctx)
	if err != nil {
		return nil, userSubjectErrorConverter.convert(err)
	}
	return convertEntUserToModelUser(us.Edges.User), nil
}

func (repo *EntRepository) RegisterUserSubject(
	ctx context.Context, subject string, userID uuid.UUID,
) error {
	err := repo.client.UserSubject.
		Create().
		SetID(subject).
		SetUserID(userID).
		Exec(ctx)
	return userSubjectErrorConverter.convert(err)
}
