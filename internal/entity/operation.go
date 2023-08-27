package entity

import "time"

type Operation struct {
	UserId      int64
	SegmentName string
	Type        OperationType
	CreatedAt   time.Time
}

type OperationType string

var (
	AddOperation OperationType = "ADD"
	DelOperation OperationType = "DEL"
)
