package segment

import (
	"context"
	"dynamic-user-segmentation/internal/repository/segment"
	"errors"
)

var (
	ErrInvalidName = errors.New("invalid segment name")
)

type Service interface {
	CreateSegment(ctx context.Context, name string) error
	DeleteSegment(ctx context.Context, name string) error
}

type service struct {
	segRepo segment.Repository
}

func New(segRepo segment.Repository) Service {
	return &service{segRepo: segRepo}
}

func (s *service) CreateSegment(ctx context.Context, name string) error {
	if len(name) == 0 {
		return ErrInvalidName
	}
	return s.segRepo.Create(ctx, name)
}

func (s *service) DeleteSegment(ctx context.Context, name string) error {
	if len(name) == 0 {
		return ErrInvalidName
	}
	return s.segRepo.Delete(ctx, name)
}
