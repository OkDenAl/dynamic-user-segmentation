package service

import (
	"bytes"
	"context"
	"dynamic-user-segmentation/internal/repository/db"
	"dynamic-user-segmentation/internal/webapi/gdrive"
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=OperationService --output=../mocks/service/operationserv_mocks --outpkg=operationserv_mocks

var (
	ErrFailedToWriteCSV = errors.New("failed to write csv")
)

type OperationService interface {
	MakeReportLink(ctx context.Context, month, year int) (string, error)
}

type operationService struct {
	operRepo db.OperationRepository
	gDrive   gdrive.GDriveApi
}

func NewOperationService(operRepo db.OperationRepository, gDrive gdrive.GDriveApi) OperationService {
	return &operationService{operRepo: operRepo, gDrive: gDrive}
}

func (s *operationService) MakeReportLink(ctx context.Context, month, year int) (string, error) {
	if !s.gDrive.IsAvailable() {
		return "", errors.New("google drive is not available")
	}
	operations, err := s.operRepo.GetAllOperationsSortedByUserId(ctx, month, year)
	if err != nil {
		return "", err
	}
	b := bytes.Buffer{}
	writer := csv.NewWriter(&b)
	for _, oper := range operations {
		err = writer.Write([]string{strconv.FormatInt(oper.UserId, 10), oper.SegmentName, string(oper.Type), oper.CreatedAt.String()})
		if err != nil {
			return "", fmt.Errorf("OperationService.MakeReportLink - writer.Write - %w", err)
		}
	}
	writer.Flush()
	if err = writer.Error(); err != nil {
		return "", ErrFailedToWriteCSV
	}
	link, err := s.gDrive.UploadCSVFile(ctx, fmt.Sprintf("%d_%d_report.csv", year, month), b.Bytes())
	if err != nil {
		return "", err
	}
	return link, nil
}
