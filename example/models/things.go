package models

import "time"

// ThingRequest model
type ThingRequest struct {
	Name string `json:"name"`
}

// ThingResponse model
type ThingResponse struct {
	Name       string    `json:"name"`
	CreateDate time.Time `json:"createDate" format:"date-time,2006-01-02"`
	UpdateDate time.Time `json:"updateDate"`
}
