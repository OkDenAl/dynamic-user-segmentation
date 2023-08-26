package segment

type segmentCreatingRequest struct {
	Name           string  `json:"name"`
	PercentOfUsers float64 `json:"percent_of_users,omitempty"`
}

type segmentDeletingRequest struct {
	Name string `json:"name"`
}
