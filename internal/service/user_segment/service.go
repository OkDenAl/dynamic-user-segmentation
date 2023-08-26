package user_segment

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/internal/repository/user_segment"
	"errors"
	"strings"
)

var (
	ErrInvalidUserId = errors.New("invalid user id")
)

type Service interface {
	AddSegmentsToUser(ctx context.Context, userId int64, segments string, ttl entity.TTL) error
	DeleteSegmentsFromUser(ctx context.Context, userId int64, segments string) error
	GetAllUserSegments(ctx context.Context, userId int64) ([]string, error)
}

type service struct {
	repo user_segment.Repository
}

func New(repo user_segment.Repository) Service {
	return &service{repo: repo}
}

func (s *service) AddSegmentsToUser(ctx context.Context, userId int64, segments string, ttl entity.TTL) error {
	if userId < 0 {
		return ErrInvalidUserId
	}
	if len(segments) == 0 {
		return nil
	}
	return s.repo.CreateMultSegsForOneUser(ctx, userId, strings.Split(segments, ","), ttl)
}

func (s *service) DeleteSegmentsFromUser(ctx context.Context, userId int64, segments string) error {
	if userId < 0 {
		return ErrInvalidUserId
	}
	if len(segments) == 0 {
		return nil
	}
	return s.repo.DeleteSegmentsFromSpecUser(ctx, userId, strings.Split(segments, ","))
}

func (s *service) GetAllUserSegments(ctx context.Context, userId int64) ([]string, error) {
	if userId < 0 {
		return nil, ErrInvalidUserId
	}
	return s.repo.GetAllUserSegments(ctx, userId)
}
