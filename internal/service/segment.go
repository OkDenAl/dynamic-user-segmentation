package service

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	segRepo "dynamic-user-segmentation/internal/repository/db"
	"errors"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=SegmentService --output=../mocks/service/segmentserv --outpkg=segmentserv_mocks

var (
	ErrInvalidSegment     = errors.New("invalid segment name")
	ErrInvalidPercentData = errors.New("invalid percent data")
)

type SegmentService interface {
	CreateSegment(ctx context.Context, name string, percentOfUsers float64) error
	DeleteSegment(ctx context.Context, name string) error
}

type segmentService struct {
	segRepo     segRepo.SegmentRepository
	userRepo    segRepo.UserRepository
	userSegRepo segRepo.UserSegmentRepository
}

func NewSegmentService(segRepo segRepo.SegmentRepository, userRepo segRepo.UserRepository, userSegRepo segRepo.UserSegmentRepository) SegmentService {
	return &segmentService{segRepo: segRepo, userRepo: userRepo, userSegRepo: userSegRepo}
}

func (s *segmentService) CreateSegment(ctx context.Context, name string, percentOfUsers float64) error {
	segment := entity.Segment{Name: name}
	if !segment.IsValid() {
		return ErrInvalidSegment
	}
	if percentOfUsers < 0 || percentOfUsers > 100 {
		return ErrInvalidPercentData
	}
	err := s.segRepo.Create(ctx, segment)
	if err != nil {
		return err
	}
	if percentOfUsers != 0 {
		usersIds, err := s.userRepo.GetNPercentOfUsersIDs(ctx, percentOfUsers)
		if err != nil {
			return err
		}
		err = s.userSegRepo.CreateOneSegForMultUsers(ctx, usersIds, segment)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *segmentService) DeleteSegment(ctx context.Context, name string) error {
	segment := entity.Segment{Name: name}
	if !segment.IsValid() {
		return ErrInvalidSegment
	}
	return s.segRepo.Delete(ctx, segment)
}
