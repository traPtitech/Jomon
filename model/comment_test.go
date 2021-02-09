package model

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	t.Parallel()

	commentText := "This is comment."
	userID := "userId"

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		reqID, err := repo.createRequest(db, userID)
		if err != nil {
			panic(err)
		}

		comment, err := commentRepo.CreateComment(reqID, commentText, userID)
		asr.NoError(err)
		asr.Equal(comment.ApplicationID, reqID)
		asr.Equal(comment.Comment, commentText)
		asr.Equal(comment.UserTrapID.TrapID, userID)

		commentText2 := "This is comment 2."

		comment2, err := commentRepo.CreateComment(reqID, commentText2, userID)
		asr.NoError(err)
		asr.Equal(comment2.Comment, commentText2)

		getComment, err := commentRepo.GetComment(reqID, comment.ID)
		asr.NoError(err)
		asr.Equal(getComment.Comment, commentText)

		getComment2, err := commentRepo.GetComment(reqID, comment2.ID)
		asr.NoError(err)
		asr.Equal(getComment2.Comment, commentText2)
	})

	t.Run("shouldFail", func(t *testing.T) {
		asr := assert.New(t)

		id, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}

		_, err = commentRepo.CreateComment(id, commentText, userID)
		asr.Error(err)
	})
}

func TestPutComment(t *testing.T) {
	t.Parallel()

	userID := "userId"
	commentText := "This is comment."

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		reqID, err := repo.createRequest(db, userID)
		if err != nil {
			panic(err)
		}

		comment, err := commentRepo.CreateComment(reqID, commentText, userID)
		if err != nil {
			panic(err)
		}

		newCommentText := "This is new comment."

		comment, err = commentRepo.PutComment(reqID, comment.ID, newCommentText)
		asr.NoError(err)
		asr.Equal(comment.Comment, newCommentText)

		app, err := repo.GetRequest(reqID, true)
		asr.NoError(err)
		asr.NotEqual(app.Comments[0].Comment, commentText)
		asr.Equal(app.Comments[0].Comment, newCommentText)
	})

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		reqID, err := repo.createRequest(db, userID)
		if err != nil {
			panic(err)
		}

		_, err = commentRepo.CreateComment(reqID, commentText, userID)
		if err != nil {
			panic(err)
		}

		commentText2 := "This is comment 2."

		comment, err := commentRepo.CreateComment(reqID, commentText2, userID)
		if err != nil {
			panic(err)
		}

		newCommentText2 := "This is new comment2."

		comment, err = commentRepo.PutComment(reqID, comment.ID, newCommentText2)
		asr.NoError(err)
		asr.Equal(comment.Comment, newCommentText2)

		app, err := repo.GetRequest(reqID, true)
		asr.NoError(err)
		asr.Len(app.Comments, 2)
		asr.NotEqual(app.Comments[1].Comment, commentText2)
		asr.Equal(app.Comments[1].Comment, newCommentText2)
	})
}

func TestDeleteComment(t *testing.T) {
	t.Parallel()

	userID := "userId"
	commentText := "This is comment."

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		reqID, err := repo.createRequest(db, userID)
		if err != nil {
			panic(err)
		}

		comment, err := commentRepo.CreateComment(reqID, commentText, userID)
		if err != nil {
			panic(err)
		}

		err = commentRepo.DeleteComment(reqID, comment.ID)
		asr.NoError(err)

		app, err := repo.GetRequest(reqID, true)
		asr.NoError(err)
		asr.Empty(app.Comments)
	})

	t.Run("shouldSuccess", func(t *testing.T) {
		asr := assert.New(t)

		reqID, err := repo.createRequest(db, userID)
		if err != nil {
			panic(err)
		}

		_, err = commentRepo.CreateComment(reqID, commentText, userID)
		if err != nil {
			panic(err)
		}

		commentText2 := "This is comment 2."

		comment, err := commentRepo.CreateComment(reqID, commentText2, userID)
		if err != nil {
			panic(err)
		}

		err = commentRepo.DeleteComment(reqID, comment.ID)
		asr.NoError(err)

		app, err := repo.GetRequest(reqID, true)
		asr.NoError(err)
		asr.Len(app.Comments, 1)
		asr.Equal(app.Comments[0].Comment, commentText)
	})
}
