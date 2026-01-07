package model

import (
	"github.com/traPtitech/Jomon/internal/ent"
	"github.com/traPtitech/Jomon/internal/service"
)

type entErrorConverter struct {
	msgBadInput string
	msgNotFound string
}

var defaultEntErrorConverter = &entErrorConverter{
	msgBadInput: "invalid input",
	msgNotFound: "specified data not found",
}

func (ec *entErrorConverter) convert(err error) error {
	switch {
	case err == nil:
		return nil
	case ent.IsConstraintError(err):
		return service.NewBadInputError(ec.msgBadInput).WithInternal(err)
	case ent.IsNotFound(err):
		return service.NewNotFoundError(ec.msgNotFound).WithInternal(err)
	case ent.IsNotLoaded(err):
		return service.NewUnexpectedError(err)
	case ent.IsNotSingular(err):
		return service.NewUnexpectedError(err)
	case ent.IsValidationError(err):
		return service.NewBadInputError(ec.msgBadInput).WithInternal(err)
	default:
		return err
	}
}
