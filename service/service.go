package service

import "github.com/traPtitech/Jomon/ent"

type Service interface {
}
type Services struct {
	EntCli *ent.Client
}

func NewServices(client *ent.Client) (Services, error) {
	return Services{client}, nil
}
