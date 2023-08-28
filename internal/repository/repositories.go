package repository

import (
	"dynamic-user-segmentation/internal/repository/db"
	"dynamic-user-segmentation/pkg/logging"
	"dynamic-user-segmentation/pkg/postgres"
)

type Repositories struct {
	db.OperationRepository
	db.SegmentRepository
	db.UserRepository
	db.UserSegmentRepository
}

func NewRepositories(pgxPool postgres.PgxPool, log logging.Logger) *Repositories {
	return &Repositories{
		OperationRepository:   db.NewOperationRepo(pgxPool, log),
		SegmentRepository:     db.NewSegmentRepo(pgxPool, log),
		UserRepository:        db.NewUserRepo(pgxPool, log),
		UserSegmentRepository: db.NewUserSegment(pgxPool, log),
	}
}
