package model

import (
	"context"
	"fmt"
	"io"
	"mime"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/file"
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

func (repo *EntRepository) OpenFile(ctx context.Context, fileID uuid.UUID) (io.ReadCloser, error) {
	file, err := repo.client.File.
		Query().
		Where(file.IDEQ(fileID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	ext, err := mime.ExtensionsByType(file.MimeType)
	if err != nil {
		return nil, err
	} else if len(ext) == 0 {
		return nil, fmt.Errorf("%s is not registered", file.MimeType)
	}

	filename := fmt.Sprintf("%s%s", file.ID.String(), ext[0])

	return repo.storage.Open(filename)
}

func (repo *EntRepository) DeleteFile(ctx context.Context, fileID uuid.UUID, requestID uuid.UUID) error {
	file, err := repo.client.File.
		Query().
		Where(file.IDEQ(fileID)).
		Only(ctx)
	if err != nil {
		return err
	}
	ext, err := mime.ExtensionsByType(file.MimeType)
	if err != nil {
		return err
	} else if len(ext) == 0 {
		return fmt.Errorf("%s is not registered", file.MimeType)
	}

	filename := fmt.Sprintf("%s%s", file.ID.String(), ext[0])

	if err := repo.client.File.
		DeleteOne(file).
		Exec(ctx); err != nil {
		return err
	}

	return repo.storage.Delete(filename)
}

func ConvertEntFileToModelFile(entfile *ent.File) *File {
	return &File{
		ID:        entfile.ID,
		Name:      entfile.Name,
		MimeType:  entfile.MimeType,
		CreatedAt: entfile.CreatedAt,
	}
}
