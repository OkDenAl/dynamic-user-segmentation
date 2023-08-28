package service

import (
	"dynamic-user-segmentation/internal/repository"
)

type Services struct {
	OperationService
	SegmentService
	UserSegmentService
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		OperationService:   NewOperationService(repos.OperationRepository),
		SegmentService:     NewSegmentService(repos.SegmentRepository, repos.UserRepository, repos.UserSegmentRepository),
		UserSegmentService: NewUserSegmentService(repos.UserSegmentRepository),
	}
}
