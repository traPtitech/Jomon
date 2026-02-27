package service

type Repository interface {
	AccountManagerRepository
	CommentRepository
	FileRepository
	ApplicationFileRepository
	ApplicationStatusRepository
	ApplicationTagRepository
	ApplicationTargetRepository
	ApplicationRepository
	TagRepository
	UserRepository
}
