package service

import (
	"github.com/traPtitech/Jomon/model"
)

type Service interface {
}
type Services struct {
	Repository model.Repository
}

func NewServices(repo model.Repository) (Services, error) {
	return Services{repo}, nil
}
