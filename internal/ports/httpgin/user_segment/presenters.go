package user_segment

import "dynamic-user-segmentation/internal/entity"

type userSegmentOperationRequest struct {
	UserId           int64      `json:"user_id"`
	SegmentsToAdd    string     `json:"segments_to_add"`
	SegmentsToDelete string     `json:"segments_to_delete"`
	ExpiresAt        entity.TTL `json:"expires_at,omitempty"`
}
