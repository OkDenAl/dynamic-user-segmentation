package service

import (
	"context"
	"dynamic-user-segmentation/internal/repository/db"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=OperationService --output=../../mocks/service/operationserv_mocks --outpkg=operationserv_mocks

type OperationService interface {
	MakeReportLink(ctx context.Context, month, year int) (string, error)
}

type operationService struct {
	opersRepo db.OperationRepository
}

func NewOperationService(opersRepo db.OperationRepository) OperationService {
	return &operationService{opersRepo: opersRepo}
}

func (s *operationService) MakeReportLink(ctx context.Context, month, year int) (string, error) {

	return "link", nil
}
