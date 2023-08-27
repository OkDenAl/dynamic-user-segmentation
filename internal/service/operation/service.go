package operation

import (
	"context"
	"dynamic-user-segmentation/internal/repository/operation"
)

type Service interface {
	MakeReportLink(ctx context.Context, month, year int) (string, error)
}

type service struct {
	opersRepo operation.Repository
}

func New(opersRepo operation.Repository) Service {
	return &service{opersRepo: opersRepo}
}

func (s *service) MakeReportLink(ctx context.Context, month, year int) (string, error) {

	return "link", nil
}
