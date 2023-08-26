package segment

import (
	"context"
	"dynamic-user-segmentation/internal/repository/segment"
	"dynamic-user-segmentation/internal/repository/user"
	"dynamic-user-segmentation/internal/repository/user_segment"
	"errors"
)

var (
	ErrInvalidName        = errors.New("invalid segment name")
	ErrInvalidPercentData = errors.New("invalid percent data")
)

type Service interface {
	CreateSegment(ctx context.Context, name string, percentOfUsers float64) error
	DeleteSegment(ctx context.Context, name string) error
}

type service struct {
	segRepo     segment.Repository
	userRepo    user.Repository
	userSegRepo user_segment.Repository
}

func New(segRepo segment.Repository, userRepo user.Repository, userSegRepo user_segment.Repository) Service {
	return &service{segRepo: segRepo, userRepo: userRepo, userSegRepo: userSegRepo}
}

func (s *service) CreateSegment(ctx context.Context, name string, percentOfUsers float64) error {
	if len(name) == 0 {
		return ErrInvalidName
	}
	if percentOfUsers < 0 || percentOfUsers > 100 {
		return ErrInvalidPercentData
	}
	err := s.segRepo.Create(ctx, name)
	if err != nil {
		return err
	}
	if percentOfUsers != 0 {
		usersIds, err := s.userRepo.GetNPercentOfUsersIDs(ctx, percentOfUsers)
		if err != nil {
			return err
		}
		err = s.userSegRepo.CreateOneSegForMultUsers(ctx, usersIds, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) DeleteSegment(ctx context.Context, name string) error {
	if len(name) == 0 {
		return ErrInvalidName
	}
	return s.segRepo.Delete(ctx, name)
}
