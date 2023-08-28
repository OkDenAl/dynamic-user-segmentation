package service

import (
	"dynamic-user-segmentation/internal/repository"
	"dynamic-user-segmentation/internal/webapi/gdrive"
)

type Services struct {
	OperationService
	SegmentService
	UserSegmentService
}

func NewServices(repos *repository.Repositories, gDrive gdrive.GDriveApi) *Services {
	return &Services{
		OperationService:   NewOperationService(repos.OperationRepository, gDrive),
		SegmentService:     NewSegmentService(repos.SegmentRepository, repos.UserRepository, repos.UserSegmentRepository),
		UserSegmentService: NewUserSegmentService(repos.UserSegmentRepository),
	}
}
