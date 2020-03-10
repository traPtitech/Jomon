package model

import (
	"testing"

	"github.com/gofrs/uuid"
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

		comment, err := commentRepo.CreateComment(appId, commentText, userId)
		asr.NoError(err)
		asr.Equal(comment.ApplicationID, appId)
		asr.Equal(comment.Comment, commentText)
		asr.Equal(comment.UserTrapID.TrapId, userId)

		commentText2 := "This is comment 2."

		comment2, err := commentRepo.CreateComment(appId, commentText2, userId)
		asr.NoError(err)
		asr.Equal(comment2.Comment, commentText2)

		getComment, err := commentRepo.GetComment(appId, comment.ID)
		asr.NoError(err)
		asr.Equal(getComment.Comment, commentText)

		getComment2, err := commentRepo.GetComment(appId, comment2.ID)
		asr.NoError(err)
		asr.Equal(getComment2.Comment, commentText2)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		id, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}

		_, err = commentRepo.CreateComment(id, commentText, userId)
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

		comment, err := commentRepo.CreateComment(appId, commentText, userId)
		if err != nil {
			panic(err)
		}

		newCommentText := "This is new comment."

		comment, err = commentRepo.PutComment(appId, comment.ID, newCommentText)
		asr.NoError(err)
		asr.Equal(comment.Comment, newCommentText)

		app, err := repo.GetApplication(appId, true)
		asr.NoError(err)
		asr.NotEqual(app.Comments[0].Comment, commentText)
		asr.Equal(app.Comments[0].Comment, newCommentText)
	})

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		appId, err := repo.createApplication(db, userId)
		if err != nil {
			panic(err)
		}

		_, err = commentRepo.CreateComment(appId, commentText, userId)
		if err != nil {
			panic(err)
		}

		commentText2 := "This is comment 2."

		comment, err := commentRepo.CreateComment(appId, commentText2, userId)
		if err != nil {
			panic(err)
		}

		newCommentText2 := "This is new comment2."

		comment, err = commentRepo.PutComment(appId, comment.ID, newCommentText2)
		asr.NoError(err)
		asr.Equal(comment.Comment, newCommentText2)

		app, err := repo.GetApplication(appId, true)
		asr.NoError(err)
		asr.Len(app.Comments, 2)
		asr.NotEqual(app.Comments[1].Comment, commentText2)
		asr.Equal(app.Comments[1].Comment, newCommentText2)
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

		comment, err := commentRepo.CreateComment(appId, commentText, userId)
		if err != nil {
			panic(err)
		}

		err = commentRepo.DeleteComment(appId, comment.ID)
		asr.NoError(err)

		app, err := repo.GetApplication(appId, true)
		asr.NoError(err)
		asr.Empty(app.Comments)
	})

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		appId, err := repo.createApplication(db, userId)
		if err != nil {
			panic(err)
		}

		_, err = commentRepo.CreateComment(appId, commentText, userId)
		if err != nil {
			panic(err)
		}

		commentText2 := "This is comment 2."

		comment, err := commentRepo.CreateComment(appId, commentText2, userId)
		if err != nil {
			panic(err)
		}

		err = commentRepo.DeleteComment(appId, comment.ID)
		asr.NoError(err)

		app, err := repo.GetApplication(appId, true)
		asr.NoError(err)
		asr.Len(app.Comments, 1)
		asr.Equal(app.Comments[0].Comment, commentText)
	})
}
