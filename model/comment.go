package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Comment struct {
	ID            int        `gorm:"type:int(11) AUTO_INCREMENT;primary_key" json:"comment_id"`
	ApplicationID uuid.UUID  `gorm:"type:char(36);not null" json:"-"`
	UserTrapID    User       `gorm:"embedded;embedded_prefix:user_" json:"user"`
	Comment       string     `gorm:"type:text;not null" json:"comment"`
	CreatedAt     time.Time  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     *time.Time `json:"-"`
}

func (com *Comment) GiveIsUserAdmin(admins []string) {
	if com == nil {
		return
	}

	com.UserTrapID.GiveIsUserAdmin(admins)
}

func GetComment(applicationId uuid.UUID, commentId int) (Comment, error) {
	comment := Comment{
		ID:            commentId,
		ApplicationID: applicationId,
	}

	if err := db.First(&comment).Error; err != nil {
		return Comment{}, err
	}
	return comment, nil
}

func CreateComment(applicationId uuid.UUID, commentBody string, userId string) (Comment, error) {
	comment := Comment{
		ApplicationID: applicationId,
		UserTrapID:    User{TrapId: userId},
		Comment:       commentBody,
	}

	if err := db.Create(&comment).Error; err != nil {
		return Comment{}, err
	}

	return comment, nil
}

func UpdateComment(applicationId uuid.UUID, commentId int, commentBody string) (Comment, error) {
	comment := Comment{
		ID:            commentId,
		ApplicationID: applicationId,
	}

	if err := db.First(&comment).Error; err != nil {
		return Comment{}, err
	}

	if err := db.Model(&comment).Update("Comment", commentBody).Error; err != nil {
		return Comment{}, err
	}

	return comment, nil
}

func DeleteComment(applicationId uuid.UUID, commentId int) error {
	comment := Comment{
		ID:            commentId,
		ApplicationID: applicationId,
	}

	if err := db.Delete(&comment).Error; err != nil {
		return err
	}
	return nil
}
