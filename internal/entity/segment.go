package entity

type Segment struct {
	Name string `json:"name"`
}

func (s *Segment) IsValid() bool {
	return s.Name != ""
}
