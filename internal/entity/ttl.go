package entity

type TTL struct {
	Years   uint `json:"years,omitempty"`
	Months  uint `json:"months,omitempty"`
	Days    uint `json:"days,omitempty"`
	Hours   uint `json:"hours,omitempty"`
	Minutes uint `json:"minutes,omitempty"`
	Seconds uint `json:"seconds,omitempty"`
}
