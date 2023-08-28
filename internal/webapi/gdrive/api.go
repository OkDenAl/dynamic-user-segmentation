package gdrive

import (
	"context"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GDriveApi interface {
	UploadCSVFile(ctx context.Context) (string, error)
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

func (g *gDriveApi) UploadCSVFile(ctx context.Context) (string, error) {
	return "", nil
}
