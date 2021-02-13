package model

import (
	"time"

	"github.com/gofrs/uuid"
)

// Comment struct of Comment
type Comment struct {
	ID        int        `gorm:"type:int(11) AUTO_INCREMENT;primary_key" json:"comment_id"`
	RequestID uuid.UUID  `gorm:"type:char(36);not null;index" json:"-"`
	CreatedBy TrapUser   `gorm:"embedded;embedded_prefix:created_by_" json:"trap_user"`
	Comment   string     `gorm:"type:text;not null" json:"comment"`
	CreatedAt time.Time  `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

// GiveIsUserAdmin check whether comment is admin or not
func (com *Comment) GiveIsUserAdmin(admins []string) {
	if com == nil {
		return
	}

	com.CreatedBy.GiveIsUserAdmin(admins)
}

// CommentRepository Repo of Comment
type CommentRepository interface {
	GetComment(requestID uuid.UUID, commentID int) (Comment, error)
	CreateComment(requestID uuid.UUID, commentText string, userID string) (Comment, error)
	PutComment(requestID uuid.UUID, commentID int, commentText string) (Comment, error)
	DeleteComment(requestID uuid.UUID, commentID int) error
}

type commentRepository struct{}

// NewCommentRepository Make CommentRepository
func NewCommentRepository() CommentRepository {
	return &commentRepository{}
}

func (*commentRepository) GetComment(requestID uuid.UUID, commentID int) (Comment, error) {
	comment := Comment{
		ID:        commentID,
		RequestID: requestID,
	}

	if err := db.First(&comment).Error; err != nil {
		return Comment{}, err
	}
	return comment, nil
}

func (*commentRepository) CreateComment(requestID uuid.UUID, commentText string, userID string) (Comment, error) {
	comment := Comment{
		RequestID: requestID,
		CreatedBy: TrapUser{TrapID: userID},
		Comment:   commentText,
	}

	if err := db.Create(&comment).Error; err != nil {
		return Comment{}, err
	}

	return comment, nil
}

func (*commentRepository) PutComment(requestID uuid.UUID, commentID int, commentText string) (Comment, error) {
	comment := Comment{
		ID:        commentID,
		RequestID: requestID,
	}

	if err := db.Model(&comment).Update("Comment", commentText).Error; err != nil {
		return Comment{}, err
	}

	return comment, nil
}

func (*commentRepository) DeleteComment(requestID uuid.UUID, commentID int) error {
	comment := Comment{
		ID:        commentID,
		RequestID: requestID,
	}

	if err := db.Delete(&comment).Error; err != nil {
		return err
	}
	return nil
}
