package service

type Service struct {
	repository Repository
	storage    Storage
	oidcClient OIDCClient
}

func New(repository Repository, storage Storage, oidcClient OIDCClient) *Service {
	return &Service{
		repository: repository,
		storage:    storage,
		oidcClient: oidcClient,
	}
}
