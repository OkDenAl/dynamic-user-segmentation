package gdrive

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var (
	ErrFileNotExists = errors.New("file doesnt exists")
)

type GDriveApi interface {
	UploadCSVFile(ctx context.Context, name string, filedata []byte) (string, error)
	IsAvailable() bool
}

type gDriveApi struct {
	service *drive.Service
}

func New(JSONCredentialsFilePath string) (GDriveApi, error) {
	gDriveService, err := drive.NewService(context.Background(), option.WithCredentialsFile(JSONCredentialsFilePath))
	if err != nil {
		return nil, err
	}
	return &gDriveApi{service: gDriveService}, nil
}

func (g *gDriveApi) IsAvailable() bool {
	return g.service != nil
}

func (g *gDriveApi) UploadCSVFile(ctx context.Context, name string, filedata []byte) (string, error) {
	fileId, err := g.getFileId(ctx, name)
	if err != nil {
		if errors.Is(err, ErrFileNotExists) {
			fileId, err = g.createFile(ctx, name, filedata)
			if err != nil {
				return "", err
			}
			return createURL(fileId), nil
		}
		return "", err
	}
	err = g.updateFile(ctx, fileId, filedata)
	if err != nil {
		return "", err
	}
	return createURL(fileId), nil
}

func (g *gDriveApi) createFile(ctx context.Context, name string, filedata []byte) (string, error) {
	file := &drive.File{
		Name:     name,
		MimeType: "text/csv",
	}
	permissions := &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}
	res, err := g.service.Files.Create(file).Context(ctx).Media(bytes.NewReader(filedata)).Do()
	if err != nil {
		return "", fmt.Errorf("gdrive.createFile - g.service.Files.Create - %w", err)
	}
	_, err = g.service.Permissions.Create(res.Id, permissions).Do()
	if err != nil {
		return "", fmt.Errorf("gdrive.createFile - g.service.Permissions.Create - %w", err)
	}
	return res.Id, nil
}

func (g *gDriveApi) updateFile(ctx context.Context, id string, filedata []byte) error {
	_, err := g.service.Files.Update(id, &drive.File{}).
		Context(ctx).
		Media(bytes.NewReader(filedata)).
		Do()
	if err != nil {
		return fmt.Errorf("gdrive.updateFile - g.service.Files.Update - %w", err)
	}
	return nil
}

func (g *gDriveApi) getFileId(ctx context.Context, name string) (string, error) {
	fileList, err := g.service.Files.List().Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("gdrive.getFileId - g.service.Permissions.Create - %w", err)
	}
	for _, file := range fileList.Files {
		if file.Name == name {
			return file.Id, nil
		}
	}
	return "", ErrFileNotExists
}

func createURL(fileId string) string {
	return fmt.Sprintf("https://drive.google.com/file/d/%s/view?usp=sharing", fileId)
}
