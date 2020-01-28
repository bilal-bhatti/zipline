package zipline

import "net/http"

// Post marker for code gen
func Post(handler interface{}) func(w http.ResponseWriter, req *http.Request) {
	panic("Implementation not generated, run zipline")
}

// Get marker for code gen
func Get(handler interface{}) func(w http.ResponseWriter, req *http.Request) {
	panic("Implementation not generated, run zipline")
}
