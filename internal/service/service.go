package service

type Service struct {
	repository Repository
	storage    Storage
}

func New(repository Repository, storage Storage) *Service {
	return &Service{
		repository: repository,
		storage:    storage,
	}
}
