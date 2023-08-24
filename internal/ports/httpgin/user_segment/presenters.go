package user_segment

type userSegmentOperationRequest struct {
	UserId           int64  `json:"user_id"`
	SegmentsToAdd    string `json:"segments_to_add"`
	SegmentsToDelete string `json:"segments_to_delete"`
}
