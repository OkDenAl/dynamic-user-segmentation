package user_segment

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/internal/repository/user_segment"
	"dynamic-user-segmentation/internal/service/segment"
	"errors"
	"strings"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=Service --output=../../mocks/service/usersegmentserv --outpkg=usersegmentserv_mocks

var (
	ErrInvalidUserId = errors.New("invalid user id")
)

type Service interface {
	AddSegmentsToUser(ctx context.Context, userId int64, segments string, ttl entity.TTL) error
	DeleteSegmentsFromUser(ctx context.Context, userId int64, segments string) error
	GetAllUserSegments(ctx context.Context, userId int64) ([]entity.Segment, error)
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
	splitted := strings.Split(segments, ",")
	segmentsArr := make([]entity.Segment, len(splitted))
	for i, segName := range splitted {
		segmentsArr[i] = entity.Segment{Name: segName}
		if !segmentsArr[i].IsValid() {
			return segment.ErrInvalidSegment
		}
	}
	return s.repo.CreateMultSegsForOneUser(ctx, userId, segmentsArr, ttl)
}

func (s *service) DeleteSegmentsFromUser(ctx context.Context, userId int64, segments string) error {
	if userId < 0 {
		return ErrInvalidUserId
	}
	if len(segments) == 0 {
		return nil
	}
	splitted := strings.Split(segments, ",")
	segmentsArr := make([]entity.Segment, len(splitted))
	for i, segName := range splitted {
		segmentsArr[i] = entity.Segment{Name: segName}
		if !segmentsArr[i].IsValid() {
			return segment.ErrInvalidSegment
		}
	}
	return s.repo.DeleteSegmentsFromSpecUser(ctx, userId, segmentsArr)
}

func (s *service) GetAllUserSegments(ctx context.Context, userId int64) ([]entity.Segment, error) {
	if userId < 0 {
		return nil, ErrInvalidUserId
	}
	return s.repo.GetAllUserSegments(ctx, userId)
}
