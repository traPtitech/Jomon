package model

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	t.Parallel()

	commentText := "This is comment."
	userId := "userId"

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		appId, err := repo.createApplication(db, userId)
		if err != nil {
			panic(err)
		}

		comment, err := CreateComment(appId, commentText, userId)
		asr.NoError(err)
		asr.Equal(comment.ApplicationID, appId)
		asr.Equal(comment.Comment, commentText)
		asr.Equal(comment.UserTrapID.TrapId, userId)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		id, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}

		_, err = CreateComment(id, commentText, userId)
		asr.Error(err)
	})
}

func TestPutComment(t *testing.T) {
	t.Parallel()

	userId := "userId"
	commentText := "This is comment."

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		appId, err := repo.createApplication(db, userId)
		if err != nil {
			panic(err)
		}

		comment, err := CreateComment(appId, commentText, userId)
		asr.NoError(err)

		newCommentText := "This is new comment."

		comment, err = PutComment(appId, comment.ID, newCommentText)
		asr.NoError(err)
		asr.Equal(comment.Comment, newCommentText)

		app, err := repo.GetApplication(appId, true)
		asr.NoError(err)
		asr.NotEqual(app.Comments[0].Comment, commentText)
		asr.Equal(app.Comments[0].Comment, newCommentText)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		appId, err := repo.createApplication(db, userId)
		if err != nil {
			panic(err)
		}

		_, err = PutComment(appId, int(randSrc.Int63()), userId)
		asr.Error(err)
		asr.True(gorm.IsRecordNotFoundError(err))
	})
}

func TestDeleteComment(t *testing.T) {
	t.Parallel()

	userId := "userId"
	commentText := "This is comment."

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		appId, err := repo.createApplication(db, userId)
		if err != nil {
			panic(err)
		}

		comment, err := CreateComment(appId, commentText, userId)
		asr.NoError(err)

		err = DeleteComment(appId, comment.ID)
		asr.NoError(err)

		app, err := repo.GetApplication(appId, true)
		asr.NoError(err)
		asr.Empty(app.Comments)
	})


	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		appId, err := repo.createApplication(db, userId)
		if err != nil {
			panic(err)
		}

		err = DeleteComment(appId, int(randSrc.Int63()))
		asr.Error(err)
		asr.True(gorm.IsRecordNotFoundError(err))
	})
}