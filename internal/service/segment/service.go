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
	repo segment.Repository
}

func New(repo segment.Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateSegment(ctx context.Context, name string) error {
	if len(name) == 0 {
		return ErrInvalidName
	}
	return s.repo.Create(ctx, name)
}

func (s *service) DeleteSegment(ctx context.Context, name string) error {
	if len(name) == 0 {
		return ErrInvalidName
	}
	return s.repo.Delete(ctx, name)
}
