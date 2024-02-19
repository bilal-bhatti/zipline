package models

import "time"

// ThingRequest model
type ThingRequest struct {
	Name string `json:"name"`
}

type (
	// ThingResponse model
	ThingResponse struct {
		Name       string    `json:"name"`
		Int        int       `json:"int"`
		Int8       int8      `json:"int8"`
		Int16      int16     `json:"int16"`
		Int32      int32     `json:"int32"`
		Int64      int64     `json:"int64"`
		UInt       uint      `json:"uint"`
		UInt8      uint8     `json:"uint8"`
		UInt16     uint16    `json:"uint16"`
		UInt32     uint32    `json:"uint32"`
		UInt64     uint64    `json:"uint64"`
		Float32    float32   `json:"float32"`
		Float64    float64   `json:"float64"`
		Bool       bool      `json:"bool"`
		CreateDate time.Time `json:"createDate" format:"date-time,2006-01-02"`
		UpdateDate time.Time `json:"updateDate"`
	}

	// //ErrorResponse is the standard error respons format (type block)
	// ErrorResponse struct {
	// 	Code    string `json:"code"`
	// 	Message string `json:"message"`
	// }
)

// ErrorResponse is the standard error respons format (independent)
type ErrorResponse struct {
	Code        int    `json:"code"`
	Status      string `json:"status"`
	Description string `json:"description,omitempty"`
}
