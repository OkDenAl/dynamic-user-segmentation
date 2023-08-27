package segment

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	segRepo "dynamic-user-segmentation/internal/repository/segment"
	"dynamic-user-segmentation/internal/repository/user"
	"dynamic-user-segmentation/internal/repository/user_segment"
	"errors"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=Service --output=../../mocks/service/segmentserv_mocks --outpkg=segmentserv_mocks

var (
	ErrInvalidSegment     = errors.New("invalid segment name")
	ErrInvalidPercentData = errors.New("invalid percent data")
)

type Service interface {
	CreateSegment(ctx context.Context, name string, percentOfUsers float64) error
	DeleteSegment(ctx context.Context, name string) error
}

type service struct {
	segRepo     segRepo.Repository
	userRepo    user.Repository
	userSegRepo user_segment.Repository
}

func New(segRepo segRepo.Repository, userRepo user.Repository, userSegRepo user_segment.Repository) Service {
	return &service{segRepo: segRepo, userRepo: userRepo, userSegRepo: userSegRepo}
}

func (s *service) CreateSegment(ctx context.Context, name string, percentOfUsers float64) error {
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

func (s *service) DeleteSegment(ctx context.Context, name string) error {
	segment := entity.Segment{Name: name}
	if !segment.IsValid() {
		return ErrInvalidSegment
	}
	return s.segRepo.Delete(ctx, segment)
}
