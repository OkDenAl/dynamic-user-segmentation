package segment

type segmentOperationRequest struct {
	Name           string  `json:"name"`
	PercentOfUsers float64 `json:"percent_of_users,omitempty"`
}
