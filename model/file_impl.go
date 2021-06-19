package model

import (
	"context"
	"fmt"
	"io"
	"mime"

	"github.com/google/uuid"
)

func (repo *EntRepository) CreateFile(ctx context.Context, src io.Reader, name string, mimetype string, requestID uuid.UUID) (*File, error) {
	id := uuid.New()

	ext, err := mime.ExtensionsByType(mimetype)
	if err != nil {
		return nil, err
	} else if len(ext) == 0 {
		return nil, fmt.Errorf("%s is not registered", mimetype)
	}

	filename := fmt.Sprintf("%s%s", id.String(), ext[0])

	err = repo.storage.Save(filename, src)
	if err != nil {
		return nil, err
	}

	created, err := repo.client.File.
		Create().
		SetID(id).
		SetName(name).
		SetMimeType(mimetype).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	_, err = repo.client.Request.
		UpdateOneID(requestID).
		AddFile(created).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	file := &File{
		ID:        created.ID,
		Name:      name,
		MimeType:  mimetype,
		CreatedAt: created.CreatedAt,
	}

	return file, nil
}
