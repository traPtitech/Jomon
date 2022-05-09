package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/file"
)

func (repo *EntRepository) CreateFile(ctx context.Context, name string, mimetype string, requestID uuid.UUID) (*File, error) {
	id := uuid.New()

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

func (repo *EntRepository) GetFile(ctx context.Context, fileID uuid.UUID) (*File, error) {
	file, err := repo.client.File.
		Query().
		Where(file.IDEQ(fileID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return ConvertEntFileToModelFile(file), nil
}

func (repo *EntRepository) DeleteFile(ctx context.Context, fileID uuid.UUID) error {
	return repo.client.File.
		DeleteOneID(fileID).
		Exec(ctx)
}

func ConvertEntFileToModelFile(entfile *ent.File) *File {
	return &File{
		ID:        entfile.ID,
		Name:      entfile.Name,
		MimeType:  entfile.MimeType,
		CreatedAt: entfile.CreatedAt,
	}
}
