package segment

import (
	"context"
	"dynamic-user-segmentation/internal/repository/segment"
	"dynamic-user-segmentation/internal/repository/user_segment"
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
	segRepo     segment.Repository
	userSegRepo user_segment.Repository
}

func New(segRepo segment.Repository, userSegRepo user_segment.Repository) Service {
	return &service{segRepo: segRepo, userSegRepo: userSegRepo}
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
	err := s.userSegRepo.DeleteSpecSegmentFromUsers(ctx, name)
	if err != nil {
		return err
	}
	return s.segRepo.Delete(ctx, name)
}
